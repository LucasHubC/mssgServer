package gateway

import (
	"context"
	"encoding/json"
	"fmt"
	"mssgServer/infra/log"
	"mssgServer/infra/metrics"
	"mssgServer/infra/redis"
	"mssgServer/infra/rpc"
	wspb "mssgServer/proto"
	"net"
	"net/http"
	"net/url"
	"strings"
	"sync"
	"sync/atomic"
	"time"

	"github.com/gorilla/websocket"
	"google.golang.org/protobuf/proto"
)

// ═══════════════════════════════════════════════════════════
//  WS 智能路由代理 —— 客户端永远只连 Gateway 一个 WS 连接
//  按消息 name 前缀实时路由到对应后端节点
// ═══════════════════════════════════════════════════════════

// ═══════════════════════════════════════════════════════════
//
//	WS 连接建立限流 —— 两层足够，简单可靠
//	L1: 全局并发 WS 连接数（防打爆进程）
//	L2: 单 IP 并发连接数（防单 IP 恶意刷）
//
// ═══════════════════════════════════════════════════════════
const (
	maxGlobalWSConns = 5000 // L1: 全局并发上限
	maxPerIPConns    = 2000 // L2: 单 IP 最大并发（压测单 IP 1000，够；生产改 100）

	// proxy WS 心跳参数（和 net/wsserver.go 保持一致）
	proxyPingInterval = 30 * time.Second
	proxyPongTimeout  = 90 * time.Second // 3 个 Ping 周期没 Pong → 判定死连接
)

var (
	ipConnCount sync.Map // map[string]int64 — 单 IP 当前连接数
)

// tryAcquireWS 尝试获取 WS 连接配额（L1+L2 两层全过才返回 true）
func tryAcquireWS(ip string) (bool, string) {
	// L1: 全局并发数
	if atomic.LoadInt64(&activeProxyConns) >= maxGlobalWSConns {
		return false, "global_full"
	}

	// L2: 单 IP 并发数（CAS 原子操作，避免竞态）
	for {
		var cur int64
		v, loaded := ipConnCount.LoadOrStore(ip, int64(1))
		if !loaded {
			// 第一次存 1，成功
			break
		}
		cur = v.(int64)
		if cur >= maxPerIPConns {
			return false, "ip_full"
		}
		if ipConnCount.CompareAndSwap(ip, cur, cur+1) {
			break
		}
		// CAS 失败，重新循环
	}

	atomic.AddInt64(&activeProxyConns, 1)
	return true, ""
}

// releaseWS 释放 WS 连接配额
func releaseWS(ip string) {
	atomic.AddInt64(&activeProxyConns, -1)
	if v, ok := ipConnCount.Load(ip); ok {
		cnt := v.(int64) - 1
		if cnt <= 0 {
			ipConnCount.Delete(ip)
		} else {
			ipConnCount.Store(ip, cnt)
		}
	}
}

var (
	uidToClientWS    sync.Map // map[string]*websocket.Conn
	kickListenerOnce sync.Once

	// backendsOffline 全局状态标记：后端是否全部离线
	// 只在状态切换时打一条日志（避免 500 个 proxySession 同时刷）
	// false = 至少有一个后端在线，true = 全挂
	backendsOffline atomic.Bool

	// lastPanicLog reader panic 限频：每 5s 最多一条
	lastPanicLog time.Time
	panicMu      sync.Mutex

	uidWarnCount atomic.Int64 // chat 消息没注入 uid 的限频计数（有用，保留）
)

var wsProxyUpgrader = websocket.Upgrader{
	ReadBufferSize:    16384, // 16K 足够 Protobuf 游戏消息
	WriteBufferSize:   16384, // 64K 是浪费,每个 WS 省 96KB
	EnableCompression: false,
	CheckOrigin:       func(r *http.Request) bool { return true },
}

// msgRouteRules 消息 name 前缀 → node type
var msgRouteRules = []struct {
	prefix   string
	nodeType string
}{
	{"account.", "login"},
	{"inventory.", "login"},
	{"room.", "room"},
	{"chat.", "chat"},
	{"group.", "chat"},
}

func routeByName(name string) string {
	for _, r := range msgRouteRules {
		if strings.HasPrefix(name, r.prefix) {
			return r.nodeType
		}
	}
	return "login"
}

var activeProxyConns int64

// ActiveProxyConns 当前活跃的 WS 代理连接数
func ActiveProxyConns() int64 { return atomic.LoadInt64(&activeProxyConns) }

// NodeProvider 后端节点地址提供者
// Gateway 进程内：*Gateway 直接实现（方法调用）
// Broker 进程间：HTTPGatewayClient 实现（HTTP 调 Gateway API）
type NodeProvider interface {
	Route(req rpc.RouteReq) *rpc.RouteResp
}

// proxySession 一个客户端 WS 代理会话
//
// 🔑 核心设计：RCU (Read-Copy-Update)
//
//	backends 用 atomic.Pointer[map] 存储整个 backend map
//	读侧（proxyClientToBackend 每条消息都读）→ atomic.Load → 完全无锁
//	写侧（dialAll 重建，只在 login token 时触发 1 次）→ 创建新 map → Store 替换 → Close 旧 WS
//	旧 map 无人引用后被 GC 回收，无锁、无竞态、读写都高效
type proxySession struct {
	backends atomic.Pointer[map[string]*websocket.Conn]
	done     chan struct{} // 🔴 client WS 断开信号，close 时触发

	// mu 保护：uid/username、lastForwarded、token（写极少）
	mu              sync.Mutex
	service         string
	client          *websocket.Conn
	nodeProvider    NodeProvider
	nodeTypes       []string
	token           string // 🔴 存储客户端初始 token，redial 后端时透传（防止登录丢失）
	uid             string
	username        string
	lastForwarded   map[string]time.Time
	lastCleanup     time.Time
	lastDialAll     time.Time
	dialAllMu       sync.Mutex // 防止多个 goroutine 同时调 dialAll
	backoffFailures int        // dialAll 全部失败的累计次数，用于指数退避
}

// backendMsg fan-in 消息
type backendMsg struct {
	data []byte
	from string
}

// RegisterUID login 成功后注册 uid→client WS
func RegisterUID(uid string, client *websocket.Conn) {
	uidToClientWS.Store(uid, client)
}

// UnregisterUID proxySession Close 时注销
func UnregisterUID(uid string) {
	if uid != "" {
		uidToClientWS.Delete(uid)
	}
}

// StartKickListener 订阅 Redis kick 消息，Close 对应的 client WS
// Broker 进程启动时调用一次（sync.Once 保证只启动一个 goroutine）
func StartKickListener() {
	kickListenerOnce.Do(func() {
		go func() {
			defer func() {
				if r := recover(); r != nil {
					log.Error("[Broker] StartKickListener PANIC", "err", r)
				}
			}()
			redis.SubscribeBroadcast(func(msg redis.BroadcastMsg) {
				if msg.Action != "kick" {
					return
				}
				if ws, ok := uidToClientWS.Load(msg.Target); ok {
					if client, ok2 := ws.(*websocket.Conn); ok2 {
						log.Info("[Broker] 🦵 admin/kick → Close client WS", "uid", msg.Target)
						client.Close()
					}
				}
			})
		}()
	})
}

func newProxySession(service string, client *websocket.Conn, np NodeProvider, nodeTypes []string, token string) *proxySession {
	ps := &proxySession{
		service:       service,
		client:        client,
		nodeProvider:  np,
		nodeTypes:     nodeTypes,
		token:         token,
		lastForwarded: make(map[string]time.Time),
		done:          make(chan struct{}),
	}
	empty := make(map[string]*websocket.Conn)
	ps.backends.Store(&empty)
	return ps
}

// 🔥 全局 dialAll 限流器：30 并发
// 并行 dial 3 backend → 每批 30×3=90 次握手，SYN 队列温和不爆
// 1000 session / 30 = 34 批 × 2s ≈ 68s 总建连时间（可接受）
var dialAllLimiter = make(chan struct{}, 30)

// dialAll 为所有 nodeType 建立 backend WS 连接（可带 token）
// ⚡ 3 个 backend 并行 dial，批次耗时从 6s 降到 2s
func (ps *proxySession) dialAll(token string) {
	dialAllLimiter <- struct{}{}
	defer func() { <-dialAllLimiter }()

	oldPtr := ps.backends.Load()
	newMap := make(map[string]*websocket.Conn, len(ps.nodeTypes))

	var wg sync.WaitGroup
	var mu sync.Mutex
	for _, nt := range ps.nodeTypes {
		wg.Add(1)
		go func(nodeType string) {
			defer wg.Done()
			defer func() {
				if r := recover(); r != nil {
					log.Error("["+ps.service+"] ⚠️ dialAll PANIC", "type", nodeType, "err", r)
				}
			}()
			routeResp := ps.nodeProvider.Route(rpc.RouteReq{Type: nodeType})
			if routeResp.Code != 0 {
				log.Warn("["+ps.service+"] ❌ Route failed", "type", nodeType, "code", routeResp.Code)
				return
			}
			backendWS, err := dialBackend(routeResp.Addr, token)
			if err != nil {
				log.Warn("["+ps.service+"] ❌ dialBackend failed", "type", nodeType, "addr", routeResp.Addr, "err", err)
				return
			}
			mu.Lock()
			newMap[nodeType] = backendWS
			mu.Unlock()
		}(nt)
	}
	wg.Wait()

	if len(newMap) == 0 {
		if !backendsOffline.Swap(true) {
			log.Error("[Broker] 💥 dialAll 全部失败", "token", tokenStr(token != ""))
		}
		if oldPtr != nil {
			for _, ws := range *oldPtr {
				ws.Close()
			}
		}
		empty := make(map[string]*websocket.Conn)
		ps.backends.Store(&empty)
		return
	}

	if backendsOffline.Swap(false) {
		log.Info("[Broker] ✅ dialAll OK", "token", tokenStr(token != ""), "backends", len(newMap))
	}
	ps.backends.Store(&newMap)
	if oldPtr != nil {
		for _, ws := range *oldPtr {
			ws.Close()
		}
	}
}

// GetBackend 读侧完全无锁！atomic.Load 拿指针 → 查 map
func (ps *proxySession) GetBackend(nt string) (*websocket.Conn, bool) {
	m := ps.backends.Load()
	if m == nil {
		return nil, false
	}
	b, ok := (*m)[nt]
	return b, ok
}

// Snapshot 快照当前所有 backend（无锁）
func (ps *proxySession) Snapshot() []*websocket.Conn {
	m := ps.backends.Load()
	if m == nil {
		return nil
	}
	result := make([]*websocket.Conn, 0, len(*m))
	for _, ws := range *m {
		result = append(result, ws)
	}
	return result
}

// Close 关闭所有 backend WS
func (ps *proxySession) Close() {
	ps.mu.Lock()
	uid := ps.uid
	ps.mu.Unlock()
	UnregisterUID(uid)
	if m := ps.backends.Load(); m != nil {
		for _, ws := range *m {
			ws.Close()
		}
	}
}

// clientIP 从 request 里提取客户端 IP（支持 X-Forwarded-For）
func clientIP(r *http.Request) string {
	if xff := r.Header.Get("X-Forwarded-For"); xff != "" {
		// 取第一个 IP
		if idx := strings.IndexByte(xff, ','); idx > 0 {
			return strings.TrimSpace(xff[:idx])
		}
		return strings.TrimSpace(xff)
	}
	ip, _, _ := net.SplitHostPort(r.RemoteAddr)
	return ip
}

// HandleWSProxy 智能 WS 代理入口（Gateway + Broker 共用这个函数！）
// service: 调用方名称，如 "Gateway" 或 "Broker"（用于日志前缀）
// np: NodeProvider — Gateway struct（进程内直接方法调用）或 HTTP 客户端（Broker 进程间 HTTP 调用）
func HandleWSProxy(w http.ResponseWriter, r *http.Request, service string, np NodeProvider) {
	// 🔒 三层限流检查（在 Upgrade 之前拦截，避免无效握手消耗资源）
	ip := clientIP(r)
	if ok, reason := tryAcquireWS(ip); !ok {
		metrics.ErrorsTotal.WithLabelValues("ws_rate_" + reason).Inc()
		log.Warn("["+service+"] 🚫 WS rejected", "ip", ip, "reason", reason, "global", activeProxyConns)
		w.WriteHeader(http.StatusTooManyRequests) // 429
		return
	}
	defer releaseWS(ip)

	clientWS, err := wsProxyUpgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Warn("["+service+"] WS upgrade failed", "err", err)
		return
	}
	metrics.Connections.Inc()
	defer metrics.Connections.Dec()
	defer clientWS.Close()

	// 🔴 WS 僵尸连接自动清理：ReadDeadline + Ping/Pong
	// 和 net/wsserver.go 一样的机制，确保 Ctrl+C 杀进程后 ReadMessage 能超时退出
	clientWS.SetPongHandler(func(string) error {
		_ = clientWS.SetReadDeadline(time.Now().Add(proxyPongTimeout))
		return nil
	})
	_ = clientWS.SetReadDeadline(time.Now().Add(proxyPongTimeout))
	go func() {
		ticker := time.NewTicker(proxyPingInterval)
		defer ticker.Stop()
		for range ticker.C {
			if err := clientWS.WriteControl(websocket.PingMessage, nil, time.Now().Add(5*time.Second)); err != nil {
				return // Ping 写失败（连接已关），退出让 reader 超时清理
			}
		}
	}()

	initialToken := r.URL.Query().Get("token")
	nodeTypes := []string{"login", "chat", "room"}

	ps := newProxySession(service, clientWS, np, nodeTypes, initialToken)
	ps.dialAll(initialToken) // 统一 dialAll，不区分场景
	defer ps.Close()

	if len(ps.Snapshot()) == 0 {
		log.Warn("[" + service + "] ⚠️ no backend available, close client conn")
		return
	}

	log.Info("["+service+"] 🎯 WS proxy connected",
		"client", clientWS.RemoteAddr(),
		"token", tokenStr(initialToken != ""),
		"backends", fmt.Sprintf("%v", mapKeys(ps.backends.Load())))

	// ④ 双向转发
	go proxyClientToBackend(ps)
	proxyBackendToClient(ps)
}

func mapKeys(m *map[string]*websocket.Conn) []string {
	keys := make([]string, 0, len(*m))
	for k := range *m {
		keys = append(keys, k)
	}
	return keys
}

func tokenStr(b bool) string {
	if b {
		return "YES"
	}
	return "NO"
}

// dialBackend 连后端 WS 节点，透传 token 参数
// timeout: 初始 1s（快速失败），重连 2s（更宽容）
func dialBackend(addr, token string) (*websocket.Conn, error) {
	parsed, err := url.Parse(addr)
	if err != nil {
		return nil, err
	}
	if token != "" {
		q := parsed.Query()
		q.Set("token", token)
		parsed.RawQuery = q.Encode()
	}
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	dialer := &websocket.Dialer{
		HandshakeTimeout: 2 * time.Second,
	}
	c, _, err := dialer.DialContext(ctx, parsed.String(), nil)
	if err != nil {
		return nil, err
	}

	// 🔴 backend WS 僵尸连接自动清理
	// login/chat/room 有 pingLoop 定期发 Ping → 这里 PongHandler 续期 deadline
	c.SetPongHandler(func(string) error {
		_ = c.SetReadDeadline(time.Now().Add(proxyPongTimeout))
		return nil
	})
	_ = c.SetReadDeadline(time.Now().Add(proxyPongTimeout))
	return c, nil
}

// proxyClientToBackend 客户端 → 后端（Protobuf 版）
// recover 放在循环内: panic 后 goroutine 继续工作,而不是退出
func proxyClientToBackend(ps *proxySession) {
	defer close(ps.done) // 🔴 退出时通知 proxyBackendToClient 收摊
	for {
		var readErr error
		func() {
			defer func() {
				if r := recover(); r != nil {
					metrics.ErrorsTotal.WithLabelValues("client_reader_panic").Inc()
					log.Error("["+ps.service+"] ⚠️ proxyClientToBackend PANIC", "err", r)
				}
			}()

			_, data, err := ps.client.ReadMessage()
			if err != nil {
				readErr = err
				return
			}

			// 🔴 一次 Unmarshal 搞定：既拿 name 路由，又保留 p 结构体后面直接填 uid
			p := &wspb.WSReq{}
			if proto.Unmarshal(data, p) != nil {
				return // Unmarshal 失败 → 丢弃，继续循环
			}
			metrics.MessagesIn.Inc()

			nt := routeByName(p.Name)
			backend, ok := ps.GetBackend(nt)
			if !ok {
				// ⚠️ backend 暂时不可用 → 丢弃这条消息，但不退出！
				// proxyBackendToClient 的重连循环需要时间重建 backend map
				// 退出 = close(ps.done) → 重连循环也被杀 → 永远恢复不了
				metrics.ErrorsTotal.WithLabelValues("no_backend_" + nt).Inc()
				return
			}

			// 直接在已解码的 p 上填 uid，再 Marshal —— 省掉 injectAuthProtobuf 的二次 Unmarshal！
			ps.mu.Lock()
			uid, username := ps.uid, ps.username
			ps.mu.Unlock()
			if uid != "" {
				if p.AuthUid == "" { // 还没填过（直连客户端没有）
					p.AuthUid = uid
					p.AuthUsername = username
					data, _ = proto.Marshal(p)
				}
			} else if strings.HasPrefix(p.Name, "chat.") {
				if n := uidWarnCount.Add(1); n <= 1 || n%100 == 0 {
					log.Warn("[Broker] ps.uid 为空,chat 消息没注入 uid", "name", p.Name, "total", n)
				}
			}

			if err := writeBackendMsg(backend, data); err != nil {
				// backend 写失败 → 丢弃，不退出
				// 旧 backend WS 可能已被 dialAll Close 了
				metrics.ErrorsTotal.WithLabelValues("backend_write").Inc()
			}
		}()

		// 正常 ReadMessage error = 客户端断开 → 退出
		if readErr != nil {
			log.Debug("["+ps.service+"] proxyClientToBackend exit", "err", readErr)
			return
		}
		// 内部函数正常 return(非 panic)→ for 循环继续下一轮
	}
}

// injectAuthProtobuf 阶段 2：直接填外层 AuthUid/AuthUsername 字段
// 兼容老客户端：如果只发了 msg_json（没 oneof msg），也往 msg_json 里补 __uid 兜底
func injectAuthProtobuf(data []byte, uid, username string) []byte {
	p := &wspb.WSReq{}
	if err := proto.Unmarshal(data, p); err != nil {
		return data
	}

	// ✅ 阶段 2 新方式：外层直接填
	p.AuthUid = uid
	p.AuthUsername = username

	// 🔁 兜底：不管有没有 oneof msg，只要 msg_json 不为空就补 __uid/__username
	// 因为 RespBodyToProto 类型断言全失败 → 全部走 msg_json 兜底
	// chatServer readMsgLoop 同时支持 AuthUid 外层 + msg_json.__uid 两种来源
	if p.MsgJson != "" {
		msg := make(map[string]interface{})
		_ = json.Unmarshal([]byte(p.MsgJson), &msg)
		msg["__uid"] = uid
		msg["__username"] = username
		j, _ := json.Marshal(msg)
		p.MsgJson = string(j)
	}

	out, err := proto.Marshal(p)
	if err != nil {
		return data
	}
	return out
}

// writeClientMsg 写消息到客户端（加写超时，防止 TCP 缓冲区满时永久阻塞）
func writeClientMsg(client *websocket.Conn, data []byte) error {
	_ = client.SetWriteDeadline(time.Now().Add(10 * time.Second))
	return client.WriteMessage(websocket.BinaryMessage, data)
}

// writeBackendMsg 写消息到后端（同样加写超时，对称防护）
func writeBackendMsg(backend *websocket.Conn, data []byte) error {
	_ = backend.SetWriteDeadline(time.Now().Add(10 * time.Second))
	return backend.WriteMessage(websocket.BinaryMessage, data)
}

// proxyBackendToClient 后端 → Broker → 客户端
// fan-in 3 backend reader → msgCh → 批量合并转发 + login 检测 + 写超时
func proxyBackendToClient(ps *proxySession) {
	for {
		// 🔴 client 已断开 → 整个 proxySession 收摊
		select {
		case <-ps.done:
			return
		default:
		}

		// ── Step 1: 快照（原子 Load，无锁！）──
		m := ps.backends.Load()
		if m == nil || len(*m) == 0 {
			// ❌ 不能 return！return = proxySession 永久退出 → 僵尸连接复现
			// 直接内联重连逻辑（此处 respBuf 还没声明，不需要 flushBatch）
			ps.mu.Lock()
			token := ps.token
			ps.mu.Unlock()
			ps.dialAll(token)
			m = ps.backends.Load()
			if m == nil || len(*m) == 0 {
				ps.backoffFailures++
				delay := time.Duration(1<<min(ps.backoffFailures, 4)) * time.Second
				if delay > 30*time.Second {
					delay = 30 * time.Second
				}
				log.Warn("["+ps.service+"] all backends offline, backoff", "delay", delay, "failures", ps.backoffFailures)
				time.Sleep(delay)
				continue
			}
			ps.backoffFailures = 0
			// 重连成功 → 回到循环顶部启动新 reader
			continue
		}
		type entry struct {
			name string
			ws   *websocket.Conn
		}
		entries := make([]entry, 0, len(*m))
		for n, ws := range *m {
			entries = append(entries, entry{n, ws})
		}

		// ── Step 2: fan-in reader goroutines ──
		msgCh := make(chan backendMsg, 2048)
		errCh := make(chan string, len(entries))

		for _, e := range entries {
			go func(name string, ws *websocket.Conn) {
				defer func() {
					select {
					case errCh <- name:
					default:
					}
					if r := recover(); r != nil {
						if now := time.Now(); now.Sub(lastPanicLog) > 5*time.Second {
							lastPanicLog = now
							log.Debug("["+name+"] reader panic", "err", r)
						}
					}
				}()
				for {
					_, data, err := ws.ReadMessage()
					if err != nil {
						return
					}
					select {
					case msgCh <- backendMsg{data: data, from: name}:
					default:
						// msgCh 满了 → 丢弃（极端保护，正常情况下 2048 buffer 够了）
						metrics.ErrorsTotal.WithLabelValues("msgch_full").Inc()
					}
				}
			}(e.name, e.ws)
		}

		// 等所有 reader 退出后 close(msgCh)
		go func() {
			for i := 0; i < len(entries); i++ {
				<-errCh
			}
			close(msgCh)
		}()

		// ── Step 3: 批量合并转发 + 去重 + login 检测 ──
		// Broker 层也做批量合并：把多条后端消息（每条可能已经是 BatchResp）
		// 再合成一个大的 BatchResp 一次 WriteMessage
		const proxyBatchMaxSize = 64
		const proxyBatchFlushWait = 3 * time.Millisecond

		// collectWSResp 把 backendMsg 解码成 WSResp 列表（支持 BatchResp 拆解）
		collectWSResp := func(msg backendMsg) []*wspb.WSResp {
			// 先尝试 BatchResp
			batch := &wspb.BatchResp{}
			if err := proto.Unmarshal(msg.data, batch); err == nil && len(batch.Msgs) > 0 {
				return batch.Msgs
			}
			// 回退单条 WSResp
			p := &wspb.WSResp{}
			if err := proto.Unmarshal(msg.data, p); err == nil {
				return []*wspb.WSResp{p}
			}
			return nil
		}

		// checkLoginAndStore 检查 login 成功 → 存 uid
		checkLoginAndStore := func(rsp *wspb.WSResp) {
			if rsp == nil || rsp.Name != "account.login" || rsp.Code != 0 {
				return
			}
			var u, un string
			if rsp.Msg != nil {
				if lr, ok := rsp.Msg.(*wspb.WSResp_AccountLoginResp); ok {
					u = lr.AccountLoginResp.Uid
					un = lr.AccountLoginResp.Username
				}
			}
			if u == "" && rsp.MsgJson != "" {
				var inner struct {
					UId      string `json:"uid"`
					Username string `json:"username"`
				}
				_ = json.Unmarshal([]byte(rsp.MsgJson), &inner)
				u = inner.UId
				un = inner.Username
			}
			if u != "" {
				ps.mu.Lock()
				if ps.uid == "" {
					ps.uid = u
					ps.username = un
					log.Info("[Broker] login uid 已存储", "uid", u)
					RegisterUID(ps.uid, ps.client)
				}
				ps.mu.Unlock()
			}
		}

		// shouldSkip 20ms 去重检查，返回 true 表示跳过
		shouldSkip := func(key string) bool {
			ps.mu.Lock()
			defer ps.mu.Unlock()
			if last, exists := ps.lastForwarded[key]; exists && time.Since(last) < 20*time.Millisecond {
				return true
			}
			ps.lastForwarded[key] = time.Now()
			if len(ps.lastForwarded) > 10000 || time.Since(ps.lastCleanup) > 5*time.Second {
				cutoff := time.Now().Add(-100 * time.Millisecond)
				for k, t := range ps.lastForwarded {
					if t.Before(cutoff) {
						delete(ps.lastForwarded, k)
					}
				}
				ps.lastCleanup = time.Now()
			}
			return false
		}

		// flushBatch 把累积的 resp 列表合并成 BatchResp → 一次 WriteMessage
		var respBuf []*wspb.WSResp
		flushBatch := func() error {
			if len(respBuf) == 0 {
				return nil
			}
			var data []byte
			var err error
			if len(respBuf) == 1 {
				data, err = proto.Marshal(respBuf[0])
			} else {
				data, err = proto.Marshal(&wspb.BatchResp{Msgs: respBuf})
			}
			if err != nil {
				return err
			}
			err = writeClientMsg(ps.client, data)
			if err != nil {
				metrics.ErrorsTotal.WithLabelValues("client_write").Inc()
				log.Warn("[Broker] ❌ WriteMessage failed", "msgs", len(respBuf), "bytes", len(data), "err", err)
				return err
			}
			metrics.MessagesOut.Add(float64(len(respBuf)))
			respBuf = respBuf[:0]
			return nil
		}

		// 主循环：从 msgCh 取消息，做过滤后累积到 respBuf，满或超时则 flush
		for msg := range msgCh {
			resps := collectWSResp(msg)
			for _, p := range resps {
				// 20ms 去重
				if shouldSkip(p.Name + "|" + string(msg.data)) {
					continue
				}
				// login 检测
				checkLoginAndStore(p)
				respBuf = append(respBuf, p)
			}

			// 满了 → 立即 flush
			if len(respBuf) >= proxyBatchMaxSize {
				if err := flushBatch(); err != nil {
					return // 写失败 → 客户端断连 → proxySession Close 处理
				}
				continue
			}

			// 没满 → 尝试 non-blocking 多拿几条，最多等 proxyBatchFlushWait
			if len(respBuf) == 0 {
				continue // respBuf 空就直接回到外层循环取下一条
			}
			select {
			case nextMsg, ok := <-msgCh:
				if !ok {
					// msgCh 关闭，flush 剩余后退出到重连逻辑
					_ = flushBatch()
					goto reconnect
				}
				nextResps := collectWSResp(nextMsg)
				for _, np := range nextResps {
					if !shouldSkip(np.Name + "|" + string(nextMsg.data)) {
						checkLoginAndStore(np)
						respBuf = append(respBuf, np)
					}
				}
				// 满了立刻 flush
				if len(respBuf) >= proxyBatchMaxSize {
					if err := flushBatch(); err != nil {
						return
					}
				}
			case <-time.After(proxyBatchFlushWait):
				// 超时 flush
				if err := flushBatch(); err != nil {
					return
				}
			}
		}

	reconnect:
		// flush 剩余
		_ = flushBatch()

		// msgCh close = 所有 backend readers 死了 → 重连
		// ✅ 用 ps.token（ProxySession 初始化时存好的初始 token），防止 backend 重建后登录态丢失
		ps.mu.Lock()
		token := ps.token
		ps.mu.Unlock()
		ps.dialAll(token)
		m = ps.backends.Load()
		if m == nil || len(*m) == 0 {
			ps.backoffFailures++
			delay := time.Duration(1<<min(ps.backoffFailures, 4)) * time.Second
			if delay > 30*time.Second {
				delay = 30 * time.Second
			}
			log.Warn("["+ps.service+"] all backends offline, backoff", "delay", delay, "failures", ps.backoffFailures)
			time.Sleep(delay)
			continue
		}
		ps.backoffFailures = 0
	}
}
