package controller

import (
	"encoding/json"
	"fmt"
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
	// 注意：ReqBody.Msg 是 interface{}，需要二次转换
	parseMsg(req.Body.Msg, body)

	// ✅ 查数据库验证用户

	user, err := dao.UserDAO.GetByUsername(body.Username)

	if user == nil {
		rsp.Body.Code = -1
		rsp.Body.Msg = map[string]string{"errMsg": "用户不存在"}
		return
	}
	if err != nil {
		rsp.Body.Code = -1
		rsp.Body.Msg = map[string]string{"errMsg": "数据库错误"}
		fmt.Println("[DB] query error:", err)
		return
	}

	if user.Password != body.Password {
		rsp.Body.Code = -1
		rsp.Body.Msg = map[string]string{"errMsg": "密码错误"}
		return
	}

	fmt.Printf(" Account login ReadMessage: %v\n", req.Body)
	loginResp := &proto.LoginResp{}
	loginResp.Username = "admin"
	loginResp.Password = "123456"
	loginResp.Session = "123456"
	loginResp.UId = "123456"
	rsp.Body.Msg = loginResp
}

// parseMsg 把 ReqBody.Msg (interface{}) 反序列化到具体结构体
func parseMsg(raw interface{}, target interface{}) error {
	data, err := json.Marshal(raw)
	if err != nil {
		return err
	}
	return json.Unmarshal(data, target)
}
