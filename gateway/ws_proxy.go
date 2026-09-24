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

// uidToClientWS Broker 侧 uid→client WS 注册表
// 做法 A 后 admin/kick 必须踢 client WS（踢 backend WS 只会触发 dialAll 重连，client 不会断）
// Broker login 成功时 RegisterUID，proxySession Close 时 UnregisterUID
// StartKickListener() 启动一个 goroutine 订阅 Redis kick 消息
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

	// 🔍 临时诊断：injectAuth 是否真的执行了（限频打印）
	diagInjected    atomic.Bool
	diagNoUid       atomic.Bool
	diagLoginStored atomic.Bool
	uidWarnCount    atomic.Int64 // chat 消息没注入 uid 的限频计数

	// 🔍 dialAll 诊断计数器：定位 no_backend 根因
	diagDialRouteFail int64 // dialAll 中 Route 返回非零 code 的次数
	diagDialConnFail  int64 // dialAll 中 dialBackend 网络错误次数
	diagDialAllEmpty  int64 // dialAll 全部 backend 都失败的次数
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

// 🔍 诊断计数器 - client reader 侧(proxyClientToBackend)
var (
	diagClientReaderStart int64
	diagClientReaderExit  int64
	diagClientReaderPanic int64
)

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

	// mu 保护：uid/username、lastForwarded（写极少）
	mu            sync.Mutex
	service       string
	client        *websocket.Conn
	nodeProvider  NodeProvider
	nodeTypes     []string
	uid           string
	username      string
	lastForwarded map[string]time.Time
	lastCleanup   time.Time
	lastDialAll   time.Time
	dialAllMu     sync.Mutex // 防止多个 goroutine 同时调 dialAll
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
		go redis.SubscribeBroadcast(func(msg redis.BroadcastMsg) {
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
	})
}

func newProxySession(service string, client *websocket.Conn, np NodeProvider, nodeTypes []string) *proxySession {
	ps := &proxySession{
		service:       service,
		client:        client,
		nodeProvider:  np,
		nodeTypes:     nodeTypes,
		lastForwarded: make(map[string]time.Time),
	}
	empty := make(map[string]*websocket.Conn)
	ps.backends.Store(&empty)
	return ps
}

// 🔥 全局 dialAll 限流器：最多同时 50 个 session 在重建 backend
// 防止 1000 个 session 同时 Route + dial 打爆 Gateway 和后端
var dialAllLimiter = make(chan struct{}, 50)

// dialAll 为所有 nodeType 建立 backend WS 连接（可带 token）
// ⚡ 串行 + 全局限流：保证 Gateway 不被打爆
func (ps *proxySession) dialAll(token string) {
	// 🔴 全局限流！先拿 token
	dialAllLimiter <- struct{}{}
	defer func() { <-dialAllLimiter }()

	oldPtr := ps.backends.Load()
	newMap := make(map[string]*websocket.Conn, len(ps.nodeTypes))
	for _, nt := range ps.nodeTypes {
		routeResp := ps.nodeProvider.Route(rpc.RouteReq{Type: nt})
		if routeResp.Code != 0 {
			atomic.AddInt64(&diagDialRouteFail, 1)
			log.Warn("["+ps.service+"] ❌ Route failed", "type", nt, "code", routeResp.Code)
			continue
		}
		backendWS, err := dialBackend(routeResp.Addr, token)
		if err != nil {
			atomic.AddInt64(&diagDialConnFail, 1)
			log.Warn("["+ps.service+"] ❌ dialBackend failed", "type", nt, "addr", routeResp.Addr, "err", err)
			continue
		}
		newMap[nt] = backendWS
		log.Info("["+ps.service+"] ✅ dialBackend success", "type", nt, "addr", routeResp.Addr)
	}

	if len(newMap) == 0 {
		atomic.AddInt64(&diagDialAllEmpty, 1)
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

// HandleWSProxy 智能 WS 代理入口（Gateway + Broker 共用这个函数！）
// service: 调用方名称，如 "Gateway" 或 "Broker"（用于日志前缀）
// np: NodeProvider — Gateway struct（进程内直接方法调用）或 HTTP 客户端（Broker 进程间 HTTP 调用）
func HandleWSProxy(w http.ResponseWriter, r *http.Request, service string, np NodeProvider) {
	// ① Upgrade 客户端连接
	clientWS, err := wsProxyUpgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Warn("["+service+"] WS upgrade failed", "err", err)
		return
	}
	atomic.AddInt64(&activeProxyConns, 1)
	defer atomic.AddInt64(&activeProxyConns, -1)
	metrics.Connections.Inc()
	defer metrics.Connections.Dec()
	defer clientWS.Close()

	// ② 初始 token（客户端可能没带，后续 account.login 成功后重建）
	initialToken := r.URL.Query().Get("token")
	nodeTypes := []string{"login", "chat", "room"}

	// ③ 建 proxySession + dial 初始 backend
	ps := newProxySession(service, clientWS, np, nodeTypes)
	ps.dialAll(initialToken)
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
// ⚠️ 必须带超时！否则后端全挂时 500 个 dial 会同时卡在 TCP connect 里 → goroutine 爆炸
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
	ctx, cancel := context.WithTimeout(context.Background(), 500*time.Millisecond)
	defer cancel()
	dialer := &websocket.Dialer{
		HandshakeTimeout: 500 * time.Millisecond,
	}
	c, _, err := dialer.DialContext(ctx, parsed.String(), nil)
	return c, err
}

// proxyClientToBackend 客户端 → 后端（Protobuf 版）
// recover 放在循环内: panic 后 goroutine 继续工作,而不是退出
func proxyClientToBackend(ps *proxySession) {
	atomic.AddInt64(&diagClientReaderStart, 1)
	defer atomic.AddInt64(&diagClientReaderExit, 1)
	for {
		var readErr error
		func() {
			defer func() {
				if r := recover(); r != nil {
					atomic.AddInt64(&diagClientReaderPanic, 1)
					metrics.ErrorsTotal.WithLabelValues("client_reader_panic").Inc()
					log.Error("["+ps.service+"] ⚠️ proxyClientToBackend PANIC", "err", r)
				}
			}()

			_, data, err := ps.client.ReadMessage()
			if err != nil {
				readErr = err
				return
			}

			// Protobuf 解码拿到 name 路由
			p := &wspb.WSReq{}
			if proto.Unmarshal(data, p) != nil {
				return
			}
			metrics.MessagesIn.Inc()

			nt := routeByName(p.Name)
			backend, ok := ps.GetBackend(nt)
			if !ok {
				metrics.ErrorsTotal.WithLabelValues("no_backend_" + nt).Inc()
				return
			}

			// 如果已经 login，注入 __uid/__username
			ps.mu.Lock()
			uid, username := ps.uid, ps.username
			ps.mu.Unlock()
			if uid != "" {
				data = injectAuthProtobuf(data, uid, username)
			} else if strings.HasPrefix(p.Name, "chat.") {
				// 限频:首次 + 每 100 条打 1 条
				if n := uidWarnCount.Add(1); n <= 1 || n%100 == 0 {
					log.Warn("[Broker] ps.uid 为空,chat 消息没注入 uid", "name", p.Name, "total", n)
				}
			}

			// BinaryMessage — Protobuf 是二进制
			if err := backend.WriteMessage(websocket.BinaryMessage, data); err != nil {
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

// proxyBackendToClient 后端 → Gateway → 客户端
// fan-in 3 backend reader → msgCh → select 转发 + 检测 login token
func proxyBackendToClient(ps *proxySession) {
	for {
		// ── Step 1: 快照（原子 Load，无锁！）──
		m := ps.backends.Load()
		if m == nil || len(*m) == 0 {
			log.Info("[" + ps.service + "] no backends, exit")
			return
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
		// 每个 backend 一个 reader,退出时通过 errCh 发信号(省掉 WaitGroup 辅助 goroutine)
		msgCh := make(chan backendMsg, 2048)
		cancel := make(chan struct{})
		errCh := make(chan string, len(entries))

		for _, e := range entries {
			go func(name string, ws *websocket.Conn) {
				defer func() {
					// 退出时通知主循环
					select {
					case errCh <- name:
					default:
					}
					// recover 捕获 gorilla 并发 Close+Read 的 panic
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
					case <-cancel:
						return
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

		// ── Step 3: 转发 + 去重 + 存 login uid（做法 A: WS 只是管道）──
		for msg := range msgCh {
			// Protobuf 解码：先尝试 BatchResp(多条合并)，再回退 WSResp(单条)
			p := &wspb.WSResp{}
			batch := &wspb.BatchResp{}
			if err := proto.Unmarshal(msg.data, batch); err == nil && len(batch.Msgs) > 0 {
				p = batch.Msgs[0]
			} else if err := proto.Unmarshal(msg.data, p); err != nil {
				continue
			}

			// 20ms 去重
			key := p.Name + "|" + string(msg.data)
			ps.mu.Lock()
			if last, exists := ps.lastForwarded[key]; exists && time.Since(last) < 20*time.Millisecond {
				ps.mu.Unlock()
				log.Info("[Broker] 🚫 去重跳过", "from", msg.from, "name", p.Name, "key", key[:30])
				continue
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
			ps.mu.Unlock()

			// ✅ 关键:不能把 BatchResp 原始 bytes 透传给 client
			// client 只认单条 WSResp!backend 用 BatchResp 合并发送只是后端间优化
			metrics.MessagesOut.Inc()
			if len(batch.Msgs) > 1 {
				// BatchResp → 拆成单条 WSResp 重新序列化逐条发
				for _, sub := range batch.Msgs {
					if sub == nil {
						continue
					}
					if subBytes, err := proto.Marshal(sub); err == nil {
						err = ps.client.WriteMessage(websocket.BinaryMessage, subBytes)
						if err != nil {
							metrics.ErrorsTotal.WithLabelValues("client_write").Inc()
							log.Warn("[Broker] ❌ WriteMessage failed", "from", msg.from, "name", sub.Name, "err", err)
							return
						}
						log.Debug("[Broker] WriteMessage → client", "from", msg.from, "name", sub.Name, "uid", ps.uid, "bytes", len(subBytes))
					}
				}
			} else {
				// 单条 WSResp(或解 BatchResp 失败 → 直接透传)
				err := ps.client.WriteMessage(websocket.BinaryMessage, msg.data)
				if err != nil {
					metrics.ErrorsTotal.WithLabelValues("client_write").Inc()
					log.Warn("[Broker] ❌ WriteMessage failed", "from", msg.from, "name", p.Name, "err", err)
					return
				}
				log.Debug("[Broker] WriteMessage → client", "from", msg.from, "name", p.Name, "uid", ps.uid, "bytes", len(msg.data))
			}

			// login 成功 → 记住 uid/username
			// ⚠️ 必须遍历 BatchResp 所有条目 —— 并发压测下 login 可能不是第一条
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
			// 先检查取出来的 p(第一条)
			checkLoginAndStore(p)
			// 如果是 BatchResp,遍历其余条目
			if len(batch.Msgs) > 1 {
				for i := 1; i < len(batch.Msgs); i++ {
					checkLoginAndStore(batch.Msgs[i])
				}
			}
		}

		// msgCh close = 所有 backend readers 死了 → 重连
		// 注:close(cancel) 和 wg.Wait() 已不需要 —— msgCh close 本身就是所有 reader 退出后的信号
		ps.dialAll("") // 不带 token，因为后端 server 鉴权靠消息里的 __uid
		m = ps.backends.Load()
		if m == nil || len(*m) == 0 {
			// 全后端离线：静默退避（状态日志已在 dialAll 里打过了）
			time.Sleep(2 * time.Second)
			continue
		}
	}
}
