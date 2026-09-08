package gateway

import (
	"encoding/json"
	"fmt"
	"mssgServer/infra/log"
	"mssgServer/infra/rpc"
	"net/http"
	"sync"
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
	g.mu.Lock()
	defer g.mu.Unlock()

	// 同 addr 重复注册 → 更新（先删掉旧的再加）
	nodes := g.nodes[req.Type]
	for i, n := range nodes {
		if n.Addr == req.Addr {
			nodes = append(nodes[:i], nodes[i+1:]...)
			break
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

	log.Info("[Gateway] node registered", "type", req.Type, "addr", req.Addr)
	return nil
}

// Heartbeat 后端节点上报心跳
func (g *Gateway) Heartbeat(req rpc.NodeHeartbeatReq) error {
	g.mu.RLock()
	nodes := g.nodes[req.Type]
	g.mu.RUnlock()

	for _, n := range nodes {
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
