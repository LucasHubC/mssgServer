package controller

import (
	"crypto/rand"
	"encoding/hex"
	"mssgServer/infra/log"
	rediscache "mssgServer/infra/redis"
	"mssgServer/net"
	"mssgServer/server/chat/proto"
	"strings"
	"time"
)

var DefaultChat = &Chat{}

type Chat struct{}

// Router 消息名统一前缀: chat.xxx
func (c *Chat) Router(r *net.Router) {
	g := r.Group("chat")
	g.AddRoute("world.send", c.sendWorld)
	g.AddRoute("private.send", c.sendPrivate)
	g.AddRoute("group.send", c.sendGroup)
	g.AddRoute("group.join", c.joinGroup)
	g.AddRoute("group.leave", c.leaveGroup)
	g.AddRoute("group.members", c.getGroupMembers)
	g.AddRoute("online.get", c.getOnline)
}

// ────────────────────────────────────────
//
//	世界频道聊天（跨服全服广播）
//
// ────────────────────────────────────────
func (c *Chat) sendWorld(req *net.WsMsgReq, rsp *net.WsMsgResp) {
	body := &proto.SendWorldReq{}
	if err := net.ParseMsg(req.Body.Msg, body); err != nil {
		rsp.Body.Code = -1
		rsp.Body.Msg = map[string]string{"errMsg": "参数解析失败"}
		return
	}

	uid, username := net.GetAuthInfo(req.Conn)
	if uid == "" {
		rsp.Body.Code = -1
		rsp.Body.Msg = map[string]string{"errMsg": "请先登录"}
		return
	}

	body.Content = strings.TrimSpace(body.Content)
	if body.Content == "" || len(body.Content) > 500 {
		rsp.Body.Code = -1
		rsp.Body.Msg = map[string]string{"errMsg": "内容不能为空且不超过 500 字"}
		return
	}

	msg := proto.WorldChatMsg{
		MsgId:     genMsgId(),
		FromUid:   uid,
		FromName:  username,
		Content:   body.Content,
		Timestamp: time.Now().Unix(),
	}

	log.Info("[Chat] world", "from", username, "msgId", msg.MsgId, "content", truncate(body.Content, 30))
	net.ConnManager.BroadcastAll("chat.world.new", msg)

	rsp.Body.Code = 0
	rsp.Body.Msg = &proto.SendWorldResp{MsgId: msg.MsgId}
}

// ────────────────────────────────────────
//
//	私聊（点对点，自动跨服）
//
// ────────────────────────────────────────
func (c *Chat) sendPrivate(req *net.WsMsgReq, rsp *net.WsMsgResp) {
	body := &proto.SendPrivateReq{}
	if err := net.ParseMsg(req.Body.Msg, body); err != nil {
		rsp.Body.Code = -1
		rsp.Body.Msg = map[string]string{"errMsg": "参数解析失败"}
		return
	}

	uid, username := net.GetAuthInfo(req.Conn)
	if uid == "" {
		rsp.Body.Code = -1
		rsp.Body.Msg = map[string]string{"errMsg": "请先登录"}
		return
	}
	if body.ToUid == "" || body.ToUid == uid {
		rsp.Body.Code = -1
		rsp.Body.Msg = map[string]string{"errMsg": "目标用户不合法"}
		return
	}

	body.Content = strings.TrimSpace(body.Content)
	if body.Content == "" || len(body.Content) > 500 {
		rsp.Body.Code = -1
		rsp.Body.Msg = map[string]string{"errMsg": "内容不能为空且不超过 500 字"}
		return
	}

	_, found := rediscache.GetUserNode(body.ToUid)

	msg := proto.PrivateChatMsg{
		MsgId:     genMsgId(),
		FromUid:   uid,
		FromName:  username,
		Content:   body.Content,
		Timestamp: time.Now().Unix(),
	}

	log.Info("[Chat] private", "from", username, "to", body.ToUid, "online", found)
	if found {
		net.ConnManager.PushToUid(body.ToUid, "chat.private.new", msg)
	}

	rsp.Body.Code = 0
	rsp.Body.Msg = &proto.SendPrivateResp{MsgId: msg.MsgId, Offline: !found}
}

// ────────────────────────────────────────
//
//	群组聊天（跨服群组消息）
//
// ────────────────────────────────────────
func (c *Chat) sendGroup(req *net.WsMsgReq, rsp *net.WsMsgResp) {
	body := &proto.SendGroupReq{}
	if err := net.ParseMsg(req.Body.Msg, body); err != nil {
		rsp.Body.Code = -1
		rsp.Body.Msg = map[string]string{"errMsg": "参数解析失败"}
		return
	}

	uid, username := net.GetAuthInfo(req.Conn)
	if uid == "" {
		rsp.Body.Code = -1
		rsp.Body.Msg = map[string]string{"errMsg": "请先登录"}
		return
	}

	body.Content = strings.TrimSpace(body.Content)
	if body.Content == "" || len(body.Content) > 500 {
		rsp.Body.Code = -1
		rsp.Body.Msg = map[string]string{"errMsg": "内容不能为空且不超过 500 字"}
		return
	}

	// 只有群成员才能发
	if !rediscache.IsGroupMember(body.GroupId, uid) {
		rsp.Body.Code = -1
		rsp.Body.Msg = map[string]string{"errMsg": "你不在这个群组里"}
		return
	}

	msg := proto.GroupChatMsg{
		MsgId:     genMsgId(),
		GroupId:   body.GroupId,
		FromUid:   uid,
		FromName:  username,
		Content:   body.Content,
		Timestamp: time.Now().Unix(),
	}

	log.Info("[Chat] group", "from", username, "group", body.GroupId, "msgId", msg.MsgId)

	// 发布到 Redis，SubscribeRemoteBroadcast 里 "group" case 会推给群内所有成员
	// （包括本机 + 其他节点的）
	rediscache.PublishBroadcast(rediscache.BroadcastMsg{
		Action: "group",
		UID:    body.GroupId,
		Name:   "chat.group.new",
		Data:   msg,
	})

	// 🔵 给自己也推一份（因为 FromNode=自己，SubscribeRemoteBroadcast 会跳过）
	if ws, ok := req.Conn.(*net.WSServer); ok {
		ws.Push("chat.group.new", msg)
	}

	rsp.Body.Code = 0
	rsp.Body.Msg = map[string]string{"msgId": msg.MsgId}
}

// ────────────────────────────────────────
//
//	加入群组
//
// ────────────────────────────────────────
func (c *Chat) joinGroup(req *net.WsMsgReq, rsp *net.WsMsgResp) {
	body := &proto.GroupMemberReq{}
	if err := net.ParseMsg(req.Body.Msg, body); err != nil {
		rsp.Body.Code = -1
		rsp.Body.Msg = map[string]string{"errMsg": "参数解析失败"}
		return
	}
	uid := net.GetUid(req.Conn)
	if uid == "" {
		rsp.Body.Code = -1
		rsp.Body.Msg = map[string]string{"errMsg": "请先登录"}
		return
	}

	rediscache.JoinGroup(body.GroupId, uid)
	log.Info("[Chat] join group", "uid", uid, "group", body.GroupId)

	members := rediscache.GetGroupMembers(body.GroupId)
	rsp.Body.Code = 0
	rsp.Body.Msg = &proto.GroupMemberResp{
		GroupId: body.GroupId,
		Members: members,
		Count:   len(members),
	}
}

// ────────────────────────────────────────
//
//	退出群组
//
// ────────────────────────────────────────
func (c *Chat) leaveGroup(req *net.WsMsgReq, rsp *net.WsMsgResp) {
	body := &proto.GroupMemberReq{}
	if err := net.ParseMsg(req.Body.Msg, body); err != nil {
		rsp.Body.Code = -1
		rsp.Body.Msg = map[string]string{"errMsg": "参数解析失败"}
		return
	}
	uid := net.GetUid(req.Conn)
	if uid == "" {
		rsp.Body.Code = -1
		rsp.Body.Msg = map[string]string{"errMsg": "请先登录"}
		return
	}

	rediscache.LeaveGroup(body.GroupId, uid)
	log.Info("[Chat] leave group", "uid", uid, "group", body.GroupId)

	rsp.Body.Code = 0
	rsp.Body.Msg = map[string]string{"msg": "ok"}
}

// ────────────────────────────────────────
//
//	查询群组成员
//
// ────────────────────────────────────────
func (c *Chat) getGroupMembers(req *net.WsMsgReq, rsp *net.WsMsgResp) {
	body := &proto.GroupMemberReq{}
	if err := net.ParseMsg(req.Body.Msg, body); err != nil {
		rsp.Body.Code = -1
		rsp.Body.Msg = map[string]string{"errMsg": "参数解析失败"}
		return
	}
	members := rediscache.GetGroupMembers(body.GroupId)
	rsp.Body.Code = 0
	rsp.Body.Msg = &proto.GroupMemberResp{
		GroupId: body.GroupId,
		Members: members,
		Count:   len(members),
	}
}

// ────────────────────────────────────────
//
//	本机在线用户列表
//
// ────────────────────────────────────────
func (c *Chat) getOnline(req *net.WsMsgReq, rsp *net.WsMsgResp) {
	uids := net.ConnManager.OnlineUids()
	rsp.Body.Code = 0
	rsp.Body.Msg = &proto.OnlineUsersResp{Total: len(uids), Uids: uids}
}

// ────────────────────────────────────────
//  工具函数
// ────────────────────────────────────────

func genMsgId() string {
	b := make([]byte, 8)
	rand.Read(b)
	return hex.EncodeToString(b)
}

func truncate(s string, n int) string {
	if len(s) <= n {
		return s
	}
	return s[:n] + "..."
}
