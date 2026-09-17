// Package redis 封装了游戏服务端需要的 Redis 操作
// 包括：用户缓存、跨节点在线状态、分布式锁、Pub/Sub 消息广播
//
// 设计原则：
//   - 所有函数在 RDB 未初始化时都是 no-op（安全返回，不阻塞启动）
//   - 业务层只需调用本包方法，不直接依赖 go-redis
package redis

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"mssgServer/config/db"
	"mssgServer/infra/log"
	"time"
)

const redisTimeout = 2 * time.Second

// redisCtx 返回带超时的 context，防止 Redis 网络抖动导致整个游戏阻塞
func redisCtx() context.Context {
	ctx, _ := context.WithTimeout(context.Background(), redisTimeout)
	return ctx
}

var (
	myNodeAddr string // 当前节点的 ws 地址，启动时 SetMyNodeAddr 设一次
)

// ────────────────────────────────────────
//  用户缓存（DAO 层用）
//  key 格式：cache:user:{username} / cache:user:uid:{uid}
// ────────────────────────────────────────

const (
	userCacheTTL = 1 * time.Hour // 用户缓存 1 小时
)

// GetUserCache 从缓存取用户数据
// 返回 (json字节, true) 表示命中，(nil, false) 表示未命中
func GetUserCache(key string) ([]byte, bool) {
	if db.RDB == nil {
		return nil, false
	}
	val, err := db.RDB.Get(redisCtx(), key).Bytes()
	if err != nil {
		return nil, false
	}
	return val, true
}

// SetUserCache 写入用户缓存
func SetUserCache(key string, data []byte) {
	if db.RDB == nil {
		return
	}
	if err := db.RDB.Set(redisCtx(), key, data, userCacheTTL).Err(); err != nil {
		log.Warn("redis SetUserCache failed", "err", err)
	}
}

// DelUserCache 删除缓存（更新用户后调用）
func DelUserCache(keys ...string) {
	if db.RDB == nil {
		return
	}
	db.RDB.Del(redisCtx(), keys...)
}

// UserCacheKeyByUsername 按用户名缓存的 key
func UserCacheKeyByUsername(username string) string {
	return "cache:user:name:" + username
}

// UserCacheKeyByUid 按 uid 缓存的 key
func UserCacheKeyByUid(uid string) string {
	return "cache:user:uid:" + uid
}

// ────────────────────────────────────────
//  跨节点在线状态（ConnManager 用）
//  key 格式：online:user:{uid} → nodeAddr
// ────────────────────────────────────────

const (
	onlineTTL = 30 * time.Second // 在线状态 30s 过期，心跳续期
)

// BindUser 将 uid 绑定到当前节点（登录成功时调用）
// 返回 true 表示绑定成功，false 表示被分布式锁拦住（另一个节点正在处理这个用户）
func BindUser(uid, nodeAddr string) bool {
	if db.RDB == nil {
		return true // 没 Redis 时安全降级：假设绑定成功，单节点场景
	}

	// SETNX 分布式锁：保证同一 uid 不会被两个节点同时绑定
	lockKey := "lock:bind:" + uid
	locked, err := db.RDB.SetNX(redisCtx(), lockKey, nodeAddr, 5*time.Second).Result()
	if err != nil || !locked {
		log.Warn("BindUser failed: another node is binding this uid", "uid", uid)
		return false
	}
	defer safeDelLock(lockKey)

	// 写入在线位置
	key := "online:user:" + uid
	err = db.RDB.Set(redisCtx(), key, nodeAddr, onlineTTL).Err()
	if err != nil {
		log.Warn("BindUser redis Set failed", "err", err)
		return true
	}

	log.Info("[Redis] user bound", "uid", uid, "node", nodeAddr)
	return true
}

// UnbindUser 解除 uid 的在线绑定（下线时调用）
func UnbindUser(uid string) {
	if db.RDB == nil {
		return
	}
	db.RDB.Del(redisCtx(), "online:user:"+uid)
	log.Info("[Redis] user unbound", "uid", uid)
}

// TouchUserHeartbeat 心跳续期（每 10s 调一次，防止在线状态过期）
func TouchUserHeartbeat(uid, nodeAddr string) {
	if db.RDB == nil {
		return
	}
	db.RDB.Set(redisCtx(), "online:user:"+uid, nodeAddr, onlineTTL)
}

// GetUserNode 查某个 uid 现在在哪个节点
// 返回 nodeAddr, true 表示找到了；"", false 表示不在线或 Redis 没开
func GetUserNode(uid string) (string, bool) {
	if db.RDB == nil {
		return "", false
	}
	node, err := db.RDB.Get(redisCtx(), "online:user:"+uid).Result()
	if err != nil {
		return "", false
	}
	return node, true
}

// ────────────────────────────────────────
//  Pub/Sub（跨节点消息广播）
//  频道：node:broadcast  — 所有节点都订阅
// ────────────────────────────────────────

const ChannelBroadcast = "node:broadcast"

// SetMyNodeAddr 设置当前节点地址（main 启动时调一次）
func SetMyNodeAddr(addr string) {
	myNodeAddr = addr
}

// GetMyNodeAddr 返回当前节点地址
func GetMyNodeAddr() string {
	return myNodeAddr
}

// BroadcastMsg 广播消息体
type BroadcastMsg struct {
	Action   string      `json:"action"` // "kick" | "push" | "broadcast"
	UID      string      `json:"uid,omitempty"`
	Name     string      `json:"name,omitempty"`
	Data     interface{} `json:"data,omitempty"`
	FromNode string      `json:"from_node,omitempty"` // 发送节点地址（自动填）
}

// PublishBroadcast 发布跨节点广播消息
// 自动填入 FromNode = 当前节点地址
func PublishBroadcast(msg BroadcastMsg) {
	if db.RDB == nil {
		return
	}
	msg.FromNode = myNodeAddr // 自动标记来源，接收方用来跳过自己发的
	data, err := json.Marshal(msg)
	if err != nil {
		return
	}
	if err := db.RDB.Publish(redisCtx(), ChannelBroadcast, data).Err(); err != nil {
		log.Warn("redis PublishBroadcast failed", "err", err)
	}
}

// SubscribeBroadcast 订阅跨节点广播频道
// 收到消息时调用 handler，阻塞直到 Redis 断开
func SubscribeBroadcast(handler func(BroadcastMsg)) {
	if db.RDB == nil {
		return
	}

	sub := db.RDB.Subscribe(redisCtx(), ChannelBroadcast)
	defer sub.Close()

	ch := sub.Channel()
	for msg := range ch {
		var bm BroadcastMsg
		if err := json.Unmarshal([]byte(msg.Payload), &bm); err != nil {
			continue
		}
		handler(bm)
	}
}

// ────────────────────────────────────────
//  分布式锁（通用）
// ────────────────────────────────────────

// TryLock 尝试获取分布式锁，key 是锁名，ttl 是过期时间
// 返回 (unlock函数, true) 表示拿到锁了；(nil, false) 表示被别人占着
func TryLock(key string, ttl time.Duration) (func(), bool) {
	if db.RDB == nil {
		return func() {}, true // 安全降级
	}

	ok, err := db.RDB.SetNX(redisCtx(), key, "1", ttl).Result()
	if err != nil || !ok {
		return nil, false
	}

	unlock := func() {
		safeDelLock(key)
	}
	return unlock, true
}

// safeDelLock 安全释放分布式锁：Lua 脚本检查锁还存在才 DEL
// 防止锁 TTL 过期后被其他节点重新获取，本节点误删别人的锁
func safeDelLock(lockKey string) {
	if db.RDB == nil {
		return
	}
	// 只用 TTL 判断：TTL > 0 说明锁还没过期（就是我们的），可以删
	// 如果 TTL == -1（无过期）或 TTL == -2（不存在），跳过
	const script = `
		local ttl = redis.call('TTL', KEYS[1])
		if ttl > 0 then
			redis.call('DEL', KEYS[1])
		end
	`
	db.RDB.Eval(redisCtx(), script, []string{lockKey})
}

// ────────────────────────────────────────
//  群组管理（Redis Set）
//  key 格式: chat:group:{groupId}  → Set<uid>
// ────────────────────────────────────────

// groupKey 群组 Redis key
func groupKey(groupId string) string {
	return "chat:group:" + groupId
}

// JoinGroup 加入群组
func JoinGroup(groupId, uid string) {
	if db.RDB == nil {
		return
	}
	if err := db.RDB.SAdd(redisCtx(), groupKey(groupId), uid).Err(); err != nil {
		log.Warn("redis JoinGroup failed", "group", groupId, "uid", uid, "err", err)
	}
}

// LeaveGroup 退出群组
func LeaveGroup(groupId, uid string) {
	if db.RDB == nil {
		return
	}
	if err := db.RDB.SRem(redisCtx(), groupKey(groupId), uid).Err(); err != nil {
		log.Warn("redis LeaveGroup failed", "group", groupId, "uid", uid, "err", err)
	}
}

// GetGroupMembers 获取群组所有成员 uid（跨服都能查到，因为 Redis 是共享的）
func GetGroupMembers(groupId string) []string {
	if db.RDB == nil {
		return nil
	}
	members, err := db.RDB.SMembers(redisCtx(), groupKey(groupId)).Result()
	if err != nil {
		log.Warn("redis GetGroupMembers failed", "group", groupId, "err", err)
		return nil
	}
	return members
}

// ────────────────────────────────────────
//  Session / Token 存储（HTTP 短连接用）
//
//  设计：
//   - 登录成功时生成 token = "tk_" + 32位hex
//   - Redis key: sess:{token} → {uid, username} (JSON)
//   - 默认过期 2 小时，防止 token 无限堆积
//   - 跨服务可用：apiLoginServer 签发的 token，wsLoginServer 也能查
// ────────────────────────────────────────

const (
	sessionTTL = 2 * time.Hour
)

// Session 会话信息（存在 Redis 里的内容）
type Session struct {
	Token    string `json:"-"` // 原始 token（Redis key，不存进 value）
	UId      string `json:"uid"`
	Username string `json:"username"`
}

// StoreSession 签发 token，存到 Redis
// 返回生成的 token（形如 "tk_a1b2c3..."）
func StoreSession(sess Session) (token string, err error) {
	// 生成 32 位随机 hex + 前缀
	b := make([]byte, 16)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	token = "tk_" + hex.EncodeToString(b)

	if db.RDB == nil {
		// 没 Redis 时安全降级：token 照样生成，但不持久化
		// 这种场景下 Token 校验会直接失败（单节点无 Redis 不支持 HTTP Token）
		log.Warn("StoreSession: Redis not initialized, token will NOT be valid")
		return token, nil
	}

	data, _ := json.Marshal(sess)
	key := "sess:" + token
	if err := db.RDB.Set(redisCtx(), key, data, sessionTTL).Err(); err != nil {
		log.Warn("StoreSession redis Set failed", "err", err)
		return "", err
	}
	return token, nil
}

// GetSession 用 token 查会话
// 返回 (Session, true) 表示找到；(nil, false) 表示 token 无效/过期
// 返回的 Session.Token 会被填上原始 token（方便调用方登出时用）
func GetSession(token string) (Session, bool) {
	if db.RDB == nil || token == "" {
		return Session{}, false
	}
	data, err := db.RDB.Get(redisCtx(), "sess:"+token).Bytes()
	if err != nil {
		return Session{}, false
	}
	var sess Session
	if err := json.Unmarshal(data, &sess); err != nil {
		return Session{}, false
	}
	sess.Token = token // ✅ 把原始 token 填回去
	return sess, true
}

// DelSession 注销 token（登出时调用）
func DelSession(token string) {
	if db.RDB == nil || token == "" {
		return
	}
	db.RDB.Del(redisCtx(), "sess:"+token)
}

// IsGroupMember 判断 uid 是否在群组里
func IsGroupMember(groupId, uid string) bool {
	if db.RDB == nil {
		return false
	}
	ok, err := db.RDB.SIsMember(redisCtx(), groupKey(groupId), uid).Result()
	if err != nil {
		return false
	}
	return ok
}
