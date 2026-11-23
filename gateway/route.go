package gateway

import (
	"crypto/md5"
	"encoding/binary"
	"sort"
	"sync/atomic"
)

// roundRobinCounter 全局轮询计数器
// 每次 round_robin 路由时 +1，取模得到节点下标
var roundRobinCounter int64

// RoutePolicy 路由策略枚举
type RoutePolicy string

const (
	PolicyRoundRobin RoutePolicy = "round_robin" // 轮询（默认）
	PolicyLeastConn  RoutePolicy = "least_conn"  // 最少连接
	PolicyUIDHash    RoutePolicy = "uid_hash"    // 按 uid 哈希（同一用户落到同一节点）
)

// PickNode 从候选节点列表中选一个
// ⚠️ 内部会复制一份再排序，不修改传入的原始切片
func PickNode(nodes []*Node, policy RoutePolicy, uid string) *Node {
	if len(nodes) == 0 {
		return nil
	}

	// 复制一份再 sort，避免意外修改原始切片（Gateway.Route 传进来的 alive 虽然是新建的，
	// 但养成好习惯，防止以后有人直接把 g.nodes[type] 传进来）
	candidates := make([]*Node, len(nodes))
	copy(candidates, nodes)

	switch policy {
	case PolicyLeastConn:
		// 在线人数最少的优先
		sort.Slice(candidates, func(i, j int) bool {
			return candidates[i].Online < candidates[j].Online
		})
		return candidates[0]

	case PolicyUIDHash:
		// 对 uid 做 md5 → 取前 2 字节当 uint16 → 模节点数
		sum := md5.Sum([]byte(uid))
		idx := int(binary.BigEndian.Uint16(sum[:2])) % len(candidates)
		return candidates[idx]

	default: // PolicyRoundRobin
		// 原子计数器 + 取模，简单确定，不依赖 sort 稳定性
		idx := int(atomic.AddInt64(&roundRobinCounter, 1)-1) % len(candidates)
		return candidates[idx]
	}
}
