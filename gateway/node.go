package gateway

import (
	"sync"
	"sync/atomic"
	"time"
)

// Node 一个后端节点的元信息
type Node struct {
	Type     string `json:"type"`      // login / game / chat
	Addr     string `json:"addr"`      // ws://127.0.0.1:8001
	HTTPAddr string `json:"http_addr"` // http://127.0.0.1:9001 (健康检查用)
	MaxConns int64  `json:"max_conns"`
	Online   int64  `json:"online"` // 当前在线

	// 内部字段
	mu            sync.RWMutex
	lastHeartbeat time.Time // 最后一次心跳时间
	score         int64     // 路由评分（越小越优先），由路由策略更新
}

// 心跳超时阈值：超过这个时间没收到心跳就认为节点死了
const nodeTimeout = 15 * time.Second

// IsAlive 判断节点是否还活着（心跳超时不算）
func (n *Node) IsAlive() bool {
	n.mu.RLock()
	defer n.mu.RUnlock()
	return time.Since(n.lastHeartbeat) < nodeTimeout
}

// Touch 收到心跳时更新时间 + 在线人数
func (n *Node) Touch(online int64) {
	n.mu.Lock()
	defer n.mu.Unlock()
	n.lastHeartbeat = time.Now()
	atomic.StoreInt64(&n.Online, online)
}

// IncrScore 每次被路由选中 +1（轮询用）
func (n *Node) IncrScore() int64 {
	return atomic.AddInt64(&n.score, 1)
}

// Score 当前评分
func (n *Node) Score() int64 {
	return atomic.LoadInt64(&n.score)
}
