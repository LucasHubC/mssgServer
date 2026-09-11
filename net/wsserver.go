package net

import (
	"encoding/json"
	"mssgServer/infra/log"
	"sync"
	"sync/atomic"
	"time"

	"github.com/gorilla/websocket"
	"golang.org/x/time/rate"
)

// workerPool 全局消息处理协程池（用 channel 做信号量）
// 上限 2048：单实例 2048 个并发 handler 执行
// 超过时 readMsgLoop 会阻塞等待，但由于 DB 连接池已有 1000，
// 正常情况下不会打满。如果打满说明 DB 是真瓶颈了。
var workerPool = make(chan struct{}, 2048)

// WSServer 服务器端
type WSServer struct {
	wsConn *websocket.Conn

	router *Router

	outChan chan *WsMsgResp

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
		wsConn:   wsConn,
		outChan:  make(chan *WsMsgResp, 256),
		property: make(map[string]interface{}),
		Seq:      0,
		limiter:  rate.NewLimiter(30, 50),
		done:     make(chan struct{}),
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
	return nil, nil
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
func (s *WSServer) safeSend(msg *WsMsgResp) bool {
	s.sendMu.Lock()
	defer s.sendMu.Unlock()
	if s.closed {
		return false
	}
	select {
	case s.outChan <- msg:
		return true
	default:
		log.Warn("outChan full, drop message", "addr", s.Addr())
		return false
	}
}

// Push 向本连接推送消息（非阻塞，outChan 满或连接已关时丢弃并记日志）
func (s *WSServer) Push(name string, data interface{}) {
	seq := atomic.AddInt64(&s.Seq, 1) // 🔴 并发安全的序列号
	rsp := &WsMsgResp{
		Body: &RespBody{
			Seq:  seq,
			Name: name,
			Code: 0,
			Msg:  data,
		},
	}
	if !s.safeSend(rsp) {
		log.Warn("Push: outChan full or closed, drop message", "addr", s.Addr(), "name", name)
	}
}

// 通道一旦建立 那么收发消息都要持续去监听通道
func (s *WSServer) Start() {
	go s.writeMsgLoop()
	go s.readMsgLoop()
}

func (s *WSServer) writeMsgLoop() {
	for msg := range s.outChan {
		data, err := json.Marshal(msg.Body)
		if err != nil {
			log.Error("marshal response failed", "err", err)
			continue
		}
		if err := s.writeMessage(websocket.TextMessage, data); err != nil {
			log.Error("write message failed", "addr", s.Addr(), "err", err)
			return
		}
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
			log.Warn("read message failed", "addr", s.Addr(), "err", err)
			break
		}

		log.Debug("recv raw", "addr", s.Addr(), "data", string(messageData))

		body := &ReqBody{}
		err = json.Unmarshal(messageData, body)
		if err != nil {
			log.Warn("unmarshal message failed", "addr", s.Addr(), "err", err, "raw", string(messageData))
			break
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
			go func() {
				defer func() { <-workerPool }()
				defer func() {
					if r := recover(); r != nil {
						log.Error("handler panic", "addr", s.Addr(), "name", body.Name, "err", r)
					}
				}()
				s.router.Run(req, rsp)
				// 响应通过 safeSend 投递，连接已关或 outChan 满时安全跳过
				if !s.safeSend(rsp) {
					log.Warn("outChan full or closed, drop response", "addr", s.Addr(), "name", body.Name)
				}
			}()
		default:
			// worker pool 满了，降级为同步执行（极端保护）
			log.Warn("worker pool full, fallback to sync", "addr", s.Addr())
			s.router.Run(req, rsp)
			s.safeSend(rsp)
		}
	}

	s.Close()
}

// Close 关闭连接，幂等（sync.Once 保证只执行一次）
// 关闭顺序：先在 sendMu 锁内置 closed=true 并 close(outChan)（让 writeMsgLoop 退出，
// 同时阻止后续 safeSend），再关 wsConn，再关 done 通知 wsHandler
func (s *WSServer) Close() {
	s.once.Do(func() {
		s.sendMu.Lock()
		s.closed = true
		close(s.outChan)
		s.sendMu.Unlock()

		s.wsConn.Close()
		close(s.done)
	})
}

// WaitClose 阻塞等连接关闭（wsHandler 里用）
func (s *WSServer) WaitClose() {
	<-s.done
}
