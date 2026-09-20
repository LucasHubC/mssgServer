package net

import (
	"mssgServer/config"
	"mssgServer/infra/log"
	rediscache "mssgServer/infra/redis"
	"sync"
	"sync/atomic"
	"time"
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

// ConnManager 全局单例，最大连接数从 config.ini [ws_server] max_conns 读取
// 默认 20000，阿里云 8核16G 可提至 50000
var ConnManager = NewConnectionManager(int64(config.File.MustInt("ws_server", "max_conns", 20000)))

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

// Unregister 注销连接，同时解绑 uid + 清理 Redis 在线状态
// 整个函数只做一次 Lock/Unlock，避免 unlock→IO→lock 之间状态被篡改
func (cm *ConnectionManager) Unregister(ws *WSServer) {
	var uidToUnbind string
	cm.mu.Lock()
	delete(cm.conns, ws)
	uidVal, _ := ws.GetProperty("uid")
	if uidStr, ok := uidVal.(string); ok && uidStr != "" {
		if cur, exists := cm.uidMap[uidStr]; exists && cur == ws {
			delete(cm.uidMap, uidStr)
			uidToUnbind = uidStr
		}
	}
	cm.mu.Unlock()

	// Redis IO 移到锁外，避免阻塞 cm.mu
	if uidToUnbind != "" {
		rediscache.UnbindUser(uidToUnbind)
	}
	atomic.AddInt64(&cm.total, -1)
}

// BindUid 登录成功后绑定 uid（同一 uid 重复登录会踢掉旧连接）
// 同时通过 Redis 告知其他节点"这个用户在我这了"
func (cm *ConnectionManager) BindUid(uid string, ws *WSServer) {
	// ① 先查 Redis 里这个 uid 原来在哪个节点（用来通知旧节点踢人）
	var oldNode string
	if ws.NodeAddr != "" {
		if node, ok := rediscache.GetUserNode(uid); ok {
			oldNode = node
		}
	}

	// ② 拿分布式锁 + 写入新的在线位置
	if ws.NodeAddr != "" {
		if !rediscache.BindUser(uid, ws.NodeAddr) {
			log.Warn("BindUid rejected by redis lock", "uid", uid)
			return
		}
	}

	// ③ 本机处理：踢本机旧连接
	cm.mu.Lock()
	if old, ok := cm.uidMap[uid]; ok && old != ws {
		delete(cm.uidMap, uid)
		delete(cm.conns, old)
		go old.Close()
	}
	cm.uidMap[uid] = ws
	cm.mu.Unlock()

	// ④ 跨节点处理：如果旧连接在别的节点，广播让它们踢
	// （BindUser 内部已经 SET 了新值，旧节点的 SubscribeRemoteBroadcast
	//   会收到这个广播 → 踢掉本机上的那个 uid 对应的连接）
	if oldNode != "" && oldNode != ws.NodeAddr {
		rediscache.PublishBroadcast(rediscache.BroadcastMsg{
			Action: "kick",
			Target: uid,
		})
		log.Info("[Redis] broadcast remote kick for", "uid", uid, "oldNode", oldNode)
	}
}

// GetByUid 按 uid 查连接
func (cm *ConnectionManager) GetByUid(uid string) *WSServer {
	cm.mu.RLock()
	defer cm.mu.RUnlock()
	return cm.uidMap[uid]
}

// Broadcast 向本机所有活跃连接推消息（非阻塞）
func (cm *ConnectionManager) Broadcast(name string, data interface{}) {
	cm.mu.RLock()
	defer cm.mu.RUnlock()
	for ws := range cm.conns {
		ws.Push(name, data) // 统一走 Push：atomic Seq + safeSend 非阻塞
	}
}

// BroadcastAll 全服广播（本机 + 通过 Redis 广播给所有其他节点）
func (cm *ConnectionManager) BroadcastAll(name string, data interface{}) {
	cm.Broadcast(name, data) // 先推本机
	rediscache.PublishBroadcast(rediscache.BroadcastMsg{
		Action: "broadcast",
		Name:   name,
		Data:   data,
	})
}

// PushToUid （可跨服私信） 推消息给指定 uid（本机直推，不在本机则通过 Redis 广播让持有该 uid 的节点推）
func (cm *ConnectionManager) PushToUid(uid, name string, data interface{}) {
	// ① 本机有 → 直接推
	if ws := cm.GetByUid(uid); ws != nil {
		ws.Push(name, data)
		return
	}
	// ② 本机没有 → 广播给所有节点，持有该 uid 的节点会处理
	rediscache.PublishBroadcast(rediscache.BroadcastMsg{
		Action: "push",
		Target: uid,
		Name:   name,
		Data:   data,
	})
}

// KickByUid 踢下线某个 uid（本机）
// 如果 uid 不在本机，会通过 Redis Pub/Sub 广播给其他节点
func (cm *ConnectionManager) KickByUid(uid string) {
	if ws := cm.GetByUid(uid); ws != nil {
		go ws.Close()
		return
	}
	// 不在本机 → 广播给其他节点
	rediscache.PublishBroadcast(rediscache.BroadcastMsg{
		Action: "kick",
		Target: uid,
	})
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

// PushToLocalGroupMembers 推给本节点上属于该群组的所有成员
// （BroadcastAll 模式：先本机推，再 PublishBroadcast 给其他节点）
func (cm *ConnectionManager) PushToLocalGroupMembers(groupId, name string, data interface{}) {
	cm.pushToLocalUids(rediscache.GetGroupMembers(groupId), name, data)
}

// PushToLocalRoomMembers 推给本节点上属于该房间的所有成员
func (cm *ConnectionManager) PushToLocalRoomMembers(roomId, name string, data interface{}) {
	cm.pushToLocalUids(rediscache.GetRoomMembers(roomId), name, data)
}

// pushToLocalUids 通用：推给本机上指定的一组 uid（group/room 广播共用）
func (cm *ConnectionManager) pushToLocalUids(uids []string, name string, data interface{}) {
	cm.mu.RLock()
	for _, uid := range uids {
		if ws := cm.uidMap[uid]; ws != nil {
			ws.Push(name, data)
		}
	}
	cm.mu.RUnlock()
}

// StartHeartbeatLoop 定时续期所有本机在线用户的 Redis 在线状态 key
//
// 为什么需要它？
//
//	Redis 里 online:user:{uid} → nodeAddr 有 30s TTL（onlineTTL）
//	用户登录时 BindUser 写一次，之后如果什么都不做，30s 后 key 就过期了
//	过期后其他节点 GetUserNode 返回 false → 私信/踢人都找不到这个用户
//
//	所以这里每 10s Touch 一次（远小于 30s TTL），只要用户还活着就不会过期
//	TouchUserHeartbeat 内部就是重新 SET online:user:{uid} → nodeAddr，TTL 重置
//
// 启动位置：loginServer.go / chatServer.go 的 main() 里各启一个
func (cm *ConnectionManager) StartHeartbeatLoop() {
	ticker := time.NewTicker(10 * time.Second)
	defer ticker.Stop()

	for range ticker.C {
		cm.mu.RLock()
		uids := make([]string, 0, len(cm.uidMap))
		for uid := range cm.uidMap {
			uids = append(uids, uid)
		}
		cm.mu.RUnlock()

		if len(uids) == 0 {
			continue
		}

		// 逐个 Touch —— Redis 网络 IO 放到锁外做，不阻塞 cm.mu
		// 虽然可以 Pipeline 批量，但 uid 数量一般不大（单机几千），逐个写够了
		for _, uid := range uids {
			rediscache.TouchUserHeartbeat(uid, cm.getNodeAddr())
		}
		log.Debug("[Heartbeat] online renewed", "count", len(uids))
	}
}

// nodeAddr 当前节点的 ws 地址（启动时 SetNodeAddr 设的）
// 用来告诉 TouchUserHeartbeat 续期时写哪个 nodeAddr
func (cm *ConnectionManager) getNodeAddr() string {
	// ConnectionManager 没有存 NodeAddr，但 ConnManager 是单例，
	// 可以从 WSServer.NodeAddr 拿到（所有 WS 连接共用同一个 NodeAddr）
	cm.mu.RLock()
	defer cm.mu.RUnlock()
	for _, ws := range cm.uidMap {
		return ws.NodeAddr // 任取一个就行，同一节点的 NodeAddr 都一样
	}
	return ""
}

// OnlineUids 当前已登录用户的 uid 列表
// 比 All() 后遍历 GetProperty 高效得多（直接读 uidMap）
func (cm *ConnectionManager) OnlineUids() []string {
	cm.mu.RLock()
	defer cm.mu.RUnlock()
	uids := make([]string, 0, len(cm.uidMap))
	for uid := range cm.uidMap {
		uids = append(uids, uid)
	}
	return uids
}

// SubscribeRemoteBroadcast 后台协程：订阅 Redis 广播频道，处理其他节点发来的消息
// 支持的 Action:
//   - "kick":     其他节点请求踢某个 uid 下线
//   - "push":     其他节点请求给某个 uid 推消息
//   - "broadcast":其他节点请求全服广播（本机部分）
func (cm *ConnectionManager) SubscribeRemoteBroadcast() {
	rediscache.SubscribeBroadcast(func(msg rediscache.BroadcastMsg) {
		// 🟢 跳过自己发的 —— 避免 BroadcastAll 本机被推两遍
		// （PublishBroadcast 会自动填 FromNode = myNodeAddr）
		if msg.FromNode != "" && msg.FromNode == rediscache.GetMyNodeAddr() {
			return
		}

		switch msg.Action {
		case "kick":
			if ws := cm.GetByUid(msg.Target); ws != nil {
				log.Debug("[Redis] remote kick received", "uid", msg.Target, "from", msg.FromNode)
				go ws.Close()
			}

		case "push":
			if ws := cm.GetByUid(msg.Target); ws != nil {
				log.Debug("[Redis] remote push received", "uid", msg.Target, "name", msg.Name, "from", msg.FromNode)
				ws.Push(msg.Name, msg.Data)
			}

		case "broadcast":
			log.Debug("[Redis] remote broadcast received", "name", msg.Name, "from", msg.FromNode)
			cm.Broadcast(msg.Name, msg.Data)

		case "group":
			// Action="group" → Target 就是 groupId
			groupIds := rediscache.GetGroupMembers(msg.Target)
			log.Debug("[Redis] remote group msg", "group", msg.Target, "members", len(groupIds), "from", msg.FromNode)
			cm.pushToLocalUids(groupIds, msg.Name, msg.Data)

		case "room":
			// Action="room" → Target 就是 roomId
			roomMembers := rediscache.GetRoomMembers(msg.Target)
			log.Debug("[Redis] remote room msg", "room", msg.Target, "members", len(roomMembers), "from", msg.FromNode)
			cm.pushToLocalUids(roomMembers, msg.Name, msg.Data)

		default:
			log.Warn("[Redis] unknown broadcast action", "action", msg.Action)
		}
	})
}
