package gateway

import (
	"encoding/json"
	"fmt"
	"mssgServer/infra/log"
	"mssgServer/infra/rpc"
	"net/http"
	"net/url"
	"sync"
	"sync/atomic"
	"time"
)

// Gateway 网关注册中心 + 路由分发
// 简化原则：Gateway 只做代理转发，不 Subscribe Redis 广播
// admin kick 通过后端 Close Gateway backend WS 传导 → Gateway reader 全退 → return → 客户端断
// admin broadcast 去重：3 后端推同 name 的消息，客户端只收 1 份
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

// ═══════════════════════════════════════════════════════════
//  节点管理核心方法
// ═══════════════════════════════════════════════════════════

// Register 后端节点注册
func (g *Gateway) Register(req rpc.NodeRegisterReq) error {
	if req.Type == "" || req.Addr == "" {
		return fmt.Errorf("invalid register params: type and addr are required")
	}

	g.mu.Lock()
	defer g.mu.Unlock()

	nodes := g.nodes[req.Type]
	for _, n := range nodes {
		if n.Addr == req.Addr {
			n.HTTPAddr = req.HTTPAddr
			n.MaxConns = req.MaxConns
			atomic.StoreInt64(&n.Online, 0)
			n.Touch(0)
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
		lastHeartbeat: time.Now(),
	}
	g.nodes[req.Type] = append(nodes, node)

	log.Info("[Gateway] node registered (new)", "type", req.Type, "addr", req.Addr)
	return nil
}

// Heartbeat 后端节点上报心跳
func (g *Gateway) Heartbeat(req rpc.NodeHeartbeatReq) error {
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

// Deregister 后端节点优雅下线时主动注销
// 比 CleanupLoop 等 10s 更快让 Gateway 摘除这个节点
func (g *Gateway) Deregister(req rpc.NodeHeartbeatReq) {
	g.mu.Lock()
	defer g.mu.Unlock()

	nodes := g.nodes[req.Type]
	for i, n := range nodes {
		if n.Addr == req.Addr {
			// 从切片中移除
			g.nodes[req.Type] = append(nodes[:i], nodes[i+1:]...)
			log.Info("[Gateway] node deregistered", "type", req.Type, "addr", req.Addr)
			if len(g.nodes[req.Type]) == 0 {
				delete(g.nodes, req.Type)
			}
			return
		}
	}
	log.Warn("[Gateway] deregister: node not found", "type", req.Type, "addr", req.Addr)
}

// CleanupLoop 定期清理超时节点
// ⚡ 优化:先 RLock 扫描,再 Lock 更新,减少对 Route 的阻塞
func (g *Gateway) CleanupLoop() {
	ticker := time.NewTicker(10 * time.Second)
	defer ticker.Stop()
	for range ticker.C {
		// Phase 1:RLock 扫描所有节点状态
		g.mu.RLock()
		toDelete := make(map[string][]int) // type → 要删的 index
		for typ, nodes := range g.nodes {
			for i, n := range nodes {
				if !n.IsAlive() {
					toDelete[typ] = append(toDelete[typ], i)
					log.Warn("[Gateway] node timeout", "type", typ, "addr", n.Addr)
				}
			}
		}
		g.mu.RUnlock()

		// Phase 2:有东西要删才拿 Lock
		if len(toDelete) == 0 {
			continue
		}
		g.mu.Lock()
		for typ, idxs := range toDelete {
			nodes := g.nodes[typ]
			if len(nodes) == 0 {
				continue
			}
			var alive []*Node
			idxSet := make(map[int]struct{}, len(idxs))
			for _, i := range idxs {
				idxSet[i] = struct{}{}
			}
			for i, n := range nodes {
				if _, remove := idxSet[i]; !remove {
					alive = append(alive, n)
				}
			}
			if len(alive) == 0 {
				delete(g.nodes, typ)
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

// ═══════════════════════════════════════════════════════════
//  HTTP Handler（节点管理 + 路由分发）
// ═══════════════════════════════════════════════════════════

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

func (g *Gateway) HandleDeregister(w http.ResponseWriter, r *http.Request) {
	var req rpc.NodeHeartbeatReq
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid body", http.StatusBadRequest)
		return
	}
	g.Deregister(req)
	w.WriteHeader(http.StatusOK)
}

// HandleDispatch GET /gateway/dispatch?uid=xxx&token=tk_xxx
// 两阶段连接的第一阶段：根据 uid 选一个 Broker 节点，返回完整 WS URL
// Broker 节点本质就是 type=broker 的 Node，复用 Route 的所有逻辑
func (g *Gateway) HandleDispatch(w http.ResponseWriter, r *http.Request) {
	uid := r.URL.Query().Get("uid")
	token := r.URL.Query().Get("token")

	log.Debug("[Gateway] 📦 dispatch request", "uid", uid, "has_token", token != "", "from", r.RemoteAddr)

	routeResp := g.Route(rpc.RouteReq{Type: "broker", UId: uid})
	if routeResp.Code != 0 {
		log.Warn("[Gateway] ❌ dispatch failed — no broker alive", "uid", uid)
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"code": -1, "msg": "no broker available",
		})
		return
	}

	// routeResp.Addr = broker 注册的 ws://addr
	brokerURL := routeResp.Addr
	if token != "" {
		parsed, err := url.Parse(brokerURL)
		if err == nil {
			q := parsed.Query()
			q.Set("token", token)
			parsed.RawQuery = q.Encode()
			brokerURL = parsed.String()
		}
	}

	log.Info("[Gateway] 📦 dispatch ok", "uid", uid, "→ broker", brokerURL)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"code":       0,
		"broker_url": brokerURL,
	})
}
