package net

import (
	"encoding/json"
)

// ────────────────────────────────────────
//  通用消息工具（所有 controller 共用）
// ────────────────────────────────────────

// ParseMsg 把 ReqBody.Msg (interface{}) 反序列化到具体结构体
// 原理：先 Marshal 成 JSON，再 Unmarshal 到目标结构体
// 这样无论客户端发的是 map[string]interface{} 还是 JSON 字符串都能解析
func ParseMsg(raw interface{}, target interface{}) error {
	if raw == nil {
		return nil
	}
	data, err := json.Marshal(raw)
	if err != nil {
		return err
	}
	return json.Unmarshal(data, target)
}

// ────────────────────────────────────────
//  通用连接属性工具（所有 controller 共用）
// ────────────────────────────────────────

// GetUid 从连接取 uid（登录时 SetProperty 设的）
// 返回 "" 表示未登录或属性丢失
func GetUid(conn WSConn) string {
	v, err := conn.GetProperty("uid")
	if err != nil || v == nil {
		return ""
	}
	uid, _ := v.(string)
	return uid
}

// GetUsername 从连接取 username
func GetUsername(conn WSConn) string {
	v, err := conn.GetProperty("username")
	if err != nil || v == nil {
		return ""
	}
	name, _ := v.(string)
	return name
}

// GetAuthInfo 一次取 uid + username（减少 GetProperty 调用次数）
func GetAuthInfo(conn WSConn) (uid, username string) {
	uid = GetUid(conn)
	username = GetUsername(conn)
	return
}
