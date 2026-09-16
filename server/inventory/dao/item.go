package dao

import (
	"fmt"
	"mssgServer/config/db"
	"mssgServer/infra/log"
)

// ItemDAO 背包数据持久化（单例）
var ItemDAO = &itemDAO{}

type itemDAO struct{}

// LoadByUid 上线时从 MySQL 加载全量背包
// 返回 map[item_id]count；没有记录返回空 map（不报错）
func (d *itemDAO) LoadByUid(uid string) (map[string]int64, error) {
	rows, err := db.DB.Query(
		"SELECT item_id, count FROM user_items WHERE uid = ? AND count > 0",
		uid,
	)
	if err != nil {
		return nil, fmt.Errorf("load inventory failed: %w", err)
	}
	defer rows.Close()

	items := make(map[string]int64)
	for rows.Next() {
		var itemID string
		var count int64
		if err := rows.Scan(&itemID, &count); err != nil {
			return nil, fmt.Errorf("scan inventory failed: %w", err)
		}
		items[itemID] = count
	}
	return items, rows.Err()
}

// SaveAll 全量保存背包（定时落盘 + 下线保存）
// 使用 REPLACE INTO 保证幂等（重复执行不产生脏数据）
// 调用方应确保 items 是经过过滤的（count > 0）
func (d *itemDAO) SaveAll(uid string, items map[string]int64) error {
	// 开启事务：先删后插，避免遗留已丢弃道具
	tx, err := db.DB.Begin()
	if err != nil {
		return fmt.Errorf("begin tx failed: %w", err)
	}
	defer func() {
		if err != nil {
			tx.Rollback()
		}
	}()

	// ① 删掉旧记录
	if _, err = tx.Exec("DELETE FROM user_items WHERE uid = ?", uid); err != nil {
		return fmt.Errorf("delete old inventory failed: %w", err)
	}

	// ② 批量 INSERT 新数据
	if len(items) > 0 {
		// 手动构造 VALUES 部分（简单可靠，Go 没批量 INSERT helper）
		placeholders := ""
		args := []interface{}{}
		first := true
		for itemID, count := range items {
			if count <= 0 {
				continue
			}
			if !first {
				placeholders += ","
			}
			placeholders += "(?, ?, ?)"
			args = append(args, uid, itemID, count)
			first = false
		}
		if len(args) > 0 {
			sql := "INSERT INTO user_items (uid, item_id, count) VALUES " + placeholders
			if _, err = tx.Exec(sql, args...); err != nil {
				return fmt.Errorf("insert inventory failed: %w", err)
			}
		}
	}

	if err = tx.Commit(); err != nil {
		return fmt.Errorf("commit inventory failed: %w", err)
	}
	log.Info("[Inventory] saved", "uid", uid, "items", len(items))
	return nil
}

// DeleteByUid 删除某用户的所有背包记录（删除账号时调用）
func (d *itemDAO) DeleteByUid(uid string) error {
	_, err := db.DB.Exec("DELETE FROM user_items WHERE uid = ?", uid)
	return err
}
