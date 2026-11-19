package controller

import (
	"mssgServer/infra/log"
	rediscache "mssgServer/infra/redis"
	"mssgServer/infra/util"
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

	uid, username := net.GetAuthInfoFromReq(req)
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

	// 本机 + 跨服广播（BroadcastAll 内部已包含 PublishBroadcast）
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

	uid, username := net.GetAuthInfoFromReq(req)
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

	// 热路径：不再打 Info 日志（3000+ QPS 下每条 IO 是瓶颈）
	// 如需统计 online 比例，用 metrics 或临时开 Debug
	log.Debug("[Chat] private", "from", username, "to", body.ToUid, "online", found)
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

	uid, username := net.GetAuthInfoFromReq(req)
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

	if body.GroupId == "" {
		rsp.Body.Code = -1
		rsp.Body.Msg = map[string]string{"errMsg": "group_id 不能为空"}
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

	// BroadcastAll 模式：先推本机组员（含自己），再广播其他节点
	net.ConnManager.PushToLocalGroupMembers(body.GroupId, "chat.group.new", msg)
	rediscache.PublishBroadcast(rediscache.BroadcastMsg{
		Action: "group",
		Target: body.GroupId,
		Name:   "chat.group.new",
		Data:   msg,
	})

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
	uid, username := net.GetAuthInfoFromReq(req)
	if uid == "" {
		rsp.Body.Code = -1
		rsp.Body.Msg = map[string]string{"errMsg": "请先登录"}
		return
	}

	if body.GroupId == "" {
		rsp.Body.Code = -1
		rsp.Body.Msg = map[string]string{"errMsg": "group_id 不能为空"}
		return
	}
	rediscache.JoinGroup(body.GroupId, uid)
	log.Info("[Chat] join group", "uid", uid, "group", body.GroupId)

	members := rediscache.GetGroupMembers(body.GroupId)

	// BroadcastAll 模式：先推本机组员（含自己，因为 JoinGroup 已写入 Redis），再广播其他节点
	notice := map[string]any{
		"type":     "user_joined",
		"groupId":  body.GroupId,
		"uid":      uid,
		"username": username,
		"members":  members,
	}
	net.ConnManager.PushToLocalGroupMembers(body.GroupId, "group.notice", notice)
	rediscache.PublishBroadcast(rediscache.BroadcastMsg{
		Action: "group", Target: body.GroupId, Name: "group.notice", Data: notice,
	})

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
	uid, username := net.GetAuthInfoFromReq(req)
	if uid == "" {
		rsp.Body.Code = -1
		rsp.Body.Msg = map[string]string{"errMsg": "请先登录"}
		return
	}

	if body.GroupId == "" {
		rsp.Body.Code = -1
		rsp.Body.Msg = map[string]string{"errMsg": "group_id 不能为空"}
		return
	}
	rediscache.LeaveGroup(body.GroupId, uid)
	log.Info("[Chat] leave group", "uid", uid, "group", body.GroupId)

	// BroadcastAll 模式：LeaveGroup 已从 Redis 移除自己，
	// PushToLocalGroupMembers 只会推给剩下的人（不会推给自己，正好符合预期）
	notice := map[string]any{
		"type":     "user_left",
		"groupId":  body.GroupId,
		"uid":      uid,
		"username": username,
		"members":  rediscache.GetGroupMembers(body.GroupId),
	}
	net.ConnManager.PushToLocalGroupMembers(body.GroupId, "group.notice", notice)
	rediscache.PublishBroadcast(rediscache.BroadcastMsg{
		Action: "group", Target: body.GroupId, Name: "group.notice", Data: notice,
	})

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
	if body.GroupId == "" {
		rsp.Body.Code = -1
		rsp.Body.Msg = map[string]string{"errMsg": "group_id 不能为空"}
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

// genMsgId 生成消息 ID（8 字节 = 16 位 hex）
func genMsgId() string { return util.RandomHex(8) }

func truncate(s string, n int) string {
	if len(s) <= n {
		return s
	}
	return s[:n] + "..."
}
