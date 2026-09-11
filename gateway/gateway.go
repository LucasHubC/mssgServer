package gateway

import (
	"encoding/json"
	"fmt"
	"mssgServer/infra/log"
	rediscache "mssgServer/infra/redis"
	"mssgServer/infra/rpc"
	"net/http"
	"sync"
	"sync/atomic"
	"time"
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
		UID:    req.UID,
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
		UID:    req.UID,
	})

	log.Info("[Gateway] admin kick", "uid", req.UID)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{"code": 0, "msg": "ok"})
}
