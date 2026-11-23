package net

import (
	"mssgServer/config"
	"mssgServer/infra/log"
	"mssgServer/infra/metrics"
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
	broadcastPoolOnce.Do(initBroadcastPool) // 启动常驻 worker pool
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
	metrics.Connections.Set(float64(len(cm.conns)))
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
	metrics.Connections.Set(float64(len(cm.conns)))
}

// BindUid 登录成功后绑定 uid（同一 uid 重复登录会踢掉本机旧连接）
// skipRedis: true = 只写本机 uidMap（做法 A 后端代理场景，Redis 锁由 login 服务独占）
//
//	false = 完整流程（login handler 用，需要 Redis 锁保证唯一性）
func (cm *ConnectionManager) BindUid(uid string, ws *WSServer, skipRedis bool) {
	bindUidDiagCount.Add(1)
	// ① 写 Redis 在线位置（login handler 才需要，后端代理跳过）
	if !skipRedis && ws.NodeAddr != "" {
		if !rediscache.BindUser(uid, ws.NodeAddr) {
			log.Warn("BindUid rejected by redis lock", "uid", uid)
			return
		}
	}

	// ② 本机处理：踢本机旧连接（同一进程内同一 uid 的旧 WS）
	cm.mu.Lock()
	if old, ok := cm.uidMap[uid]; ok && old != ws {
		delete(cm.uidMap, uid)
		delete(cm.conns, old)
		go old.Close()
	}
	cm.uidMap[uid] = ws
	metrics.UIDMapSize.Set(float64(len(cm.uidMap)))
	metrics.Connections.Set(float64(len(cm.conns))) // ✅ 踢旧后立刻同步
	cm.mu.Unlock()

	// ⚠️ 不再做跨节点 kick！
	// Gateway 场景：login/chat/room 后端同时 BindUid 同一个 uid 是正常代理行为，
	// 每个后端都以为"用户从别的节点迁过来"→ 都发 kick → Gateway 收到→ 踢掉客户端
	// 跨服踢人统一链路：admin/kick → Redis kick → 所有后端 Subscribe kick
	//   → 后端 Close Gateway 的 backend WS → Gateway reader 全退 → return → 客户端断 ✅
}

// GetByUid 按 uid 查连接
func (cm *ConnectionManager) GetByUid(uid string) *WSServer {
	cm.mu.RLock()
	defer cm.mu.RUnlock()
	return cm.uidMap[uid]
}

// ── 常驻广播 Worker Pool ─────────────────────────────────────────
// 预创建 16 个 worker goroutine 循环从 channel 取任务
// 避免每条 world 消息创建 16 个 goroutine（高并发下 goroutine 爆炸）

type broadcastTask struct {
	name    string
	data    interface{}
	targets []*WSServer // 预分片好的目标列表
}

const broadcastWorkerCount = 16
const broadcastTaskBuffer = 8192 // 16分片 × 512 缓冲/分片，扛 world 广播洪峰

var broadcastTaskCh chan broadcastTask
var broadcastPoolOnce sync.Once
var broadcastDropCount atomic.Int64

func initBroadcastPool() {
	broadcastTaskCh = make(chan broadcastTask, broadcastTaskBuffer)
	for i := 0; i < broadcastWorkerCount; i++ {
		go runBroadcastWorker(broadcastTaskCh)
	}
}

func runBroadcastWorker(ch <-chan broadcastTask) {
	for task := range ch {
		for _, ws := range task.targets {
			ws.Push(task.name, task.data)
		}
	}
}

// broadcastDiagCount 诊断用：统计 Broadcast 调用次数
var broadcastDiagCount atomic.Int64
var bindUidDiagCount atomic.Int64

// Broadcast 向本机所有已登录用户推消息
// ⚡ 业界主流优化：异步 + 分批
//   - 异步：不阻塞调用方（handler 立刻返回处理下一条）
//   - 分批：uidMap 分成 16 份，每份一个 goroutine 并行推
//     1000 客户端 × 10 条 world/s = 10000 条 world/s
//     每条 × 1000 用户 = 1000 万次 push/s
//     分 16 份并行 → 每份 62.5 万次 push/s → 可控
func (cm *ConnectionManager) Broadcast(name string, data interface{}) {
	cm.mu.RLock()
	uids := make([]*WSServer, 0, len(cm.uidMap))
	for _, ws := range cm.uidMap {
		uids = append(uids, ws)
	}
	cm.mu.RUnlock()

	// 🔍 诊断：定期打印 uidMap 大小（第一次 + 每 5000 次）
	broadcastDiagCount.Add(1)
	if broadcastDiagCount.Load() <= 3 || broadcastDiagCount.Load()%5000 == 0 {
		log.Debug("[Broadcast] diag", "uidMap_size", len(uids), "total_conns", cm.Total(), "bindUid_count", bindUidDiagCount.Load())
	}

	if len(uids) == 0 {
		return
	}

	// 分片 + 投递到常驻 worker pool（零 goroutine 创建）
	chunk := (len(uids) + broadcastWorkerCount - 1) / broadcastWorkerCount
	dropped := 0
	for i := 0; i < broadcastWorkerCount; i++ {
		start := i * chunk
		if start >= len(uids) {
			break
		}
		end := start + chunk
		if end > len(uids) {
			end = len(uids)
		}
		task := broadcastTask{name: name, data: data, targets: uids[start:end]}
		select {
		case broadcastTaskCh <- task:
		default:
			// channel 满了 → 丢弃分片（广播丢包用户可接受）
			dropped++
		}
	}
	if dropped > 0 {
		if broadcastDropCount.Add(int64(dropped))%1000 == 0 {
			log.Warn("[Broadcast] task channel full, dropped shards", "dropped", broadcastDropCount.Load(), "name", name)
		}
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
	for {
		func() {
			defer func() {
				if r := recover(); r != nil {
					log.Error("[HeartbeatLoop] PANIC", "err", r)
				}
			}()
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

				nodeAddr := rediscache.GetMyNodeAddr()
				rediscache.TouchUserHeartbeatBatch(uids, nodeAddr)
				log.Debug("[Heartbeat] online renewed", "count", len(uids))
			}
		}()
		time.Sleep(5 * time.Second) // panic 后冷却再重启
	}
}

// nodeAddr 获取当前节点地址，优先用 redis 包的全局值
// （SetMyNodeAddr 在 main 启动时已设好，redis.GetMyNodeAddr() 永远返回非空）
func (cm *ConnectionManager) getNodeAddr() string {
	return rediscache.GetMyNodeAddr()
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
	for {
		func() {
			defer func() {
				if r := recover(); r != nil {
					log.Error("[SubscribeRemoteBroadcast] PANIC", "err", r)
				}
			}()
			rediscache.SubscribeBroadcast(func(msg rediscache.BroadcastMsg) {
				// 🟢 跳过自己发的 —— 避免 BroadcastAll 本机被推两遍
				// （PublishBroadcast 会自动填 FromNode = myNodeAddr）
				if msg.FromNode != "" && msg.FromNode == rediscache.GetMyNodeAddr() {
					return
				}

				switch msg.Action {
				case "kick":
					// kick 统一链路：admin/kick → Redis kick → 所有后端 Subscribe kick
					//   → 后端 Close Gateway 的 backend WS → Gateway reader 全退
					//   → proxyBackendToClient return → HandleWSProxy defer Close client WS
					cm.mu.Lock()
					if ws, ok := cm.uidMap[msg.Target]; ok {
						log.Info("[Redis] remote kick: closing WS", "uid", msg.Target, "from", msg.FromNode)
						delete(cm.uidMap, msg.Target)
						delete(cm.conns, ws)
						cm.mu.Unlock()
						rediscache.UnbindUser(msg.Target)
						go ws.Close() // Close Gateway backend WS → Gateway reader 全退
					} else {
						cm.mu.Unlock()
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
		}()
		time.Sleep(5 * time.Second) // panic 后冷却再重启
	}
}
