package gateway

import (
	"encoding/json"
	"fmt"
	"mssgServer/infra/log"
	rediscache "mssgServer/infra/redis"
	"mssgServer/infra/rpc"
	"net/http"
	"net/url"
	"strings"
	"sync"
	"sync/atomic"
	"time"

	"github.com/gorilla/websocket"
)

// Gateway 网关注册中心 + 路由分发
type Gateway struct {
	mu     sync.RWMutex
	nodes  map[string][]*Node // type → []node
	policy RoutePolicy
}

func NewGateway(policy RoutePolicy) *Gateway {
	return &Gateway{
		nodes:  make(map[string][]*Node),
		policy: policy,
	}
}

// Register 后端节点注册
func (g *Gateway) Register(req rpc.NodeRegisterReq) error {

	// 参数校验：防止非法 key 污染 map
	if req.Type == "" || req.Addr == "" {
		return fmt.Errorf("invalid register params: type and addr are required")
	}

	g.mu.Lock()
	defer g.mu.Unlock()

	nodes := g.nodes[req.Type]
	// 找：同 addr 已存在 → 原地更新元信息，保留 score 让轮询公平性不丢失
	for _, n := range nodes {
		if n.Addr == req.Addr {
			n.HTTPAddr = req.HTTPAddr
			n.MaxConns = req.MaxConns
			atomic.StoreInt64(&n.Online, 0)
			n.Touch(0) // 重置心跳时间，让 CleanupLoop 不会误判
			log.Info("[Gateway] node re-registered (updated)", "type", req.Type, "addr", req.Addr)
			return nil
		}
	}
	node := &Node{
		Type:          req.Type,
		Addr:          req.Addr,
		HTTPAddr:      req.HTTPAddr,
		MaxConns:      req.MaxConns,
		Online:        0,
		lastHeartbeat: time.Now(), // 注册时算第一次心跳
	}
	g.nodes[req.Type] = append(nodes, node)

	log.Info("[Gateway] node registered (new)", "type", req.Type, "addr", req.Addr)
	return nil
}

// Heartbeat 后端节点上报心跳
func (g *Gateway) Heartbeat(req rpc.NodeHeartbeatReq) error {
	// 全程持读锁遍历，避免 CleanupLoop 并发清理导致 slice 变了
	// Heartbeat 内部只调用 n.Touch（自己有加锁），不会阻塞读
	g.mu.RLock()
	defer g.mu.RUnlock()

	for _, n := range g.nodes[req.Type] {
		if n.Addr == req.Addr {
			n.Touch(req.Online)
			return nil
		}
	}
	return fmt.Errorf("node not found: %s@%s", req.Type, req.Addr)
}

// Route 给客户端分配一个节点
func (g *Gateway) Route(req rpc.RouteReq) *rpc.RouteResp {
	g.mu.RLock()
	nodes := g.nodes[req.Type]
	g.mu.RUnlock()

	// 过滤掉心跳超时的（已死）
	var alive []*Node
	for _, n := range nodes {
		if n.IsAlive() {
			alive = append(alive, n)
		}
	}

	if len(alive) == 0 {
		return &rpc.RouteResp{Code: -1, Msg: "no available node for type: " + req.Type}
	}

	picked := PickNode(alive, g.policy, req.UId)
	return &rpc.RouteResp{Code: 0, Addr: picked.Addr}
}

// CleanupLoop 定期清理超时节点
func (g *Gateway) CleanupLoop() {
	ticker := time.NewTicker(10 * time.Second)
	defer ticker.Stop()
	for range ticker.C {
		g.mu.Lock()
		for typ, nodes := range g.nodes {
			var alive []*Node
			for _, n := range nodes {
				if n.IsAlive() {
					alive = append(alive, n)
				} else {
					log.Warn("[Gateway] node timeout, removing", "type", typ, "addr", n.Addr)
				}
			}
			if len(alive) == 0 {
				delete(g.nodes, typ) // 全部死光 → 彻底删掉这个 key
			} else {
				g.nodes[typ] = alive
			}
		}
		g.mu.Unlock()
	}
}

// Status 调试用：打印所有节点状态
func (g *Gateway) Status() string {
	g.mu.RLock()
	defer g.mu.RUnlock()
	b, _ := json.MarshalIndent(g.nodes, "", "  ")
	return string(b)
}

// ──── HTTP Handler ────

func (g *Gateway) HandleRegister(w http.ResponseWriter, r *http.Request) {
	var req rpc.NodeRegisterReq
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid body", http.StatusBadRequest)
		return
	}
	if err := g.Register(req); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}

func (g *Gateway) HandleHeartbeat(w http.ResponseWriter, r *http.Request) {
	var req rpc.NodeHeartbeatReq
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid body", http.StatusBadRequest)
		return
	}
	if err := g.Heartbeat(req); err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	w.WriteHeader(http.StatusOK)
}

func (g *Gateway) HandleRoute(w http.ResponseWriter, r *http.Request) {
	req := rpc.RouteReq{
		Type: r.URL.Query().Get("type"),
		UId:  r.URL.Query().Get("uid"),
	}
	resp := g.Route(req)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

func (g *Gateway) HandleStatus(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	fmt.Fprint(w, g.Status())
}

// ──── 管理接口（可选，给运营后台用）────

// AdminBroadcastReq 全服广播请求
type AdminBroadcastReq struct {
	Name string      `json:"name"` // 消息名，如 "system.notice"
	Data interface{} `json:"data"` // 消息体
}

// AdminPushReq 跨服推指定 uid
type AdminPushReq struct {
	UID  string      `json:"uid"`
	Name string      `json:"name"`
	Data interface{} `json:"data"`
}

// AdminKickReq 跨服踢人
type AdminKickReq struct {
	UID string `json:"uid"`
}

// HandleAdminBroadcast POST /admin/broadcast
// 通过 Redis Pub/Sub 向所有后端节点广播
func (g *Gateway) HandleAdminBroadcast(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	var req AdminBroadcastReq
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid body", http.StatusBadRequest)
		return
	}
	if req.Name == "" {
		http.Error(w, "name is required", http.StatusBadRequest)
		return
	}

	rediscache.PublishBroadcast(rediscache.BroadcastMsg{
		Action: "broadcast",
		Name:   req.Name,
		Data:   req.Data,
	})

	log.Info("[Gateway] admin broadcast", "name", req.Name)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{"code": 0, "msg": "ok"})
}

// HandleAdminPush POST /admin/push
// 跨服推消息给指定 uid
func (g *Gateway) HandleAdminPush(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	var req AdminPushReq
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid body", http.StatusBadRequest)
		return
	}
	if req.UID == "" || req.Name == "" {
		http.Error(w, "uid and name are required", http.StatusBadRequest)
		return
	}

	rediscache.PublishBroadcast(rediscache.BroadcastMsg{
		Action: "push",
		Target: req.UID,
		Name:   req.Name,
		Data:   req.Data,
	})

	log.Info("[Gateway] admin push", "uid", req.UID, "name", req.Name)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{"code": 0, "msg": "ok"})
}

// HandleAdminKick POST /admin/kick
// 跨服踢某个 uid 下线
func (g *Gateway) HandleAdminKick(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	var req AdminKickReq
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid body", http.StatusBadRequest)
		return
	}
	if req.UID == "" {
		http.Error(w, "uid is required", http.StatusBadRequest)
		return
	}

	rediscache.PublishBroadcast(rediscache.BroadcastMsg{
		Action: "kick",
		Target: req.UID,
	})

	log.Info("[Gateway] admin kick", "uid", req.UID)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{"code": 0, "msg": "ok"})
}

// ═══════════════════════════════════════════════════════════
//  WS 智能路由代理 —— 客户端永远只连 Gateway 一个 WS 连接
//  按消息 name 前缀实时路由到对应后端节点
// ═══════════════════════════════════════════════════════════

var wsProxyUpgrader = websocket.Upgrader{
	ReadBufferSize:  4096,
	WriteBufferSize: 4096,
	CheckOrigin:     func(r *http.Request) bool { return true },
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

// HandleWSProxy 智能 WS 代理入口
// 客户端: ws://gateway:9000/ws?token=xxx
// 所有消息按 name 前缀实时路由到 login/chat/room 后端节点
func (g *Gateway) HandleWSProxy(w http.ResponseWriter, r *http.Request) {
	// ① Upgrade 客户端连接
	clientWS, err := wsProxyUpgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Warn("[Gateway] WS upgrade failed", "err", err)
		return
	}
	atomic.AddInt64(&activeProxyConns, 1)
	defer atomic.AddInt64(&activeProxyConns, -1)
	defer clientWS.Close()

	// ② 收集 query 参数
	query := r.URL.Query()
	token := query.Get("token")

	// ③ 为每种 nodeType 各 Route 一个后端节点（负载均衡！）
	nodeTypes := []string{"login", "chat", "room"}
	backendConns := make(map[string]*websocket.Conn)

	for _, nt := range nodeTypes {
		routeResp := g.Route(rpc.RouteReq{Type: nt})
		if routeResp.Code != 0 {
			log.Warn("[Gateway] no alive node for type", "type", nt)
			continue
		}
		backendWS, err := dialBackend(routeResp.Addr, token)
		if err != nil {
			log.Warn("[Gateway] dial backend failed", "type", nt, "addr", routeResp.Addr, "err", err)
			continue
		}
		backendConns[nt] = backendWS
		defer backendWS.Close()
	}

	if len(backendConns) == 0 {
		log.Warn("[Gateway] no backend available, close client conn")
		return
	}

	log.Info("[Gateway] WS proxy connected", "backendTypes", fmt.Sprintf("%v", mapKeys(backendConns)))

	// ④ 双向转发
	go proxyClientToBackend(clientWS, backendConns)
	proxyBackendToClient(clientWS, backendConns)
}

func mapKeys(m map[string]*websocket.Conn) []string {
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	return keys
}

// dialBackend 连后端 WS 节点，透传 token 参数
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
	c, _, err := websocket.DefaultDialer.Dial(parsed.String(), nil)
	return c, err
}

// proxyClientToBackend 客户端 → Gateway → 后端
func proxyClientToBackend(client *websocket.Conn, backends map[string]*websocket.Conn) {
	for {
		_, data, err := client.ReadMessage()
		if err != nil {
			return // 客户端断开
		}

		// 解析 name 字段
		var env struct {
			Name string `json:"name"`
		}
		if err := json.Unmarshal(data, &env); err != nil {
			continue
		}

		// 按 name 前缀选后端类型
		nodeType := routeByName(env.Name)
		backend, ok := backends[nodeType]
		if !ok {
			continue
		}

		backend.WriteMessage(websocket.TextMessage, data)
	}
}

// proxyBackendToClient 后端 → Gateway → 客户端
func proxyBackendToClient(client *websocket.Conn, backends map[string]*websocket.Conn) {
	// 从所有 backend 并发读，谁先有数据就发给客户端
	type msgWithSource struct {
		data []byte
	}

	channels := make([]<-chan msgWithSource, 0, len(backends))
	for _, b := range backends {
		ch := make(chan msgWithSource, 32)
		go func(ws *websocket.Conn, c chan<- msgWithSource) {
			defer close(c)
			for {
				_, data, err := ws.ReadMessage()
				if err != nil {
					return
				}
				c <- msgWithSource{data: data}
			}
		}(b, ch)
		channels = append(channels, ch)
	}

	// 简单 fan-in：从所有 channel 接收，任一关闭则退出
	for {
		var gotOne bool
		for _, ch := range channels {
			select {
			case msg, ok := <-ch:
				if !ok {
					return // 后端断开，整条链路关
				}
				if err := client.WriteMessage(websocket.TextMessage, msg.data); err != nil {
					return
				}
				gotOne = true
			default:
			}
		}
		if !gotOne {
			// 短暂无数据，避免空转
			select {
			case <-time.After(10 * time.Millisecond):
			}
		}
	}
}
