package net

import (
	"sync"
	"sync/atomic"
)

// ConnectionManager 全局连接管理器
// 职责：维护 uid → *WSServer 的映射，支持推送、踢下线、广播
type ConnectionManager struct {
	mu       sync.RWMutex
	conns    map[*WSServer]struct{} // 所有活跃连接
	uidMap   map[string]*WSServer   // uid → 连接（快速寻址）
	total    int64                  // 当前总连接数
	maxConns int64                  // 最大允许连接数
}

// ConnManager 全局单例，默认最多 20000 连接
var ConnManager = NewConnectionManager(20000)

func NewConnectionManager(maxConns int64) *ConnectionManager {
	return &ConnectionManager{
		conns:    make(map[*WSServer]struct{}),
		uidMap:   make(map[string]*WSServer),
		maxConns: maxConns,
	}
}

// MaxConns 当前上限
func (cm *ConnectionManager) MaxConns() int64 {
	cm.mu.RLock()
	defer cm.mu.RUnlock()
	return cm.maxConns
}

// Register 注册新连接，返回 false 表示连接数已满
func (cm *ConnectionManager) Register(ws *WSServer) bool {
	cm.mu.Lock()
	defer cm.mu.Unlock()
	if int64(len(cm.conns)) >= cm.maxConns {
		return false
	}
	cm.conns[ws] = struct{}{}
	atomic.AddInt64(&cm.total, 1)
	return true
}

// Unregister 注销连接，同时解绑 uid
func (cm *ConnectionManager) Unregister(ws *WSServer) {
	cm.mu.Lock()
	defer cm.mu.Unlock()
	delete(cm.conns, ws)
	uid, _ := ws.GetProperty("uid")
	if uid != nil {
		if cur, ok := cm.uidMap[uid.(string)]; ok && cur == ws {
			delete(cm.uidMap, uid.(string))
		}
	}
	atomic.AddInt64(&cm.total, -1)
}

// BindUid 登录成功后绑定 uid（同一 uid 重复登录会踢掉旧连接）
func (cm *ConnectionManager) BindUid(uid string, ws *WSServer) {
	cm.mu.Lock()
	defer cm.mu.Unlock()
	if old, ok := cm.uidMap[uid]; ok && old != ws {
		delete(cm.uidMap, uid)
		delete(cm.conns, old)
		go old.Close()
	}
	cm.uidMap[uid] = ws
}

// GetByUid 按 uid 查连接
func (cm *ConnectionManager) GetByUid(uid string) *WSServer {
	cm.mu.RLock()
	defer cm.mu.RUnlock()
	return cm.uidMap[uid]
}

// Broadcast 向所有活跃连接推消息（非阻塞，满了跳过）
func (cm *ConnectionManager) Broadcast(name string, data interface{}) {
	cm.mu.RLock()
	defer cm.mu.RUnlock()
	for ws := range cm.conns {
		select {
		case ws.outChan <- &WsMsgResp{
			Body: &RespBody{Seq: ws.Seq, Name: name, Code: 0, Msg: data},
		}:
		default:
		}
	}
}

// KickByUid 踢下线某个 uid
func (cm *ConnectionManager) KickByUid(uid string) {
	if ws := cm.GetByUid(uid); ws != nil {
		go ws.Close()
	}
}

// Total 当前在线连接数
func (cm *ConnectionManager) Total() int64 {
	return atomic.LoadInt64(&cm.total)
}

// SetMaxConns 动态调整连接数上限
func (cm *ConnectionManager) SetMaxConns(max int64) {
	cm.mu.Lock()
	defer cm.mu.Unlock()
	cm.maxConns = max
}

// All 当前所有活跃连接（调试用）
func (cm *ConnectionManager) All() []*WSServer {
	cm.mu.RLock()
	defer cm.mu.RUnlock()
	list := make([]*WSServer, 0, len(cm.conns))
	for ws := range cm.conns {
		list = append(list, ws)
	}
	return list
}
