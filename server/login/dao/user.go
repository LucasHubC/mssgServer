package dao

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"mssgServer/config/db"
	rediscache "mssgServer/infra/redis"
	"strings"
)

// User 用户表对应的数据结构
type User struct {
	Id       int64  `json:"id"`
	UId      string `json:"uid"`
	Username string `json:"username"`
	Password string `json:"password"`
}

// UserDAO 账户相关数据库操作（单例）
var UserDAO = &userDAO{}

type userDAO struct{}

// ────────────────────────────────────────
//  查询类
// ────────────────────────────────────────

// GetById 按主键查单条
func (d *userDAO) GetById(id int64) (*User, error) {
	u := &User{}
	err := db.DB.QueryRow(
		"SELECT id, uid, username, password FROM users WHERE id = ?",
		id,
	).Scan(&u.Id, &u.UId, &u.Username, &u.Password)

	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return u, nil
}

// GetByUid 按业务 UId 查（带 Redis 缓存）
func (d *userDAO) GetByUid(uid string) (*User, error) {
	// ① 先查缓存
	key := rediscache.UserCacheKeyByUid(uid)
	if data, ok := rediscache.GetUserCache(key); ok {
		var u User
		if err := json.Unmarshal(data, &u); err == nil {
			return &u, nil
		}
	}

	// ② 查 MySQL
	u := &User{}
	err := db.DB.QueryRow(
		"SELECT id, uid, username, password FROM users WHERE uid = ?",
		uid,
	).Scan(&u.Id, &u.UId, &u.Username, &u.Password)

	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	// ③ 回写缓存
	data, _ := json.Marshal(u)
	rediscache.SetUserCache(key, data)
	return u, nil
}

// GetByUsername 按用户名查（登录时用，带 Redis 缓存）
func (d *userDAO) GetByUsername(username string) (*User, error) {
	// ① 先查缓存
	key := rediscache.UserCacheKeyByUsername(username)
	if data, ok := rediscache.GetUserCache(key); ok {
		var u User
		if err := json.Unmarshal(data, &u); err == nil {
			return &u, nil
		}
	}

	// ② 查 MySQL
	u := &User{}
	err := db.DB.QueryRow(
		"SELECT id, uid, username, password FROM users WHERE username = ?",
		username,
	).Scan(&u.Id, &u.UId, &u.Username, &u.Password)

	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	// ③ 回写缓存（同时缓存 uid key，方便 GetByUid 也能命中）
	data, _ := json.Marshal(u)
	rediscache.SetUserCache(key, data)
	rediscache.SetUserCache(rediscache.UserCacheKeyByUid(u.UId), data)
	return u, nil
}

// ExistsByUsername 用户名是否已存在（注册时用）
func (d *userDAO) ExistsByUsername(username string) (bool, error) {
	var count int
	err := db.DB.QueryRow(
		"SELECT COUNT(*) FROM users WHERE username = ?",
		username,
	).Scan(&count)
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

// ExistsByUid UId 是否已存在
func (d *userDAO) ExistsByUid(uid string) (bool, error) {
	var count int
	err := db.DB.QueryRow(
		"SELECT COUNT(*) FROM users WHERE uid = ?",
		uid,
	).Scan(&count)
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

// List 分页查用户
func (d *userDAO) List(page, pageSize int) ([]User, int, error) {
	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 {
		pageSize = 10
	}
	offset := (page - 1) * pageSize

	// 先查总数
	var total int
	if err := db.DB.QueryRow("SELECT COUNT(*) FROM users").Scan(&total); err != nil {
		return nil, 0, fmt.Errorf("count users failed: %w", err)
	}

	// 再查分页数据
	rows, err := db.DB.Query(
		"SELECT id, uid, username, password FROM users ORDER BY id DESC LIMIT ? OFFSET ?",
		pageSize, offset,
	)
	if err != nil {
		return nil, 0, fmt.Errorf("list users failed: %w", err)
	}
	defer rows.Close()

	list := make([]User, 0, pageSize)
	for rows.Next() {
		var u User
		if err := rows.Scan(&u.Id, &u.UId, &u.Username, &u.Password); err != nil {
			return nil, 0, fmt.Errorf("scan user failed: %w", err)
		}
		list = append(list, u)
	}
	if err := rows.Err(); err != nil {
		return nil, 0, err
	}
	return list, total, nil
}

// Search 按用户名模糊搜索
func (d *userDAO) Search(keyword string) ([]User, error) {
	keyword = strings.TrimSpace(keyword)
	if keyword == "" {
		return []User{}, nil
	}

	rows, err := db.DB.Query(
		"SELECT id, uid, username, password FROM users WHERE username LIKE ? ORDER BY id DESC",
		"%"+keyword+"%",
	)
	if err != nil {
		return nil, fmt.Errorf("search users failed: %w", err)
	}
	defer rows.Close()

	list := make([]User, 0)
	for rows.Next() {
		var u User
		if err := rows.Scan(&u.Id, &u.UId, &u.Username, &u.Password); err != nil {
			return nil, err
		}
		list = append(list, u)
	}
	return list, rows.Err()
}

// Count 统计总用户数
func (d *userDAO) Count() (int, error) {
	var total int
	err := db.DB.QueryRow("SELECT COUNT(*) FROM users").Scan(&total)
	return total, err
}

// ────────────────────────────────────────
//  写入类
// ────────────────────────────────────────

// Create 新增用户，返回自增主键 ID
// 同时会把自增 ID 写回 u.Id（调用者可以直接从 u.Id 拿到）
func (d *userDAO) Create(u *User) (int64, error) {
	result, err := db.DB.Exec(
		"INSERT INTO users (uid, username, password) VALUES (?, ?, ?)",
		u.UId, u.Username, u.Password,
	)
	if err != nil {
		return 0, fmt.Errorf("create user failed: %w", err)
	}
	// 把自增 ID 写回 u.Id，方便调用者后续使用
	u.Id, _ = result.LastInsertId()
	return u.Id, nil
}

// Update 修改用户基础信息（同步删缓存）
func (d *userDAO) Update(u *User) error {
	_, err := db.DB.Exec(
		"UPDATE users SET username = ?, password = ? WHERE id = ?",
		u.Username, u.Password, u.Id,
	)
	if err != nil {
		return fmt.Errorf("update user failed: %w", err)
	}
	// 删缓存，下次读会从 DB 拉最新的
	rediscache.DelUserCache(rediscache.UserCacheKeyByUid(u.UId))
	rediscache.DelUserCache(rediscache.UserCacheKeyByUsername(u.Username))
	return nil
}

// UpdatePassword 修改密码（只传 uid 和新密码即可）
func (d *userDAO) UpdatePassword(uid, newPwd string) error {
	// ① 先查一次，拿到 username（后面删缓存要用）
	u, err := d.GetByUid(uid)
	if err != nil {
		return fmt.Errorf("query user failed: %w", err)
	}
	if u == nil {
		return errors.New("user not found, no rows updated")
	}

	// ② 执行 UPDATE
	result, err := db.DB.Exec(
		"UPDATE users SET password = ? WHERE uid = ?",
		newPwd, uid,
	)
	if err != nil {
		return fmt.Errorf("update password failed: %w", err)
	}
	if affected, _ := result.RowsAffected(); affected == 0 {
		return errors.New("user not found, no rows updated")
	}

	// ③ 删缓存（两个 key 都删干净）
	rediscache.DelUserCache(
		rediscache.UserCacheKeyByUid(uid),
		rediscache.UserCacheKeyByUsername(u.Username),
	)
	return nil
}

// UpdateUsername 修改用户名
func (d *userDAO) UpdateUsername(uid, newName string) error {
	// ① 先查一次，拿到旧 username（后面删缓存要用）
	oldUser, err := d.GetByUid(uid)
	if err != nil {
		return fmt.Errorf("query user failed: %w", err)
	}
	if oldUser == nil {
		return errors.New("user not found, no rows updated")
	}

	// ② 执行 UPDATE
	result, err := db.DB.Exec(
		"UPDATE users SET username = ? WHERE uid = ?",
		newName, uid,
	)
	if err != nil {
		return fmt.Errorf("update username failed: %w", err)
	}
	if affected, _ := result.RowsAffected(); affected == 0 {
		return errors.New("user not found, no rows updated")
	}

	// ③ 删旧缓存 + 预写新缓存（合并一次 DB 查询搞定）
	rediscache.DelUserCache(
		rediscache.UserCacheKeyByUid(uid),
		rediscache.UserCacheKeyByUsername(oldUser.Username),
	)
	// 预写新 username 缓存
	oldUser.Username = newName
	data, _ := json.Marshal(oldUser)
	rediscache.SetUserCache(rediscache.UserCacheKeyByUsername(newName), data)
	rediscache.SetUserCache(rediscache.UserCacheKeyByUid(uid), data)

	return nil
}

// DeleteById 按主键删
func (d *userDAO) DeleteById(id int64) error {
	// ① 先查一下拿到缓存 key 需要的值（无缓存时直接查 DB）
	u, err := d.GetById(id)
	if err != nil {
		return fmt.Errorf("query user failed: %w", err)
	}
	if u == nil {
		return errors.New("user not found, no rows deleted")
	}

	// ② 删 DB + 删缓存（并发安全：删 DB 后立即删缓存，下次读取自然回源）
	_, err = db.DB.Exec("DELETE FROM users WHERE id = ?", id)
	if err != nil {
		return fmt.Errorf("delete user failed: %w", err)
	}

	// ③ 删缓存
	rediscache.DelUserCache(
		rediscache.UserCacheKeyByUid(u.UId),
		rediscache.UserCacheKeyByUsername(u.Username),
	)
	return nil
}

// DeleteByUid 按业务 UId 删
func (d *userDAO) DeleteByUid(uid string) error {
	// ① 先查一下拿到 username 以删干净缓存
	u, err := d.GetByUid(uid)
	if err != nil {
		return fmt.Errorf("query user failed: %w", err)
	}
	if u == nil {
		return errors.New("user not found, no rows deleted")
	}

	// ② 删 DB
	_, err = db.DB.Exec("DELETE FROM users WHERE uid = ?", uid)
	if err != nil {
		return fmt.Errorf("delete user failed: %w", err)
	}

	// ③ 删缓存
	rediscache.DelUserCache(
		rediscache.UserCacheKeyByUid(uid),
		rediscache.UserCacheKeyByUsername(u.Username),
	)
	return nil
}

// ────────────────────────────────────────
//  事务示例（注册用户 + 初始化金币 原子操作）
// ────────────────────────────────────────

// CreateWithBalance 事务：创建用户同时初始化金币账户
// 演示 DAO 如何配合 db.DB.Begin() 做跨表原子操作
func (d *userDAO) CreateWithBalance(u *User, initBalance int64) error {
	tx, err := db.DB.Begin()
	if err != nil {
		return fmt.Errorf("begin tx failed: %w", err)
	}
	defer func() {
		if err != nil {
			tx.Rollback()
		}
	}()

	// 1. 插入用户
	res, err := tx.Exec(
		"INSERT INTO users (uid, username, password) VALUES (?, ?, ?)",
		u.UId, u.Username, u.Password,
	)
	if err != nil {
		return fmt.Errorf("insert user failed: %w", err)
	}
	userId, _ := res.LastInsertId()

	// 2. 初始化金币（假设你有个 balance 表）
	_, err = tx.Exec(
		"INSERT INTO balance (user_id, coin) VALUES (?, ?)",
		userId, initBalance,
	)
	if err != nil {
		return fmt.Errorf("insert balance failed: %w", err)
	}

	return tx.Commit()
}
