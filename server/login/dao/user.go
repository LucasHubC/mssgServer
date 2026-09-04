package dao

import (
	"database/sql"
	"errors"
	"fmt"
	"mssgServer/config/db"
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

// GetByUid 按业务 UId 查（uid 是业务唯一标识，比如 "1001"）
func (d *userDAO) GetByUid(uid string) (*User, error) {
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
	return u, nil
}

// GetByUsername 按用户名查（登录时用）
func (d *userDAO) GetByUsername(username string) (*User, error) {
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
func (d *userDAO) Create(u *User) (int64, error) {
	result, err := db.DB.Exec(
		"INSERT INTO users (uid, username, password) VALUES (?, ?, ?)",
		u.UId, u.Username, u.Password,
	)
	if err != nil {
		return 0, fmt.Errorf("create user failed: %w", err)
	}
	return result.LastInsertId()
}

// Update 修改用户基础信息
func (d *userDAO) Update(u *User) error {
	_, err := db.DB.Exec(
		"UPDATE users SET username = ?, password = ? WHERE id = ?",
		u.Username, u.Password, u.Id,
	)
	if err != nil {
		return fmt.Errorf("update user failed: %w", err)
	}
	return nil
}

// UpdatePassword 修改密码（只传 uid 和新密码即可）
func (d *userDAO) UpdatePassword(uid, newPwd string) error {
	result, err := db.DB.Exec(
		"UPDATE users SET password = ? WHERE uid = ?",
		newPwd, uid,
	)
	if err != nil {
		return fmt.Errorf("update password failed: %w", err)
	}
	affected, _ := result.RowsAffected()
	if affected == 0 {
		return errors.New("user not found, no rows updated")
	}
	return nil
}

// UpdateUsername 修改用户名
func (d *userDAO) UpdateUsername(uid, newName string) error {
	result, err := db.DB.Exec(
		"UPDATE users SET username = ? WHERE uid = ?",
		newName, uid,
	)
	if err != nil {
		return fmt.Errorf("update username failed: %w", err)
	}
	affected, _ := result.RowsAffected()
	if affected == 0 {
		return errors.New("user not found, no rows updated")
	}
	return nil
}

// DeleteById 按主键删
func (d *userDAO) DeleteById(id int64) error {
	result, err := db.DB.Exec("DELETE FROM users WHERE id = ?", id)
	if err != nil {
		return fmt.Errorf("delete user failed: %w", err)
	}
	affected, _ := result.RowsAffected()
	if affected == 0 {
		return errors.New("user not found, no rows deleted")
	}
	return nil
}

// DeleteByUid 按业务 UId 删
func (d *userDAO) DeleteByUid(uid string) error {
	result, err := db.DB.Exec("DELETE FROM users WHERE uid = ?", uid)
	if err != nil {
		return fmt.Errorf("delete user failed: %w", err)
	}
	affected, _ := result.RowsAffected()
	if affected == 0 {
		return errors.New("user not found, no rows deleted")
	}
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
