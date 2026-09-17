package dao

import (
	"fmt"
	"mssgServer/config/db"
	"mssgServer/infra/log"
	"time"
)

// ItemOpLogDAO 道具操作日志（单例）
// 设计原则：写失败不阻塞主流程，异步写
var ItemOpLogDAO = &itemOpLogDAO{}

type itemOpLogDAO struct{}

// Insert 插入一条操作日志
// 调用方应异步调用（go ItemOpLogDAO.InsertSafe(...)），避免阻塞主线程
func (d *itemOpLogDAO) Insert(uid, itemID string, opType int8, count, beforeCount, afterCount int64, source string) error {
	// ✅ 用 DBLog（game_log 库），不是 DB（mssg 库）
	if db.DBLog == nil {
		return fmt.Errorf("log database not initialized")
	}
	const sql = "INSERT INTO item_op_log (uid, item_id, op_type, count, before_count, after_count, source) " +
		"VALUES (?, ?, ?, ?, ?, ?, ?)"
	_, err := db.DBLog.Exec(sql, uid, itemID, opType, count, beforeCount, afterCount, source)
	if err != nil {
		return fmt.Errorf("insert item_op_log failed: %w", err)
	}
	return nil
}

// InsertSafe 安全写入（带 panic recover + 3秒超时）
// 设计：同步执行（日志写 DB 毫秒级），进程退出时 Shutdown 完再 Close 不会丢
// 如果以后写入变慢，可改成有界 worker pool + WaitGroup
func (d *itemOpLogDAO) InsertSafe(uid, itemID string, opType int8, count, beforeCount, afterCount int64, source string) {
	// ① panic recover
	defer func() {
		if r := recover(); r != nil {
			log.Error("[ItemOpLog] panic recovered", "err", r, "uid", uid, "item_id", itemID)
		}
	}()
	// ② 3秒超时（防慢查询拖垮）
	done := make(chan error, 1)
	go func() {
		done <- d.Insert(uid, itemID, opType, count, beforeCount, afterCount, source)
	}()
	select {
	case err := <-done:
		if err != nil {
			log.Warn("[ItemOpLog] write failed (ignored)", "err", err, "uid", uid, "item_id", itemID)
		}
	case <-time.After(3 * time.Second):
		log.Warn("[ItemOpLog] write timeout (ignored)", "uid", uid, "item_id", itemID)
	}
}

// QueryByUid 分页查询某用户的道具操作记录（运营后台用）
func (d *itemOpLogDAO) QueryByUid(uid string, page, pageSize int) ([]map[string]interface{}, error) {
	if pageSize <= 0 || pageSize > 200 {
		pageSize = 200
	}
	if page < 1 {
		page = 1
	}
	offset := (page - 1) * pageSize

	rows, err := db.DBLog.Query(
		"SELECT id, item_id, op_type, count, before_count, after_count, source, created_at "+
			"FROM item_op_log WHERE uid = ? ORDER BY id DESC LIMIT ? OFFSET ?",
		uid, pageSize, offset,
	)
	if err != nil {
		return nil, fmt.Errorf("query item_op_log failed: %w", err)
	}
	defer rows.Close()

	columns, _ := rows.Columns()
	result := make([]map[string]interface{}, 0)
	for rows.Next() {
		values := make([]interface{}, len(columns))
		valuePtrs := make([]interface{}, len(columns))
		for i := range values {
			valuePtrs[i] = &values[i]
		}
		if err := rows.Scan(valuePtrs...); err != nil {
			return nil, fmt.Errorf("scan item_op_log failed: %w", err)
		}
		row := make(map[string]interface{})
		for i, col := range columns {
			row[col] = values[i]
		}
		result = append(result, row)
	}
	return result, rows.Err()
}

// DeleteBefore 清理指定日期之前的旧日志（定时任务用）
func (d *itemOpLogDAO) DeleteBefore(before time.Time) (int64, error) {
	if db.DBLog == nil {
		return 0, fmt.Errorf("log database not initialized")
	}
	res, err := db.DBLog.Exec("DELETE FROM item_op_log WHERE created_at < ?", before)
	if err != nil {
		return 0, fmt.Errorf("delete old item_op_log failed: %w", err)
	}
	return res.RowsAffected()
}
