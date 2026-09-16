package net

import (
	"encoding/json"
	rediscache "mssgServer/infra/redis"
	"net/http"
	"strings"
)

// ────────────────────────────────────────
//  HTTP Token 鉴权（所有 HTTP controller 共用）
// ────────────────────────────────────────

// AuthFromHeader 从 Authorization header 提取 token 并查 Redis session
// 支持 "Bearer tk_xxx" 和 "tk_xxx" 两种格式
// 返回 (Session, true) 表示通过；(零值, false) 表示失败
func AuthFromHeader(r *http.Request) (rediscache.Session, bool) {
	auth := r.Header.Get("Authorization")
	if auth == "" {
		return rediscache.Session{}, false
	}
	if strings.HasPrefix(auth, "Bearer ") {
		auth = strings.TrimPrefix(auth, "Bearer ")
	}
	if auth == "" {
		return rediscache.Session{}, false
	}
	return rediscache.GetSession(auth)
}

// AuthUid 快捷版：直接返回 uid，失败返回 ("", false)
func AuthUid(r *http.Request) (string, bool) {
	sess, ok := AuthFromHeader(r)
	if !ok {
		return "", false
	}
	return sess.UId, true
}

// ────────────────────────────────────────
//  HTTP 统一响应工具
// ────────────────────────────────────────

// WriteJSONWithStatus 写 JSON 响应（带自定义状态码）
// 其他 Write* 函数底层都调这个
func WriteJSONWithStatus(w http.ResponseWriter, status int, v interface{}) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(v)
}

// WriteJSON 写 JSON 响应（200 OK）
func WriteJSON(w http.ResponseWriter, v interface{}) {
	WriteJSONWithStatus(w, http.StatusOK, v)
}

// WriteUnauthorized 写 401 响应
func WriteUnauthorized(w http.ResponseWriter, msg string) {
	WriteJSONWithStatus(w, http.StatusUnauthorized, map[string]any{"code": -1, "errMsg": msg})
}

// WriteError 写 500 响应
func WriteError(w http.ResponseWriter, msg string) {
	WriteJSONWithStatus(w, http.StatusInternalServerError, map[string]any{"code": -1, "errMsg": msg})
}

// WriteBadRequest 写 400 响应
func WriteBadRequest(w http.ResponseWriter, msg string) {
	WriteJSONWithStatus(w, http.StatusBadRequest, map[string]any{"code": -1, "errMsg": msg})
}

// WriteNotFound 写 404 响应
func WriteNotFound(w http.ResponseWriter, msg string) {
	WriteJSONWithStatus(w, http.StatusNotFound, map[string]any{"code": -1, "errMsg": msg})
}

// WriteConflict 写 409 响应
func WriteConflict(w http.ResponseWriter, msg string) {
	WriteJSONWithStatus(w, http.StatusConflict, map[string]any{"code": -1, "errMsg": msg})
}
