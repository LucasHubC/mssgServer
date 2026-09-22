package gateway

import (
	"encoding/json"
	"mssgServer/infra/log"
	rediscache "mssgServer/infra/redis"
	"net/http"
)

// ═══════════════════════════════════════════════════════════
//  管理接口（给运营后台用）
//  全部通过 Redis Pub/Sub 跨服通信，Gateway 自身不 Subscribe
// ═══════════════════════════════════════════════════════════

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