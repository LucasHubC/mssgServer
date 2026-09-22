package gateway

import (
	"context"
	"encoding/json"
	"fmt"
	"mssgServer/infra/log"
	"mssgServer/infra/redis"
	"mssgServer/infra/rpc"
	"net/http"
	"net/url"
	"strings"
	"sync"
	"sync/atomic"
	"time"

	"github.com/gorilla/websocket"
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
)

var wsProxyUpgrader = websocket.Upgrader{
	ReadBufferSize:    65536,
	WriteBufferSize:   65536,
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

	// mu 保护：uid/username、lastForwarded（写极少）
	mu            sync.Mutex
	service       string
	client        *websocket.Conn
	nodeProvider  NodeProvider
	nodeTypes     []string
	uid           string // login 后从响应里存下来，用于业务消息鉴权注入
	username      string
	lastForwarded map[string]time.Time
	lastCleanup   time.Time
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

// dialAll 为所有 nodeType 建立 backend WS 连接（可带 token）
// ⚡ 网络 IO 全部在锁外！锁内只有 atomic.Store（ns 级）
func (ps *proxySession) dialAll(token string) {
	// ① 原子读旧指针
	oldPtr := ps.backends.Load()

	// ② 锁外做所有网络操作（dialBackend 可能几百 ms！）
	newMap := make(map[string]*websocket.Conn, len(ps.nodeTypes))
	for _, nt := range ps.nodeTypes {
		routeResp := ps.nodeProvider.Route(rpc.RouteReq{Type: nt})
		if routeResp.Code != 0 {
			continue // 静默跳过（dialAll 全失败状态机会统一打一条 ERROR）
		}
		backendWS, err := dialBackend(routeResp.Addr, token)
		if err != nil {
			continue // 静默跳过（后端刚关 Gateway 还没剔节点，dialAll 全失败时状态机统一处理）
		}
		newMap[nt] = backendWS
		log.Info("["+ps.service+"] ✅ dialBackend success", "type", nt, "addr", routeResp.Addr)
	}

	// ③ 保护：如果一个都 dial 不上，Close 旧 WS + Store empty map
	// ⚠️ 必须 Store empty！否则 RefreshLoop 里 len(*m) > 0 → 跳过 sleep → 无限循环！
	if len(newMap) == 0 {
		// 状态切换才打日志：在线→离线
		if !backendsOffline.Swap(true) {
			log.Error("[Broker] 💥 所有后端离线！dialAll 失败，保持旧连接等待恢复")
		}
		// Close 所有旧 WS（它们大概率已经断了）
		if oldPtr != nil {
			for _, ws := range *oldPtr {
				ws.Close()
			}
		}
		// 关键！Store empty map → RefreshLoop 下一轮检测到 len==0 → sleep + continue ✅
		empty := make(map[string]*websocket.Conn)
		ps.backends.Store(&empty)
		return
	}

	// dialAll 至少成了一个 — 状态切换打日志：离线→在线
	if backendsOffline.Swap(false) {
		log.Info("[Broker] ✅ 后端恢复在线，dialAll 成功")
	}

	// ④ 原子替换 + Close 旧 WS
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
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	dialer := &websocket.Dialer{
		HandshakeTimeout: 3 * time.Second,
	}
	c, _, err := dialer.DialContext(ctx, parsed.String(), nil)
	return c, err
}

// proxyClientToBackend 客户端 → 后端
// 核心改动：转发业务消息前，在 JSON 里注入 __uid/__username（后端 server 收到后自动设 WS Property）
func proxyClientToBackend(ps *proxySession) {
	for {
		_, data, err := ps.client.ReadMessage()
		if err != nil {
			return
		}

		var env struct {
			Name string `json:"name"`
		}
		if json.Unmarshal(data, &env) != nil {
			continue
		}

		backend, ok := ps.GetBackend(routeByName(env.Name))
		if !ok {
			continue
		}

		// 如果已经 login，注入 __uid/__username 到 msg 里
		ps.mu.Lock()
		uid, username := ps.uid, ps.username
		ps.mu.Unlock()
		if uid != "" {
			data = injectAuth(data, uid, username)
			// 🔍 临时诊断：第一次 injectAuth 时打印（限频）
			if env.Name != "account.login" && !diagInjected.Swap(true) {
				log.Info("[Broker] ✅ injectAuth OK", "name", env.Name, "uid", uid)
			}
		} else if env.Name != "account.login" && !diagNoUid.Swap(true) {
			log.Warn("[Broker] ❌ 业务消息但 uid 为空！login 成功了吗？", "name", env.Name)
		}

		// WriteMessage 失败：可能是 dialAll 刚 Close 了旧 WS（安全竞态）
		// 重新 GetBackend 拿最新 WS 重试一次，这条消息大概率能发出去
		if err := backend.WriteMessage(websocket.TextMessage, data); err != nil {
			log.Debug("["+ps.service+"] write failed, retry with new backend", "err", err)
			if newBackend, ok2 := ps.GetBackend(routeByName(env.Name)); ok2 && newBackend != backend {
				if err2 := newBackend.WriteMessage(websocket.TextMessage, data); err2 != nil {
					log.Debug("["+ps.service+"] retry also failed", "err", err2)
				}
			}
		}
	}
}

// injectAuth 在 JSON 的 msg 对象里注入 __uid 和 __username
// 输入: {"name":"chat.world.send","msg":{"content":"hi"}}
// 输出: {"name":"chat.world.send","msg":{"content":"hi","__uid":"xxx","__username":"yyy"}}
func injectAuth(data []byte, uid, username string) []byte {
	var raw map[string]interface{}
	if err := json.Unmarshal(data, &raw); err != nil {
		return data
	}
	msg, ok := raw["msg"].(map[string]interface{})
	if !ok {
		msg = make(map[string]interface{})
		raw["msg"] = msg
	}
	msg["__uid"] = uid
	msg["__username"] = username
	out, err := json.Marshal(raw)
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
		// 简单无限 ReadMessage + recover 防 gorilla 并发 panic
		msgCh := make(chan backendMsg, 2048)
		cancel := make(chan struct{})
		var wg sync.WaitGroup

		for _, e := range entries {
			wg.Add(1)
			go func(name string, ws *websocket.Conn) {
				defer wg.Done()
				// ⚠️ 唯一的保护：recover 捕获 gorilla 并发 Close+Read 的 panic
				// panic 后 return → wg.Wait() 完成 → close(msgCh) → 重连
				defer func() {
					if r := recover(); r != nil {
						// 🔇 限频：每 5s 最多一条 panic 日志
						if now := time.Now(); now.Sub(lastPanicLog) > 5*time.Second {
							lastPanicLog = now
							log.Debug("["+name+"] reader panic (WS closed concurrently)", "err", r)
						}
					}
				}()
				for {
					_, data, err := ws.ReadMessage()
					if err != nil {
						return // EOF / 网络断 / gorilla panic 被 recover 后走到这
					}
					select {
					case msgCh <- backendMsg{data: data, from: name}:
					case <-cancel:
						return
					}
				}
			}(e.name, e.ws)
		}

		go func() { wg.Wait(); close(msgCh) }()

		// ── Step 3: 转发 + 去重 + 存 login uid（做法 A: WS 只是管道）──
		// for range msgCh: channel close 自动退出循环
		for msg := range msgCh {
			// 解析基础信息（name/code/uid/username）
			var tmp struct {
				Name string `json:"name"`
				Code int    `json:"code"`
				Msg  struct {
					UId      string `json:"uid"`
					Username string `json:"username"`
				} `json:"msg"`
			}
			if json.Unmarshal(msg.data, &tmp) != nil {
				continue
			}

			// 20ms 去重
			key := tmp.Name + "|" + string(msg.data)
			ps.mu.Lock()
			if last, exists := ps.lastForwarded[key]; exists && time.Since(last) < 20*time.Millisecond {
				ps.mu.Unlock()
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

			// 转发给客户端
			if err := ps.client.WriteMessage(websocket.TextMessage, msg.data); err != nil {
				return // 客户端断了 → 整个 proxy 结束
			}

			// login 成功 → 记住 uid/username（用于后续 proxyClientToBackend 注入鉴权）
			if tmp.Name == "account.login" && tmp.Code == 0 && tmp.Msg.UId != "" {
				ps.mu.Lock()
				if ps.uid == "" {
					ps.uid = tmp.Msg.UId
					ps.username = tmp.Msg.Username
					// 🔍 临时诊断：login 存 uid 成功（限频打印）
					if !diagLoginStored.Swap(true) {
						log.Info("[Broker] 🎫 login uid 已存储", "uid", tmp.Msg.UId)
					}
					RegisterUID(ps.uid, ps.client)
				}
				ps.mu.Unlock()
			}
		}

		// msgCh close = 所有 backend readers 死了 → 重连
		close(cancel)
		wg.Wait()
		ps.dialAll("") // 不带 token，因为后端 server 鉴权靠消息里的 __uid
		m = ps.backends.Load()
		if m == nil || len(*m) == 0 {
			// 全后端离线：静默退避（状态日志已在 dialAll 里打过了）
			time.Sleep(2 * time.Second)
			continue
		}
	}
}
