package controller

import (
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"mssgServer/infra/log"
	rediscache "mssgServer/infra/redis"
	"mssgServer/net"
	"mssgServer/server/login/dao"
	"mssgServer/server/login/proto"
	"net/http"
	"strings"
)

var DefaultAccount = &Account{}

type Account struct {
}

// Router 贯由注册（长连接模式）
func (a *Account) Router(r *net.Router) {
	g := r.Group("account")
	g.AddRoute("login", a.login)
	g.AddRoute("register", a.register)
}

// ─────────────────────────────────────────
//
//	HTTP REST 路由注册（短连接模式）
//	纯标准库 net/http，不引入任何抽象层
//	调用方（apiLoginServer）自己起 http.Server
//
// ─────────────────────────────────────────
func (a *Account) RegisterHttpRoutes(mux *http.ServeMux) {
	// ── 公开接口（不需要 token）──
	mux.HandleFunc("/api/login", a.loginHttp)
	mux.HandleFunc("/api/register", a.registerHttp)
	mux.HandleFunc("/api/user/", a.userInfoHttp)

	// ── 需要鉴权的接口（必须带 Authorization header）──
	mux.HandleFunc("/api/me", a.meHttp)         // GET /api/me 获取当前用户信息
	mux.HandleFunc("/api/logout", a.logoutHttp) // POST /api/logout 登出
}

func (a *Account) login(req *net.WsMsgReq, rsp *net.WsMsgResp) {
	rsp.Body.Code = 0

	// 解析请求参数
	body := &proto.LoginReq{}
	if err := parseMsg(req.Body.Msg, body); err != nil {
		rsp.Body.Code = -1
		rsp.Body.Msg = map[string]string{"errMsg": "参数解析失败"}
		return
	}

	// 查数据库验证用户
	user, err := dao.UserDAO.GetByUsername(body.Username)
	if err != nil {
		log.Error("login db error", "err", err)
		rsp.Body.Code = -1
		rsp.Body.Msg = map[string]string{"errMsg": "数据库错误"}
		return
	}
	if user == nil {
		rsp.Body.Code = -1
		rsp.Body.Msg = map[string]string{"errMsg": "用户不存在"}
		return
	}
	if user.Password != body.Password {
		rsp.Body.Code = -1
		rsp.Body.Msg = map[string]string{"errMsg": "密码错误"}
		return
	}

	// ✅ 登录成功：存连接属性 + 通知 ConnectionManager 绑定 uid
	log.Info("user login", "uid", user.UId, "username", user.Username, "addr", req.Conn.Addr())

	req.Conn.SetProperty("uid", user.UId)
	req.Conn.SetProperty("username", user.Username)

	// 注意：BindUid 需要把 WSConn 强转回 *WSServer
	if ws, ok := req.Conn.(*net.WSServer); ok {
		net.ConnManager.BindUid(user.UId, ws)
	}

	rsp.Body.Msg = &proto.LoginResp{
		Username: user.Username,
		Session:  "sess_" + user.UId,
		UId:      user.UId,
	}
}

// register 用户注册
func (a *Account) register(req *net.WsMsgReq, rsp *net.WsMsgResp) {
	body := &proto.RegisterReq{}
	if err := parseMsg(req.Body.Msg, body); err != nil {
		rsp.Body.Code = -1
		rsp.Body.Msg = map[string]string{"errMsg": "参数解析失败"}
		return
	}

	// ① 参数校验
	body.Username = strings.TrimSpace(body.Username)
	body.Password = strings.TrimSpace(body.Password)
	if body.Username == "" || body.Password == "" {
		rsp.Body.Code = -1
		rsp.Body.Msg = map[string]string{"errMsg": "用户名和密码不能为空"}
		return
	}
	if len(body.Username) < 3 || len(body.Username) > 32 {
		rsp.Body.Code = -1
		rsp.Body.Msg = map[string]string{"errMsg": "用户名长度 3-32 位"}
		return
	}
	if len(body.Password) < 6 || len(body.Password) > 128 {
		rsp.Body.Code = -1
		rsp.Body.Msg = map[string]string{"errMsg": "密码长度 6-128 位"}
		return
	}

	// ② 用户名重复检查
	exists, err := dao.UserDAO.ExistsByUsername(body.Username)
	if err != nil {
		log.Error("register check username failed", "err", err)
		rsp.Body.Code = -1
		rsp.Body.Msg = map[string]string{"errMsg": "服务器错误"}
		return
	}
	if exists {
		rsp.Body.Code = -1
		rsp.Body.Msg = map[string]string{"errMsg": "用户名已存在"}
		return
	}

	// ③ 生成业务 UId（时间戳 + 随机 hex）
	u := &dao.User{
		UId:      genUId(),
		Username: body.Username,
		Password: body.Password, // 🟡 后面改成 bcrypt 哈希
	}

	// ④ 入库
	if _, err := dao.UserDAO.Create(u); err != nil {
		log.Error("register create user failed", "err", err)
		rsp.Body.Code = -1
		rsp.Body.Msg = map[string]string{"errMsg": "注册失败，请重试"}
		return
	}

	log.Info("user registered", "uid", u.UId, "username", u.Username)

	rsp.Body.Code = 0
	rsp.Body.Msg = &proto.RegisterResp{
		Username: u.Username,
		UId:      u.UId,
	}
}

// ──── HTTP REST Handlers ────

// loginHttp POST /api/login
// HTTP 短连接登录：成功后只返回 JSON，不绑定连接（因为 HTTP 请求完就断了）
func (a *Account) loginHttp(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, `{"code":-1,"errMsg":"method not allowed"}`, http.StatusMethodNotAllowed)
		return
	}

	body := &proto.LoginReq{}
	if err := json.NewDecoder(r.Body).Decode(body); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]any{"code": -1, "errMsg": "参数解析失败"})
		return
	}

	user, err := dao.UserDAO.GetByUsername(body.Username)
	if err != nil {
		log.Error("loginHttp db error", "err", err)
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]any{"code": -1, "errMsg": "数据库错误"})
		return
	}
	if user == nil {
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(map[string]any{"code": -1, "errMsg": "用户不存在"})
		return
	}
	if user.Password != body.Password {
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(map[string]any{"code": -1, "errMsg": "密码错误"})
		return
	}

	// ✅ 登录成功：签发 token 存 Redis
	token, err := rediscache.StoreSession(rediscache.Session{
		UId:      user.UId,
		Username: user.Username,
	})
	if err != nil {
		log.Error("loginHttp StoreSession failed", "err", err)
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]any{"code": -1, "errMsg": "签发登录凭证失败"})
		return
	}

	log.Info("user login (http)", "uid", user.UId, "username", user.Username)

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	json.NewEncoder(w).Encode(map[string]any{
		"code": 0,
		"data": map[string]any{
			"token":    token,
			"uid":      user.UId,
			"username": user.Username,
		},
	})
}

// registerHttp POST /api/register
func (a *Account) registerHttp(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, `{"code":-1,"errMsg":"method not allowed"}`, http.StatusMethodNotAllowed)
		return
	}

	body := &proto.RegisterReq{}
	if err := json.NewDecoder(r.Body).Decode(body); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]any{"code": -1, "errMsg": "参数解析失败"})
		return
	}

	body.Username = strings.TrimSpace(body.Username)
	body.Password = strings.TrimSpace(body.Password)
	if body.Username == "" || body.Password == "" {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]any{"code": -1, "errMsg": "用户名和密码不能为空"})
		return
	}
	if len(body.Username) < 3 || len(body.Username) > 32 {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]any{"code": -1, "errMsg": "用户名长度 3-32 位"})
		return
	}
	if len(body.Password) < 6 || len(body.Password) > 128 {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]any{"code": -1, "errMsg": "密码长度 6-128 位"})
		return
	}

	exists, err := dao.UserDAO.ExistsByUsername(body.Username)
	if err != nil {
		log.Error("registerHttp check failed", "err", err)
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]any{"code": -1, "errMsg": "服务器错误"})
		return
	}
	if exists {
		w.WriteHeader(http.StatusConflict)
		json.NewEncoder(w).Encode(map[string]any{"code": -1, "errMsg": "用户名已存在"})
		return
	}

	u := &dao.User{
		UId:      genUId(),
		Username: body.Username,
		Password: body.Password,
	}
	if _, err := dao.UserDAO.Create(u); err != nil {
		log.Error("registerHttp create failed", "err", err)
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]any{"code": -1, "errMsg": "注册失败"})
		return
	}

	log.Info("user registered (http)", "uid", u.UId, "username", u.Username)
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	json.NewEncoder(w).Encode(map[string]any{
		"code": 0,
		"data": proto.RegisterResp{
			Username: u.Username,
			UId:      u.UId,
		},
	})
}

// userInfoHttp GET /api/user/{uid}
func (a *Account) userInfoHttp(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, `{"code":-1,"errMsg":"method not allowed"}`, http.StatusMethodNotAllowed)
		return
	}

	// 路径: /api/user/test_xxx → 提取最后一段
	parts := strings.Split(strings.TrimSuffix(r.URL.Path, "/"), "/")
	uid := parts[len(parts)-1]
	if uid == "" {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]any{"code": -1, "errMsg": "uid 缺失"})
		return
	}

	user, err := dao.UserDAO.GetByUid(uid)
	if err != nil {
		log.Error("userInfoHttp db error", "err", err)
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]any{"code": -1, "errMsg": "数据库错误"})
		return
	}
	if user == nil {
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(map[string]any{"code": -1, "errMsg": "用户不存在"})
		return
	}

	// ✅ HTTP API 绝对不能返回密码！
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	json.NewEncoder(w).Encode(map[string]any{
		"code": 0,
		"data": map[string]string{
			"uid":      user.UId,
			"username": user.Username,
		},
	})
}

// ──── 鉴权辅助 ────

// authFromHeader 从 Authorization header 提取 token
// 支持两种格式：
//
//	Authorization: Bearer tk_abc123...
//	Authorization: tk_abc123...
//
// 返回 (Session, true) 表示通过；(nil, false) 表示失败
func authFromHeader(r *http.Request) (rediscache.Session, bool) {
	auth := r.Header.Get("Authorization")
	if auth == "" {
		return rediscache.Session{}, false
	}
	// 去掉 "Bearer " 前缀
	if strings.HasPrefix(auth, "Bearer ") {
		auth = strings.TrimPrefix(auth, "Bearer ")
	}
	if auth == "" {
		return rediscache.Session{}, false
	}
	return rediscache.GetSession(auth)
}

// writeUnauthorized 统一写 401 响应
func writeUnauthorized(w http.ResponseWriter, msg string) {
	w.WriteHeader(http.StatusUnauthorized)
	json.NewEncoder(w).Encode(map[string]any{"code": -1, "errMsg": msg})
}

// ──── 需要鉴权的 HTTP Handlers ────

// meHttp GET /api/me —— 获取当前登录用户信息（必须带 token）
func (a *Account) meHttp(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, `{"code":-1,"errMsg":"method not allowed"}`, http.StatusMethodNotAllowed)
		return
	}

	sess, ok := authFromHeader(r)
	if !ok {
		writeUnauthorized(w, "token 无效或已过期")
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	json.NewEncoder(w).Encode(map[string]any{
		"code": 0,
		"data": map[string]string{
			"uid":      sess.UId,
			"username": sess.Username,
		},
	})
}

// logoutHttp POST /api/logout —— 登出，销毁 token
func (a *Account) logoutHttp(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, `{"code":-1,"errMsg":"method not allowed"}`, http.StatusMethodNotAllowed)
		return
	}

	// 从 Authorization 提取 token
	auth := r.Header.Get("Authorization")
	if strings.HasPrefix(auth, "Bearer ") {
		auth = strings.TrimPrefix(auth, "Bearer ")
	}
	if auth == "" {
		writeUnauthorized(w, "请先登录")
		return
	}

	rediscache.DelSession(auth)

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	json.NewEncoder(w).Encode(map[string]any{"code": 0, "msg": "ok"})
}

// genUId 生成业务 UId: test_ + 12 位 hex = 16 字符
func genUId() string {
	b := make([]byte, 6)
	rand.Read(b)
	return "test_" + hex.EncodeToString(b)
}

// parseMsg 把 ReqBody.Msg (interface{}) 反序列化到具体结构体
func parseMsg(raw interface{}, target interface{}) error {
	data, err := json.Marshal(raw)
	if err != nil {
		return err
	}
	return json.Unmarshal(data, target)
}
