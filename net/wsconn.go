package net

import (
	"encoding/json"
	"fmt"
	"mssgServer/infra/log"
	"mssgServer/proto"
)

// ────────────────────────────────────────────────────
//  ReqBody / RespBody — Handler 层统一结构
//  阶段 2 核心：AuthUid/AuthUsername 从业务 payload 提到外层
//  Msg 优先是 protobuf typed struct（零 JSON），兜底是 map[string]interface{}
// ────────────────────────────────────────────────────

type ReqBody struct {
	Seq          int64       `json:"seq"`
	Name         string      `json:"name"`
	Msg          interface{} `json:"msg"`
	Proxy        string      `json:"proxy"`
	MsgJSON      string      `json:"-"` // 兜底：老版本 msg_json 字符串
	AuthUid      string      `json:"-"` // ✅ 阶段 2：外层鉴权字段
	AuthUsername string      `json:"-"` // ✅ 阶段 2
}

type RespBody struct {
	Seq     int64       `json:"seq"`
	Name    string      `json:"name"`
	Code    int         `json:"code"`
	Msg     interface{} `json:"msg"`
	MsgJSON string      `json:"-"`
}

// WsMsgReq / WsMsgResp — handler 层结构，传输层不感知
type WsMsgReq struct {
	Body *ReqBody
	Conn WSConn
}

type WsMsgResp struct {
	Body *RespBody
}

// ────────────────────────────────────────────────────
//  入站转换：proto.WSReq → ReqBody
//  优先级：oneof msg（typed struct）> msg_json（兜底）
// ────────────────────────────────────────────────────

func ReqBodyFromProto(p *proto.WSReq) (*ReqBody, error) {
	b := &ReqBody{
		Seq:          p.Seq,
		Name:         p.Name,
		Proxy:        p.Proxy,
		MsgJSON:      p.MsgJson,
		AuthUid:      p.AuthUid,      // ✅ 阶段 2：直接填外层鉴权
		AuthUsername: p.AuthUsername, // ✅ 阶段 2
	}

	// 阶段 2：优先从 oneof msg 取 typed struct
	if m := p.GetMsg(); m != nil {
		b.Msg = unwrapReqOneof(m)
		return b, nil
	}

	// 兜底：老版本 msg_json
	if p.MsgJson != "" {
		if err := json.Unmarshal([]byte(p.MsgJson), &b.Msg); err != nil {
			return nil, err
		}
	}
	return b, nil
}

// RespBodyFromProto：proto.WSResp → RespBody（入站用，BackendToClient）
func RespBodyFromProto(p *proto.WSResp) (*RespBody, error) {
	b := &RespBody{
		Seq:     p.Seq,
		Name:    p.Name,
		Code:    int(p.Code),
		MsgJSON: p.MsgJson,
	}

	if m := p.GetMsg(); m != nil {
		b.Msg = unwrapRespOneof(m)
		return b, nil
	}

	if p.MsgJson != "" {
		if err := json.Unmarshal([]byte(p.MsgJson), &b.Msg); err != nil {
			return nil, err
		}
	}
	return b, nil
}

// ────────────────────────────────────────────────────
//  出站转换：RespBody → proto.WSResp
//  智能路由：按 name 匹配 oneof 字段，typed struct 直接塞；兜底 marshal 到 msg_json
// ────────────────────────────────────────────────────

func RespBodyToProto(b *RespBody) (*proto.WSResp, error) {
	p := &proto.WSResp{
		Seq:  b.Seq,
		Name: b.Name,
		Code: int32(b.Code),
	}

	// 🔍 诊断:确认 Name 有没有
	if b.Name == "" && b.Code == 0 {
		log.Info("[RespBodyToProto] ⚠️ Name 为空", "code", b.Code, "msg_type", fmt.Sprintf("%T", b.Msg))
	}

	if b.Msg == nil {
		return p, nil
	}

	// 优先：typed struct 匹配 oneof
	if err := setRespOneof(p, b.Name, b.Msg); err == nil {
		return p, nil
	}

	// 兜底：marshal 到 msg_json
	if b.MsgJSON != "" {
		p.MsgJson = b.MsgJSON
	} else {
		j, err := json.Marshal(b.Msg)
		if err != nil {
			return nil, err
		}
		p.MsgJson = string(j)
	}
	return p, nil
}

// ReqBodyToProto：ReqBody → proto.WSReq（出站用，Broker injectAuth 场景）
func ReqBodyToProto(b *ReqBody) (*proto.WSReq, error) {
	p := &proto.WSReq{
		Seq:          b.Seq,
		Name:         b.Name,
		Proxy:        b.Proxy,
		AuthUid:      b.AuthUid,
		AuthUsername: b.AuthUsername,
	}

	if b.Msg == nil {
		if b.MsgJSON != "" {
			p.MsgJson = b.MsgJSON
		}
		return p, nil
	}

	// 优先：typed struct 匹配 oneof
	if err := setReqOneof(p, b.Name, b.Msg); err == nil {
		return p, nil
	}

	// 兜底：marshal 到 msg_json
	j, err := json.Marshal(b.Msg)
	if err != nil {
		return nil, err
	}
	p.MsgJson = string(j)
	return p, nil
}

// ────────────────────────────────────────────────────
//  Oneof 智能路由表 — name → proto 字段
//  handler 只需要把 typed struct 塞到 RespBody.Msg
//  这里按 name 自动匹配 oneof 字段
// ────────────────────────────────────────────────────

func setReqOneof(p *proto.WSReq, name string, msg interface{}) error {
	switch name {
	case "account.login":
		m, ok := msg.(*proto.AccountLogin)
		if !ok {
			return typeErr(name, msg)
		}
		p.Msg = &proto.WSReq_AccountLogin{AccountLogin: m}
	case "account.register":
		m, ok := msg.(*proto.AccountRegister)
		if !ok {
			return typeErr(name, msg)
		}
		p.Msg = &proto.WSReq_AccountRegister{AccountRegister: m}
	case "chat.world.send":
		m, ok := msg.(*proto.ChatWorldSend)
		if !ok {
			return typeErr(name, msg)
		}
		p.Msg = &proto.WSReq_ChatWorldSend{ChatWorldSend: m}
	case "chat.private.send":
		m, ok := msg.(*proto.ChatPrivateSend)
		if !ok {
			return typeErr(name, msg)
		}
		p.Msg = &proto.WSReq_ChatPrivateSend{ChatPrivateSend: m}
	case "chat.group.send":
		m, ok := msg.(*proto.ChatGroupSend)
		if !ok {
			return typeErr(name, msg)
		}
		p.Msg = &proto.WSReq_ChatGroupSend{ChatGroupSend: m}
	case "chat.group.join", "chat.group.leave", "chat.group.members":
		m, ok := msg.(*proto.GroupMember)
		if !ok {
			return typeErr(name, msg)
		}
		p.Msg = &proto.WSReq_GroupMemberReq{GroupMemberReq: m}
	case "room.create":
		m, ok := msg.(*proto.RoomCreate)
		if !ok {
			return typeErr(name, msg)
		}
		p.Msg = &proto.WSReq_RoomCreate{RoomCreate: m}
	case "room.join":
		m, ok := msg.(*proto.RoomJoin)
		if !ok {
			return typeErr(name, msg)
		}
		p.Msg = &proto.WSReq_RoomJoin{RoomJoin: m}
	case "room.leave":
		m, ok := msg.(*proto.RoomLeave)
		if !ok {
			return typeErr(name, msg)
		}
		p.Msg = &proto.WSReq_RoomLeave{RoomLeave: m}
	case "room.kick":
		m, ok := msg.(*proto.RoomKick)
		if !ok {
			return typeErr(name, msg)
		}
		p.Msg = &proto.WSReq_RoomKick{RoomKick: m}
	case "room.delete":
		m, ok := msg.(*proto.RoomDelete)
		if !ok {
			return typeErr(name, msg)
		}
		p.Msg = &proto.WSReq_RoomDelete{RoomDelete: m}
	case "room.send":
		m, ok := msg.(*proto.RoomSend)
		if !ok {
			return typeErr(name, msg)
		}
		p.Msg = &proto.WSReq_RoomSend{RoomSend: m}
	default:
		return errTypeMismatch // 兜底走 json
	}
	return nil
}

// unwrapReqOneof：把 oneof wrapper 解开成具体 proto.Message
// WSReq_AccountLogin{AccountLogin: *AccountLogin} → *AccountLogin
// 否则 json.Marshal 会输出 {"accountLogin":{...}}，handler ParseMsg 解析不到！
func unwrapReqOneof(m interface{}) interface{} {
	switch v := m.(type) {
	case *proto.WSReq_AccountLogin:
		return v.AccountLogin
	case *proto.WSReq_AccountRegister:
		return v.AccountRegister
	case *proto.WSReq_ChatWorldSend:
		return v.ChatWorldSend
	case *proto.WSReq_ChatPrivateSend:
		return v.ChatPrivateSend
	case *proto.WSReq_ChatGroupSend:
		return v.ChatGroupSend
	case *proto.WSReq_GroupMemberReq:
		return v.GroupMemberReq
	case *proto.WSReq_RoomCreate:
		return v.RoomCreate
	case *proto.WSReq_RoomJoin:
		return v.RoomJoin
	case *proto.WSReq_RoomLeave:
		return v.RoomLeave
	case *proto.WSReq_RoomKick:
		return v.RoomKick
	case *proto.WSReq_RoomDelete:
		return v.RoomDelete
	case *proto.WSReq_RoomSend:
		return v.RoomSend
	default:
		return m
	}
}

func setRespOneof(p *proto.WSResp, name string, msg interface{}) error {
	switch name {
	case "account.login":
		m, ok := msg.(*proto.AccountLoginResp)
		if !ok {
			return typeErr(name, msg)
		}
		p.Msg = &proto.WSResp_AccountLoginResp{AccountLoginResp: m}
	case "account.register":
		m, ok := msg.(*proto.AccountRegisterResp)
		if !ok {
			return typeErr(name, msg)
		}
		p.Msg = &proto.WSResp_AccountRegisterResp{AccountRegisterResp: m}
	case "chat.world.new":
		m, ok := msg.(*proto.ChatWorldNew)
		if !ok {
			return typeErr(name, msg)
		}
		p.Msg = &proto.WSResp_ChatWorldNew{ChatWorldNew: m}
	case "chat.private.new":
		m, ok := msg.(*proto.ChatPrivateNew)
		if !ok {
			return typeErr(name, msg)
		}
		p.Msg = &proto.WSResp_ChatPrivateNew{ChatPrivateNew: m}
	case "chat.group.new":
		m, ok := msg.(*proto.ChatGroupNew)
		if !ok {
			return typeErr(name, msg)
		}
		p.Msg = &proto.WSResp_ChatGroupNew{ChatGroupNew: m}
	case "chat.group.members":
		m, ok := msg.(*proto.GroupMemberResp)
		if !ok {
			return typeErr(name, msg)
		}
		p.Msg = &proto.WSResp_GroupMemberResp{GroupMemberResp: m}
	case "chat.online.get":
		m, ok := msg.(*proto.ChatOnlineResp)
		if !ok {
			return typeErr(name, msg)
		}
		p.Msg = &proto.WSResp_ChatOnlineResp{ChatOnlineResp: m}
	case "room.create":
		m, ok := msg.(*proto.RoomCreateResp)
		if !ok {
			return typeErr(name, msg)
		}
		p.Msg = &proto.WSResp_RoomCreateResp{RoomCreateResp: m}
	case "room.join":
		m, ok := msg.(*proto.RoomJoinResp)
		if !ok {
			return typeErr(name, msg)
		}
		p.Msg = &proto.WSResp_RoomJoinResp{RoomJoinResp: m}
	case "room.chat.new":
		m, ok := msg.(*proto.RoomChatNew)
		if !ok {
			return typeErr(name, msg)
		}
		p.Msg = &proto.WSResp_RoomChatNew{RoomChatNew: m}
	case "room.notice":
		// room.notice 目前 handler 里用 map[string]any 推送
		// setRespOneof 无法识别 map 类型，fallback json.Marshal 是预期内的
		// 如果以后 RoomNotice 有 protobuf 类型定义可以再加 type switch
		m, ok := msg.(*proto.RoomNotice)
		if !ok {
			return typeErr(name, msg)
		}
		p.Msg = &proto.WSResp_RoomNotice{RoomNotice: m}
	case "system.broadcast":
		m, ok := msg.(*proto.SystemBroadcast)
		if !ok {
			return typeErr(name, msg)
		}
		p.Msg = &proto.WSResp_SystemBroadcast{SystemBroadcast: m}
	case "system.push":
		m, ok := msg.(*proto.SystemPush)
		if !ok {
			return typeErr(name, msg)
		}
		p.Msg = &proto.WSResp_SystemPush{SystemPush: m}
	default:
		return errTypeMismatch
	}
	return nil
}

// unwrapRespOneof：同 unwrapReqOneof，解开 Resp oneof wrapper
func unwrapRespOneof(m interface{}) interface{} {
	switch v := m.(type) {
	case *proto.WSResp_AccountLoginResp:
		return v.AccountLoginResp
	case *proto.WSResp_AccountRegisterResp:
		return v.AccountRegisterResp
	case *proto.WSResp_ChatWorldNew:
		return v.ChatWorldNew
	case *proto.WSResp_ChatPrivateNew:
		return v.ChatPrivateNew
	case *proto.WSResp_ChatGroupNew:
		return v.ChatGroupNew
	case *proto.WSResp_GroupMemberResp:
		return v.GroupMemberResp
	case *proto.WSResp_ChatOnlineResp:
		return v.ChatOnlineResp
	case *proto.WSResp_RoomCreateResp:
		return v.RoomCreateResp
	case *proto.WSResp_RoomJoinResp:
		return v.RoomJoinResp
	case *proto.WSResp_RoomChatNew:
		return v.RoomChatNew
	case *proto.WSResp_RoomNotice:
		return v.RoomNotice
	case *proto.WSResp_SystemBroadcast:
		return v.SystemBroadcast
	case *proto.WSResp_SystemPush:
		return v.SystemPush
	default:
		return m
	}
}

// errTypeMismatch sentinel：oneof 类型断言失败，兜底 msg_json
var errTypeMismatch = fmt.Errorf("type mismatch for name")

func typeErr(name string, msg interface{}) error {
	return errTypeMismatch
}

// request 请求中放参数  取参数
type WSConn interface {
	SetProperty(key string, value interface{})
	GetProperty(key string) (interface{}, error)
	RemoveProperty(key string)
	Addr() string
	Push(name string, data interface{})
}
