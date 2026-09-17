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
	"errors"
	"fmt"
	"mssgServer/config/db"
	"mssgServer/infra/log"
	"strconv"

	"github.com/redis/go-redis/v9"
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
	db.RDB.SAdd(redisCtx(), dirtyUidsKey, uid)
}

// GetDirtyUids 取出所有脏 uid（定时落盘时调用）
func GetDirtyUids() []string {
	if db.RDB == nil {
		return nil
	}
	uids, err := db.RDB.SMembers(redisCtx(), dirtyUidsKey).Result()
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
	db.RDB.SRem(redisCtx(), dirtyUidsKey, uid)
}

// ────────────────────────────────────────
//  背包数据操作
// ────────────────────────────────────────

// LoadInventory 上线时从 MySQL 全量加载进 Redis
// 设计：用 Lua 原子写入，只在 Redis Hash 为空时才写（避免覆盖其他连接刚写的新数据）
//   - 如果 Redis 已有数据（另一个登录先写了）→ 跳过，直接返回
//   - 如果 Redis 空（正常冷启动）→ 全量写入 MySQL 数据
func LoadInventory(uid string, items map[string]int64) error {
	if db.RDB == nil {
		return fmt.Errorf("redis not initialized")
	}
	key := inventoryKey(uid)

	// Lua 脚本：只有 Hash 为空时才写入（原子操作，避免 Del+Set 之间被其他线程插入新数据）
	// 返回 1=写入了 MySQL 数据，0=跳过（Redis 已有数据）
	if len(items) == 0 {
		// MySQL 也没数据，如果 Redis 也空就跳过，Redis 有数据就保留
		script := `
			if redis.call('HLEN', KEYS[1]) == 0 then
				return 1
			end
			return 0
		`
		result, err := db.RDB.Eval(redisCtx(), script, []string{key}).Int64()
		if err != nil {
			return err
		}
		if result == 1 {
			return nil // Redis 空 + MySQL 空 → 什么都不做
		}
		return nil // Redis 有数据 → 保留，不覆盖
	}

	// 有 MySQL 数据 → Lua 原子判断+写入
	// ARGV 格式: count, item1, val1, item2, val2, ...
	args := make([]interface{}, 0, 1+len(items)*2)
	args = append(args, len(items))
	for itemID, count := range items {
		args = append(args, itemID, strconv.FormatInt(count, 10))
	}
	script := `
		if redis.call('HLEN', KEYS[1]) > 0 then
			return 0  -- Redis 已有数据，跳过
		end
		local n = tonumber(ARGV[1])
		for i = 1, n do
			redis.call('HSET', KEYS[1], ARGV[2*i], ARGV[2*i+1])
		end
		return 1  -- 写入成功
	`
	_, err := db.RDB.Eval(redisCtx(), script, []string{key}, args...).Int64()
	return err
}

// AddItem 增加道具（Lua 脚本保证 HINCRBY + MarkDirty 原子性）
// 数量可为负（扣减也走这里，但推荐用 ConsumeItem 有足够性检查）
func AddItem(uid, itemID string, delta int64) (int64, error) {
	if db.RDB == nil {
		return 0, fmt.Errorf("redis not initialized")
	}
	key := inventoryKey(uid)
	script := `
		local new = redis.call('HINCRBY', KEYS[1], ARGV[1], ARGV[2])
		redis.call('SADD', KEYS[2], ARGV[3])
		return new
	`
	newCount, err := db.RDB.Eval(redisCtx(), script,
		[]string{key, dirtyUidsKey},
		itemID, strconv.FormatInt(delta, 10), uid,
	).Int64()
	if err != nil {
		log.Warn("redis AddItem failed", "uid", uid, "item", itemID, "err", err)
		return 0, err
	}
	return newCount, nil
}

// ConsumeItem 消耗道具（Lua 脚本保证"检查-扣减-返回新值"原子性）
// 返回 (success, newCount, err)
//
//	success=true,  newCount=扣减后数量, err=nil   → 成功
//	success=false, newCount=扣减前数量, err=nil   → 数量不足
//	success=false, newCount=0,        err!=nil   → Redis 错误
func ConsumeItem(uid, itemID string, need int64) (bool, int64, error) {
	if db.RDB == nil {
		return false, 0, fmt.Errorf("redis not initialized")
	}
	key := inventoryKey(uid)

	// Lua 脚本：检查 → 扣减 → 标记 dirty → 返回 {success(0/1), after_or_before}
	// 成功返回 {1, afterCount}，失败返回 {0, beforeCount}
	// KEYS: [1]=背包hash, [2]=dirty UID set
	// ARGV: [1]=itemID, [2]=need, [3]=uid
	script := `
		local cur = tonumber(redis.call('HGET', KEYS[1], ARGV[1]) or '0')
		if cur < tonumber(ARGV[2]) then
			return {0, cur}
		end
		redis.call('HINCRBY', KEYS[1], ARGV[1], -tonumber(ARGV[2]))
		redis.call('SADD', KEYS[2], ARGV[3])
		return {1, cur - tonumber(ARGV[2])}
	`
	result, err := db.RDB.Eval(redisCtx(), script,
		[]string{key, dirtyUidsKey},
		itemID, strconv.FormatInt(need, 10), uid).Slice()
	if err != nil {
		log.Warn("redis ConsumeItem failed", "uid", uid, "item", itemID, "err", err)
		return false, 0, err
	}
	// 解析 Lua 返回的数组 [success, count]
	if len(result) >= 2 {
		success, _ := strconv.ParseInt(fmt.Sprint(result[0]), 10, 64)
		count, _ := strconv.ParseInt(fmt.Sprint(result[1]), 10, 64)
		if success == 1 {
			return true, count, nil
		}
		return false, count, nil // 数量不足，count = before
	}
	return false, 0, fmt.Errorf("unexpected lua result length")
}

// GetItemCount 查询单个道具数量
func GetItemCount(uid, itemID string) (int64, error) {
	if db.RDB == nil {
		return 0, fmt.Errorf("redis not initialized")
	}
	val, err := db.RDB.HGet(redisCtx(), inventoryKey(uid), itemID).Result()
	if err != nil {
		if errors.Is(err, redis.Nil) {
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
	raw, err := db.RDB.HGetAll(redisCtx(), inventoryKey(uid)).Result()
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
// 用 Lua 脚本原子清理，避免 HGetAll + N 次 HDel 的多次 RT
func RemoveZeroCount(uid string) {
	if db.RDB == nil {
		return
	}
	script := `
		local all = redis.call('HGETALL', KEYS[1])
		for i = 1, #all, 2 do
			local count = tonumber(all[i+1])
			if count <= 0 then
				redis.call('HDEL', KEYS[1], all[i])
			end
		end
	`
	db.RDB.Eval(redisCtx(), script, []string{inventoryKey(uid)})
}
