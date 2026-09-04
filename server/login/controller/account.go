package controller

import (
	"encoding/json"
	"mssgServer/infra/log"
	"mssgServer/net"
	"mssgServer/server/login/dao"
	"mssgServer/server/login/proto"
)

var DefaultAccount = &Account{}

type Account struct {
}

func (a *Account) Router(r *net.Router) {
	g := r.Group("account")
	g.AddRoute("login", a.login)
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
		Password: user.Password,
		Session:  "sess_" + user.UId,
		UId:      user.UId,
	}
}

// parseMsg 把 ReqBody.Msg (interface{}) 反序列化到具体结构体
func parseMsg(raw interface{}, target interface{}) error {
	data, err := json.Marshal(raw)
	if err != nil {
		return err
	}
	return json.Unmarshal(data, target)
}
