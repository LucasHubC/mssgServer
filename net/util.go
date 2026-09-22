package net

import (
	"encoding/json"
)

// ────────────────────────────────────────
//  通用消息工具
// ────────────────────────────────────────

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
// 旧版：直连场景用（server.go 鉴权，没有 WsMsgReq）
// ────────────────────────────────────────

func GetUid(conn WSConn) string {
	v, err := conn.GetProperty("uid")
	if err != nil || v == nil {
		return ""
	}
	uid, _ := v.(string)
	return uid
}

func GetUsername(conn WSConn) string {
	v, err := conn.GetProperty("username")
	if err != nil || v == nil {
		return ""
	}
	name, _ := v.(string)
	return name
}

func GetAuthInfo(conn WSConn) (uid, username string) {
	uid = GetUid(conn)
	username = GetUsername(conn)
	return
}

// ────────────────────────────────────────
// 🔑 做法 A 新版：handler 里用这个！
// Broker 注入 __uid/__username 到 msg 里
// 后端一条 WS 代理 N 客户端 → Property 互相覆盖
// 必须从 msg.__uid 取，不能从 Property 取！
// ────────────────────────────────────────

func GetAuthInfoFromReq(req *WsMsgReq) (uid, username string) {
	if req.Body != nil && req.Body.Msg != nil {
		if m, ok := req.Body.Msg.(map[string]interface{}); ok {
			if v, _ := m["__uid"].(string); v != "" {
				uid = v
			}
			if v, _ := m["__username"].(string); v != "" {
				username = v
			}
		}
	}
	// fallback 到 Property（兼容直连场景）
	if uid == "" {
		if v, err := req.Conn.GetProperty("uid"); err == nil && v != nil {
			if s, ok := v.(string); ok {
				uid = s
			}
		}
	}
	if username == "" {
		if v, err := req.Conn.GetProperty("username"); err == nil && v != nil {
			if s, ok := v.(string); ok {
				username = s
			}
		}
	}
	return
}

func GetUidFromReq(req *WsMsgReq) string {
	uid, _ := GetAuthInfoFromReq(req)
	return uid
}

func GetUsernameFromReq(req *WsMsgReq) string {
	_, username := GetAuthInfoFromReq(req)
	return username
}
