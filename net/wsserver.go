package net

import (
	"fmt"
	"mssgServer/infra/log"
	"mssgServer/infra/metrics"
	wspb "mssgServer/proto"
	"sync"
	"sync/atomic"
	"time"

	"github.com/gorilla/websocket"
	"golang.org/x/time/rate"
	"google.golang.org/protobuf/proto"
)

// workerPool 全局消息处理协程池（用 channel 做信号量）
// 上限 2048：单实例 2048 个并发 handler 执行
// 超过时 readMsgLoop 会阻塞等待，但由于 DB 连接池已有 1000，
// 正常情况下不会打满。如果打满说明 DB 是真瓶颈了。
var workerPool = make(chan struct{}, 2048)

// handlerWg 全局 handler WaitGroup — 用于 graceful shutdown 时等所有在跑的 handler 跑完
var handlerWg sync.WaitGroup

// 诊断用 sync.Once：仅打印一次首次收到 Broker 注入的 __uid、首次 BindUid 的日志
var (
	uidMsgDiagOnce    sync.Once
	bindUidDiagOnceWS sync.Once
)

// WaitHandlers 阻塞等所有正在执行的 handler 跑完（graceful shutdown 时调）
func WaitHandlers() {
	handlerWg.Wait()
}

const (
	// pingInterval 心跳发送间隔
	pingInterval = 30 * time.Second
	// pongTimeout 多久没收到 Pong 就认为连接已死
	pongTimeout = 60 * time.Second
)

// WSServer 服务器端
type WSServer struct {
	wsConn    *websocket.Conn
	startTime time.Time // 连接建立时间，用于 ConnectionLifetime 埋点

	router *Router

	outChan chan *WsMsgResp

	// Seq 序列号 —— 不再用 atomic！
	// writeMsgLoop 是单 goroutine 消费 outChan，Seq 自增在那里完成，完全不需要 CAS
	// Seq 序列号 —— 不再用 atomic！
	// writeMsgLoop 是单 goroutine 消费 outChan，Seq 自增在那里完成，完全不需要 CAS
	Seq int64

	property map[string]interface{}

	propertyLock sync.RWMutex

	// limiter 每连接令牌桶限流器
	// 每秒最多 30 条消息，突发最多 50 条
	limiter *rate.Limiter

	// done 通知 wsHandler 连接已关闭
	done chan struct{}

	// once 保证 Close 只执行一次
	once sync.Once

	// writeMu 保护 wsConn 的写操作（gorilla 要求同一连接同一时刻只有一个 writer）
	// writeMsgLoop / 心跳 Ping / PingHandler 回 Pong 都要走这把锁
	writeMu sync.Mutex

	// sendMu 保护 outChan 的发送 + closed 标志，避免向已关闭的 channel 发送导致 panic
	// Close 必须在持锁状态下置 closed=true 再 close(outChan)
	sendMu sync.Mutex
	closed bool

	// NodeAddr 当前节点的 ws 地址（如 "ws://127.0.0.1:8003"）
	// 由 Server.wsHandler 赋值，Redis 跨节点在线状态同步用
	NodeAddr string
}

func newWsServer(wsConn *websocket.Conn) *WSServer {
	return &WSServer{
		startTime: time.Now(),
		wsConn:    wsConn,
		outChan:   make(chan *WsMsgResp, 16384), // 16384 指针 × 8B = 128KB/conn，扛 2s world 广播洪峰，比 32768 省 128KB/conn
		property:  make(map[string]interface{}),
		Seq:       0,
		limiter:   rate.NewLimiter(100, 500), // 正常玩家 1s 5-10 条足够，100/s 防刷 DB，突发 500 条扛洪峰
		done:      make(chan struct{}),
	}
}

func (s *WSServer) Router(router *Router) {
	s.router = router

}

func (s *WSServer) SetProperty(key string, value interface{}) {
	s.propertyLock.Lock()
	defer s.propertyLock.Unlock()
	s.property[key] = value
}

func (s *WSServer) GetProperty(key string) (interface{}, error) {
	s.propertyLock.RLock()
	defer s.propertyLock.RUnlock()
	if value, ok := s.property[key]; ok {
		return value, nil
	}
	return nil, fmt.Errorf("property %s not found", key)
}

func (s *WSServer) RemoveProperty(key string) {
	s.propertyLock.Lock()
	defer s.propertyLock.Unlock()
	delete(s.property, key)
}
func (s *WSServer) Addr() string {
	return s.wsConn.RemoteAddr().String()
}

// writeMessage 加锁后写 WebSocket，避免 writeMsgLoop / 心跳 / PingHandler 并发写冲突
// gorilla/websocket 要求：同一连接同一时刻只能有一个 writer，SetWriteDeadline 也要在同把锁内
func (s *WSServer) writeMessage(msgType int, data []byte) error {
	s.writeMu.Lock()
	defer s.writeMu.Unlock()
	_ = s.wsConn.SetWriteDeadline(time.Now().Add(10 * time.Second))
	return s.wsConn.WriteMessage(msgType, data)
}

// safeSend 向 outChan 投递消息，连接已关闭或通道满时返回 false（绝不 panic）
// sendMu 与 Close 互斥，保证不会向已关闭的 channel 发送
var dropLogCount int64 // 全局丢包计数（限频用）

func (s *WSServer) safeSend(msg *WsMsgResp) bool {
	s.sendMu.Lock()
	if s.closed {
		s.sendMu.Unlock()
		return false
	}
	select {
	case s.outChan <- msg:
		s.sendMu.Unlock()
		return true
	default:
		s.sendMu.Unlock()
		metrics.DroppedMessages.WithLabelValues("outchan_full").Inc()
		if atomic.AddInt64(&dropLogCount, 1)%100 == 0 {
			log.Warn("outChan full, drop message", "addr", s.Addr(), "name", msg.Body.Name, "totalDropped", atomic.LoadInt64(&dropLogCount))
		}
		return false
	}
}

// Push 向本连接推送消息（非阻塞，outChan 满或连接已关时丢弃）
// Seq 在 writeMsgLoop 里统一赋值（单 goroutine，不需要 atomic）
func (s *WSServer) Push(name string, data interface{}) {
	rsp := &WsMsgResp{
		Body: &RespBody{
			Seq:  0, // 占位，writeMsgLoop 里统一填
			Name: name,
			Code: 0,
			Msg:  data,
		},
	}
	s.safeSend(rsp)
}

// 通道一旦建立 那么收发消息都要持续去监听通道
func (s *WSServer) Start() {
	// ✅ 加 Pong 处理：收到 Pong 时延长 ReadDeadline
	s.wsConn.SetPongHandler(func(string) error {
		_ = s.wsConn.SetReadDeadline(time.Now().Add(pongTimeout))
		return nil
	})
	// 初始 ReadDeadline
	_ = s.wsConn.SetReadDeadline(time.Now().Add(pongTimeout))

	go s.pingLoop() // 心跳
	go s.writeMsgLoop()
	go s.readMsgLoop()
}

// pingLoop 心跳循环 — 每 30s 发一次 Ping，连接断开时退出
func (s *WSServer) pingLoop() {
	ticker := time.NewTicker(pingInterval)
	defer ticker.Stop()
	for {
		select {
		case <-ticker.C:
			if err := s.writeMessage(websocket.PingMessage, nil); err != nil {
				log.Warn("heartbeat ping failed", "addr", s.Addr(), "err", err)
				s.Close()
				return
			}
		case <-s.done:
			return // 连接已关闭
		}
	}
}

const (
	batchMaxSize   = 32                   // 每帧最多合并 32 条，减少一半 TCP syscall
	batchFlushWait = 5 * time.Millisecond // 没满但等了 5ms 立刻刷，低负载也能攒更多
)

// writeMsgLoop 写连接 — 支持批量合并减少 socket syscall
// 每 batchMaxSize 条或 batchFlushWait 超时，合并成 BatchResp 一次写入
func (s *WSServer) writeMsgLoop() {
	batch := make([]*wspb.WSResp, 0, batchMaxSize)

	flush := func() {
		if len(batch) == 0 {
			return
		}
		// 🔴 统一给每条消息赋 Seq（单 goroutine，直接 s.Seq++）
		for _, p := range batch {
			p.Seq = s.Seq
			s.Seq++
		}
		metrics.BatchSize.Observe(float64(len(batch)))
		var data []byte
		var err error
		if len(batch) == 1 {
			// 单条直接发 WSResp，不用 BatchResp 包装（省 proto 编解码开销）
			data, err = proto.Marshal(batch[0])
		} else {
			// 多条合并成 BatchResp
			data, err = proto.Marshal(&wspb.BatchResp{Msgs: batch})
		}
		if err == nil {
			writeStart := time.Now()
			if werr := s.writeMessage(websocket.BinaryMessage, data); werr != nil {
				metrics.ErrorsTotal.WithLabelValues("write").Inc()
				log.Error("write message failed", "addr", s.Addr(), "err", werr)
			}
			metrics.WriteDuration.Observe(time.Since(writeStart).Seconds())
			metrics.MessagesOut.Add(float64(len(batch)))
			metrics.BytesOut.Add(float64(len(data)))
		}
		batch = batch[:0]
	}

	for {
		// 拿第一条
		msg, ok := <-s.outChan
		if !ok {
			flush()
			return
		}
		p, err := RespBodyToProto(msg.Body)
		if err != nil {
			log.Error("proto convert resp failed", "err", err)
		} else {
			batch = append(batch, p)
		}

		// 尝试再拿 batchMaxSize-1 条（非阻塞 + 超时）
		for len(batch) < batchMaxSize {
			select {
			case msg2, ok := <-s.outChan:
				if !ok {
					flush()
					return
				}
				p2, err2 := RespBodyToProto(msg2.Body)
				if err2 == nil {
					batch = append(batch, p2)
				}
			case <-time.After(batchFlushWait):
				flush()
				goto next
			}
		}
		flush()
	next:
	}
}

func (s *WSServer) readMsgLoop() {
	//先读消息 再处理消息

	defer func() {
		if err := recover(); err != nil {
			log.Error("readMsgLoop panic", "addr", s.Addr(), "err", err)
			s.Close()
		}
	}()

	for {
		_, messageData, err := s.wsConn.ReadMessage() // 阻塞等待客户端发来消息

		if err != nil {
			metrics.ErrorsTotal.WithLabelValues("read").Inc()
			log.Warn("read message failed", "addr", s.Addr(), "err", err)
			break
		}
		metrics.MessagesIn.Inc()
		metrics.BytesIn.Add(float64(len(messageData)))

		// Protobuf 解码
		p := &wspb.WSReq{}
		err = proto.Unmarshal(messageData, p)
		if err != nil {
			log.Warn("proto unmarshal message failed", "addr", s.Addr(), "err", err)
			break
		}
		body, err := ReqBodyFromProto(p)
		if err != nil {
			log.Warn("ReqBodyFromProto failed", "addr", s.Addr(), "err", err)
			break
		}

		// ── 做法 A 核心：从 Broker 注入的 uid 自动设 WS Property + BindUid ──
		// ⚠️ 双重来源：
		//   ① 阶段 2 新方式: 外层 body.AuthUid (Broker injectAuthProtobuf 直接填)
		//   ② 阶段 1 兜底:   msg_json 里的 __uid (老客户端)
		// 直连场景 (HTTP Gateway 或测试工具): 两个都没有，跳过
		// ⚠️ SetProperty ≠ BindUid！前者让 handler 能读 uid，后者让 admin/kick 能找到 WS
		var uid, username string
		uid = body.AuthUid // ✅ 阶段 2：优先读外层鉴权字段
		if uid == "" && body.Msg != nil {
			// 🔁 兜底：从 msg_json 解出来的 map 里找 __uid(阶段 1)
			if m, ok := body.Msg.(map[string]interface{}); ok {
				if u, _ := m["__uid"].(string); u != "" {
					uid = u
					if un, _ := m["__username"].(string); un != "" {
						username = un
					}
				}
			}
		}
		if uid != "" {
			s.SetProperty("uid", uid)
			if username != "" {
				s.SetProperty("username", username)
			}
			// BindUid 只绑定一次（同一 WS 收到多条带 uid 的消息只第一次生效）
			if _, err := s.GetProperty("__bound"); err != nil {
				ConnManager.BindUid(uid, s, true) // skipRedis=true: 后端代理只写本机 uidMap
				s.SetProperty("__bound", true)
			}
		}

		// ✅ 限流检查：Allow 不阻塞，超限直接断开
		if !s.limiter.Allow() {
			log.Warn("rate limit exceeded, disconnecting", "addr", s.Addr())
			break
		}

		// 路由（异步化：丢到 worker pool 执行，不阻塞 read loop）
		req := &WsMsgReq{Body: body, Conn: s}
		rsp := &WsMsgResp{Body: &RespBody{Name: body.Name, Seq: req.Body.Seq}}

		select {
		case workerPool <- struct{}{}:
			// 抢到 worker，异步执行 handler
			handlerWg.Add(1)
			go func() {
				defer handlerWg.Done()
				defer func() { <-workerPool }()
				metrics.ActiveHandlers.Inc()
				handlerStart := time.Now()
				defer metrics.ActiveHandlers.Dec()
				defer func() {
					if r := recover(); r != nil {
						metrics.HandlersPanic.Inc()
						metrics.ErrorsTotal.WithLabelValues("handler_panic").Inc()
						log.Error("handler panic", "addr", s.Addr(), "name", body.Name, "err", r)
					}
				}()
				s.router.Run(req, rsp)
				metrics.HandlerDuration.Observe(time.Since(handlerStart).Seconds())
				// 响应通过 safeSend 投递，连接已关或 outChan 满时安全跳过
				if !s.safeSend(rsp) {
					metrics.DroppedMessages.WithLabelValues("response_outchan_full").Inc()
					log.Warn("outChan full or closed, drop response", "addr", s.Addr(), "name", body.Name)
				}
			}()
		default:
			// worker pool 满了，降级为同步执行（极端保护）
			log.Warn("worker pool full, fallback to sync", "addr", s.Addr())
			handlerWg.Add(1)
			func() {
				defer handlerWg.Done()
				s.router.Run(req, rsp)
				s.safeSend(rsp)
			}()
		}
	}

	s.Close()
}

// Close 幂等关闭：sendMu 锁住 → closed=true → close(outChan) → 解锁 → 关 done → 关 ws
func (s *WSServer) Close() {
	s.once.Do(func() {
		metrics.ConnectionLifetime.Observe(time.Since(s.startTime).Seconds())
		s.sendMu.Lock()
		s.closed = true
		close(s.outChan)
		s.sendMu.Unlock()
		close(s.done)
		s.wsConn.Close()
	})
}

// WaitClose 阻塞等连接关闭（wsHandler 里用）
func (s *WSServer) WaitClose() {
	<-s.done
}
