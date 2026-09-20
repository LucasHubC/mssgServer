package controller

import (
	"encoding/json"
	"mssgServer/infra/log"
	rediscache "mssgServer/infra/redis"
	"mssgServer/infra/util"
	"mssgServer/net"
	invDAO "mssgServer/server/inventory/dao"
	"mssgServer/server/login/dao"
	"mssgServer/server/login/proto"
	"net/http"
	"strings"

	"golang.org/x/crypto/bcrypt"
)

var DefaultAccount = &Account{}

type Account struct {
}

// Router 路由注册（长连接模式）
// account.* = 账号登录/注册
// room.* 已独立到 server/room/controller/room.go
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
	if err := net.ParseMsg(req.Body.Msg, body); err != nil {
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
	if err := verifyPassword(user.Password, body.Password); err != nil {
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

	// ✅ 登录成功：从 MySQL 加载背包进 Redis（防止 Redis 重启导致道具丢失）
	// 注意：items 必须提前声明在 if 外面，否则作用域挡住后面塞进 LoginResp
	items := make(map[string]int64)
	if loaded, err := invDAO.ItemDAO.LoadByUid(user.UId); err != nil {
		log.Warn("load inventory on login failed", "uid", user.UId, "err", err)
		// 降级：items 保持 make() 的空 map，客户端收到 {} 而不是 null
	} else {
		items = loaded
		if err := rediscache.LoadInventory(user.UId, items); err != nil {
			log.Warn("load inventory to redis failed", "uid", user.UId, "err", err)
		} else if len(items) > 0 {
			log.Info("inventory loaded on login", "uid", user.UId, "items", len(items))
		}
	}

	// ✅ 登录成功：签发 token 存 Redis（用于跨 Server 连接鉴权，比如 chatServer 的 ws://host:8004?token=xxx）
	token, err := rediscache.StoreSession(rediscache.Session{
		UId:      user.UId,
		Username: user.Username,
	})
	if err != nil {
		log.Warn("StoreSession failed", "err", err, "uid", user.UId)
	}

	rsp.Body.Msg = &proto.LoginResp{
		Username: user.Username,
		Token:    token,
		UId:      user.UId,
		Items:    items,
	}
}

// register 用户注册
func (a *Account) register(req *net.WsMsgReq, rsp *net.WsMsgResp) {
	body := &proto.RegisterReq{}
	if err := net.ParseMsg(req.Body.Msg, body); err != nil {
		rsp.Body.Code = -1
		rsp.Body.Msg = map[string]string{"errMsg": "参数解析失败"}
		return
	}

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
		Password: hashPassword(body.Password),
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
		net.WriteBadRequest(w, "method not allowed")
		return
	}

	body := &proto.LoginReq{}
	if err := json.NewDecoder(r.Body).Decode(body); err != nil {
		net.WriteBadRequest(w, "参数解析失败")
		return
	}

	user, err := dao.UserDAO.GetByUsername(body.Username)
	if err != nil {
		log.Error("loginHttp db error", "err", err)
		net.WriteError(w, "数据库错误")
		return
	}
	if user == nil {
		net.WriteNotFound(w, "用户不存在")
		return
	}
	if err := verifyPassword(user.Password, body.Password); err != nil {
		net.WriteUnauthorized(w, "密码错误")
		return
	}

	// ✅ 登录成功：签发 token 存 Redis
	token, err := rediscache.StoreSession(rediscache.Session{
		UId:      user.UId,
		Username: user.Username,
	})
	if err != nil {
		log.Error("loginHttp StoreSession failed", "err", err)
		net.WriteError(w, "签发登录凭证失败")
		return
	}

	log.Info("user login (http)", "uid", user.UId, "username", user.Username)

	// ✅ HTTP 登录也加载背包进 Redis（同上）
	items := make(map[string]int64)
	if loaded, err := invDAO.ItemDAO.LoadByUid(user.UId); err != nil {
		log.Warn("load inventory on http login failed", "uid", user.UId, "err", err)
	} else {
		items = loaded
		if err := rediscache.LoadInventory(user.UId, items); err != nil {
			log.Warn("load inventory to redis failed (http)", "uid", user.UId, "err", err)
		}
	}

	net.WriteJSON(w, map[string]any{
		"code": 0,
		"data": map[string]any{
			"token":    token,
			"uid":      user.UId,
			"username": user.Username,
			"items":    items,
		},
	})
}

// registerHttp POST /api/register
func (a *Account) registerHttp(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		net.WriteBadRequest(w, "method not allowed")
		return
	}

	body := &proto.RegisterReq{}
	if err := json.NewDecoder(r.Body).Decode(body); err != nil {
		net.WriteBadRequest(w, "参数解析失败")
		return
	}

	body.Username = strings.TrimSpace(body.Username)
	body.Password = strings.TrimSpace(body.Password)
	if body.Username == "" || body.Password == "" {
		net.WriteBadRequest(w, "用户名和密码不能为空")
		return
	}
	if len(body.Username) < 3 || len(body.Username) > 32 {
		net.WriteBadRequest(w, "用户名长度 3-32 位")
		return
	}
	if len(body.Password) < 6 || len(body.Password) > 128 {
		net.WriteBadRequest(w, "密码长度 6-128 位")
		return
	}

	exists, err := dao.UserDAO.ExistsByUsername(body.Username)
	if err != nil {
		log.Error("registerHttp check failed", "err", err)
		net.WriteError(w, "服务器错误")
		return
	}
	if exists {
		net.WriteConflict(w, "用户名已存在")
		return
	}

	u := &dao.User{
		UId:      genUId(),
		Username: body.Username,
		Password: hashPassword(body.Password),
	}
	if _, err := dao.UserDAO.Create(u); err != nil {
		log.Error("registerHttp create failed", "err", err)
		net.WriteError(w, "注册失败")
		return
	}

	log.Info("user registered (http)", "uid", u.UId, "username", u.Username)
	net.WriteJSON(w, map[string]any{
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
		net.WriteBadRequest(w, "method not allowed")
		return
	}

	// 路径: /api/user/test_xxx → 提取最后一段
	parts := strings.Split(strings.TrimSuffix(r.URL.Path, "/"), "/")
	uid := parts[len(parts)-1]
	if uid == "" {
		net.WriteBadRequest(w, "uid 缺失")
		return
	}

	user, err := dao.UserDAO.GetByUid(uid)
	if err != nil {
		log.Error("userInfoHttp db error", "err", err)
		net.WriteError(w, "数据库错误")
		return
	}
	if user == nil {
		net.WriteNotFound(w, "用户不存在")
		return
	}

	// ✅ HTTP API 绝对不能返回密码！
	net.WriteJSON(w, map[string]any{
		"code": 0,
		"data": map[string]string{
			"uid":      user.UId,
			"username": user.Username,
		},
	})
}

// ──── 鉴权辅助 ────

// ──── 需要鉴权的 HTTP Handlers ────

// meHttp GET /api/me —— 获取当前登录用户信息（必须带 token）
func (a *Account) meHttp(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		net.WriteBadRequest(w, "method not allowed")
		return
	}

	sess, ok := net.AuthFromHeader(r)
	if !ok {
		net.WriteUnauthorized(w, "token 无效或已过期")
		return
	}

	net.WriteJSON(w, map[string]any{
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
		net.WriteBadRequest(w, "method not allowed")
		return
	}

	// 从 Authorization 提取 token
	sess, ok := net.AuthFromHeader(r)
	if !ok {
		net.WriteUnauthorized(w, "请先登录")
		return
	}

	rediscache.DelSession(sess.Token)
	net.WriteJSON(w, map[string]any{"code": 0, "msg": "ok"})
}

// genMsgIdShort 短消息 ID：4 字节 = 8 位 hex（account 自己也用）
func genMsgIdShort() string { return util.RandomHex(4) }

const (
	uidPrefix = "uid_" // 业务 UId 前缀
)

// genUId 生成业务 UId: uid_ + 12 位 hex = 16 字符
func genUId() string { return uidPrefix + util.RandomHex(6) }

// hashPassword 用 bcrypt 哈希密码（注册时调）
func hashPassword(password string) string {
	hashed, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		// 理论上不会发生（DefaultCost=10，成本很低），万一出错回退存原文
		log.Warn("bcrypt hash failed, fallback to plaintext", "err", err)
		return password
	}
	return string(hashed)
}

// verifyPassword 比对密码（登录时调）
func verifyPassword(hashed, input string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashed), []byte(input))
}
