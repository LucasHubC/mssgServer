// Package redis — 背包热数据操作
//
// 存储结构:
//   user:inventory:{uid}        → Hash { item_id: count, ... }    全量背包
//   inventory:dirty:uids        → Set<uid>                        修改过、需要落盘的 uid
//
// 设计:
//   - 高频增删查走 Redis（HINCRBY / Lua 原子扣减）
//   - 修改过的 uid 自动标记 dirty，定时任务扫描落 MySQL
//   - 所有函数在 RDB==nil 时安全降级（返回错误，不 panic）

package redis

import (
	"fmt"
	"mssgServer/config/db"
	"mssgServer/infra/log"
	"strconv"
)

// inventoryKey 背包 Redis key
func inventoryKey(uid string) string {
	return "user:inventory:" + uid
}

const dirtyUidsKey = "inventory:dirty:uids"

// MarkDirty 标记某个 uid 的背包已修改，需要后续落盘
func MarkDirty(uid string) {
	if db.RDB == nil {
		return
	}
	db.RDB.SAdd(ctx, dirtyUidsKey, uid)
}

// GetDirtyUids 取出所有脏 uid（定时落盘时调用）
func GetDirtyUids() []string {
	if db.RDB == nil {
		return nil
	}
	uids, err := db.RDB.SMembers(ctx, dirtyUidsKey).Result()
	if err != nil {
		log.Warn("redis GetDirtyUids failed", "err", err)
		return nil
	}
	return uids
}

// ClearDirty 清除某个 uid 的脏标记（落盘成功后调用）
func ClearDirty(uid string) {
	if db.RDB == nil {
		return
	}
	db.RDB.SRem(ctx, dirtyUidsKey, uid)
}

// ────────────────────────────────────────
//  背包数据操作
// ────────────────────────────────────────

// LoadInventory 上线时从 MySQL 全量加载进 Redis
// items: {item_id: count} — 从 DAO 层查出来的
func LoadInventory(uid string, items map[string]int64) error {
	if db.RDB == nil {
		return fmt.Errorf("redis not initialized")
	}
	key := inventoryKey(uid)
	// 先清掉旧的（防止脏数据）
	db.RDB.Del(ctx, key)
	if len(items) == 0 {
		return nil
	}
	pipe := db.RDB.Pipeline()
	for itemID, count := range items {
		pipe.HSet(ctx, key, itemID, count)
	}
	_, err := pipe.Exec(ctx)
	return err
}

// AddItem 增加道具（原子 HINCRBY，数量可为负）
func AddItem(uid, itemID string, delta int64) (int64, error) {
	if db.RDB == nil {
		return 0, fmt.Errorf("redis not initialized")
	}
	key := inventoryKey(uid)
	newCount, err := db.RDB.HIncrBy(ctx, key, itemID, delta).Result()
	if err != nil {
		log.Warn("redis AddItem failed", "uid", uid, "item", itemID, "err", err)
		return 0, err
	}
	MarkDirty(uid)
	return newCount, nil
}

// ConsumeItem 消耗道具（Lua 脚本保证"检查-扣减"原子性）
// 返回 (true, nil) = 成功扣减；(false, nil) = 数量不足；(false, err) = Redis 错误
func ConsumeItem(uid, itemID string, need int64) (bool, error) {
	if db.RDB == nil {
		return false, fmt.Errorf("redis not initialized")
	}
	key := inventoryKey(uid)

	// Lua 脚本：检查 → 扣减 → 返回 1(成功)/0(不足)
	script := `
		local cur = tonumber(redis.call('HGET', KEYS[1], ARGV[1]) or '0')
		if cur < tonumber(ARGV[2]) then
			return 0
		end
		redis.call('HINCRBY', KEYS[1], ARGV[1], -tonumber(ARGV[2]))
		return 1
	`
	result, err := db.RDB.Eval(ctx, script,
		[]string{key}, itemID, strconv.FormatInt(need, 10)).Int64()
	if err != nil {
		log.Warn("redis ConsumeItem failed", "uid", uid, "item", itemID, "err", err)
		return false, err
	}
	if result == 1 {
		MarkDirty(uid)
		return true, nil
	}
	return false, nil // 数量不足
}

// GetItemCount 查询单个道具数量
func GetItemCount(uid, itemID string) (int64, error) {
	if db.RDB == nil {
		return 0, fmt.Errorf("redis not initialized")
	}
	val, err := db.RDB.HGet(ctx, inventoryKey(uid), itemID).Result()
	if err != nil {
		if err.Error() == "redis: nil" {
			return 0, nil // 没有这个道具 = 0
		}
		return 0, err
	}
	return strconv.ParseInt(val, 10, 64)
}

// GetAllItems 全量读取背包（下线保存 / 返回给客户端时调用）
func GetAllItems(uid string) (map[string]int64, error) {
	if db.RDB == nil {
		return nil, fmt.Errorf("redis not initialized")
	}
	raw, err := db.RDB.HGetAll(ctx, inventoryKey(uid)).Result()
	if err != nil {
		return nil, err
	}
	items := make(map[string]int64, len(raw))
	for k, v := range raw {
		n, _ := strconv.ParseInt(v, 10, 64)
		if n > 0 { // 过滤掉 0 或负数（消耗完的道具）
			items[k] = n
		}
	}
	return items, nil
}

// RemoveZeroCount 清理背包里 count<=0 的道具（保存前顺手清）
func RemoveZeroCount(uid string) {
	if db.RDB == nil {
		return
	}
	all, err := db.RDB.HGetAll(ctx, inventoryKey(uid)).Result()
	if err != nil {
		return
	}
	for field, v := range all {
		n, _ := strconv.ParseInt(v, 10, 64)
		if n <= 0 {
			db.RDB.HDel(ctx, inventoryKey(uid), field)
		}
	}
}
