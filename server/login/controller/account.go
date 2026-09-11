package controller

import (
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"mssgServer/infra/log"
	"mssgServer/net"
	"mssgServer/server/login/dao"
	"mssgServer/server/login/proto"
	"strings"
)

var DefaultAccount = &Account{}

type Account struct {
}

func (a *Account) Router(r *net.Router) {
	g := r.Group("account")
	g.AddRoute("login", a.login)
	g.AddRoute("register", a.register)
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
