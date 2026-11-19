package controller

import (
	"mssgServer/infra/log"
	rediscache "mssgServer/infra/redis"
	"mssgServer/infra/util"
	"mssgServer/net"
	"mssgServer/server/room/proto"
	"strings"
	"time"
)

var DefaultRoom = &Room{}

type Room struct{}

const (
	roomPrefix   = "rm_" // 房间 ID 前缀
	msgIdHexLen  = 4     // 短消息 ID: 4 字节 = 8 位 hex
	roomIdHexLen = 6     // 房间 ID: 6 字节 = 12 位 hex
)

// genRoomId 生成房间 ID: rm_ + 12 位 hex
func genRoomId() string { return roomPrefix + util.RandomHex(roomIdHexLen) }

// genMsgIdShort 短消息 ID
func genMsgIdShort() string { return util.RandomHex(msgIdHexLen) }

// Router 消息名统一前缀: room.xxx
func (r *Room) Router(rtr *net.Router) {
	g := rtr.Group("room")
	g.AddRoute("create", r.wsCreateRoom)
	g.AddRoute("join", r.wsJoinRoom)
	g.AddRoute("leave", r.wsLeaveRoom)
	g.AddRoute("kick", r.wsKickRoom)
	g.AddRoute("delete", r.wsDeleteRoom)
	g.AddRoute("send", r.wsRoomSend)
}

// ──── 房间管理 WS Handlers ────

// wsCreateRoom room.create —— 创建房间
func (r *Room) wsCreateRoom(req *net.WsMsgReq, rsp *net.WsMsgResp) {
	rsp.Body.Code = 0

	body := &proto.CreateRoomReq{}
	if err := net.ParseMsg(req.Body.Msg, body); err != nil {
		rsp.Body.Code = -1
		rsp.Body.Msg = map[string]string{"errMsg": "参数解析失败"}
		return
	}
	uid, _ := net.GetAuthInfoFromReq(req)
	if uid == "" {
		rsp.Body.Code = -1
		rsp.Body.Msg = map[string]string{"errMsg": "请先登录"}
		return
	}

	body.Name = strings.TrimSpace(body.Name)
	if body.Name == "" {
		rsp.Body.Code = -1
		rsp.Body.Msg = map[string]string{"errMsg": "房间名不能为空"}
		return
	}
	if body.MaxSize <= 0 {
		body.MaxSize = 10
	}
	if body.MaxSize > 100 {
		body.MaxSize = 100
	}

	roomId := genRoomId()
	if err := rediscache.CreateRoom(roomId, body.Name, uid, body.MaxSize); err != nil {
		rsp.Body.Code = -1
		rsp.Body.Msg = map[string]string{"errMsg": "创建房间失败"}
		return
	}
	log.Info("[WS] room created", "room", roomId, "owner", uid, "name", body.Name)

	rsp.Body.Msg = &proto.CreateRoomResp{RoomId: roomId, Owner: uid, Name: body.Name}

	if ws, ok := req.Conn.(*net.WSServer); ok {
		ws.Push("room.notice", map[string]any{
			"type":   "room_created",
			"roomId": roomId,
			"name":   body.Name,
		})
	}
}

// wsJoinRoom room.join —— 加入房间
func (r *Room) wsJoinRoom(req *net.WsMsgReq, rsp *net.WsMsgResp) {
	rsp.Body.Code = 0

	body := &proto.JoinRoomReq{}
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
	if body.RoomId == "" {
		rsp.Body.Code = -1
		rsp.Body.Msg = map[string]string{"errMsg": "room_id 不能为空"}
		return
	}

	if err := rediscache.JoinRoom(body.RoomId, uid); err != nil {
		rsp.Body.Code = -1
		msg := "加入房间失败"
		switch err {
		case rediscache.ErrRoomNotFound:
			msg = "房间不存在"
		case rediscache.ErrRoomFull:
			msg = "房间已满"
		}
		rsp.Body.Msg = map[string]string{"errMsg": msg}
		return
	}

	info, _ := rediscache.GetRoomInfo(body.RoomId)
	members := rediscache.GetRoomMembers(body.RoomId)
	resp := &proto.JoinRoomResp{RoomId: body.RoomId, Members: members}
	if info != nil {
		resp.Name = info.Name
		resp.Owner = info.Owner
		resp.Max = info.MaxSize
	}
	resp.Size = len(members)
	rsp.Body.Msg = resp

	notice := map[string]any{
		"type":     "user_joined",
		"roomId":   body.RoomId,
		"uid":      uid,
		"username": username,
		"members":  members,
	}
	net.ConnManager.PushToLocalRoomMembers(body.RoomId, "room.notice", notice)
	rediscache.PublishBroadcast(rediscache.BroadcastMsg{
		Action: "room", Target: body.RoomId, Name: "room.notice", Data: notice,
	})
	log.Info("[WS] user joined room", "room", body.RoomId, "uid", uid)
}

// wsLeaveRoom room.leave —— 离开房间
func (r *Room) wsLeaveRoom(req *net.WsMsgReq, rsp *net.WsMsgResp) {
	rsp.Body.Code = 0

	body := &proto.LeaveRoomReq{}
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
	if body.RoomId == "" {
		rsp.Body.Code = -1
		rsp.Body.Msg = map[string]string{"errMsg": "room_id 不能为空"}
		return
	}

	rediscache.LeaveRoom(body.RoomId, uid)

	notice := map[string]any{
		"type":     "user_left",
		"roomId":   body.RoomId,
		"uid":      uid,
		"username": username,
		"members":  rediscache.GetRoomMembers(body.RoomId),
	}
	net.ConnManager.PushToLocalRoomMembers(body.RoomId, "room.notice", notice)
	rediscache.PublishBroadcast(rediscache.BroadcastMsg{
		Action: "room", Target: body.RoomId, Name: "room.notice", Data: notice,
	})

	rsp.Body.Msg = map[string]string{"msg": "ok"}
	log.Info("[WS] user left room", "room", body.RoomId, "uid", uid)
}

// wsKickRoom room.kick —— 房主踢人
func (r *Room) wsKickRoom(req *net.WsMsgReq, rsp *net.WsMsgResp) {
	rsp.Body.Code = 0

	body := &proto.KickRoomReq{}
	if err := net.ParseMsg(req.Body.Msg, body); err != nil {
		rsp.Body.Code = -1
		rsp.Body.Msg = map[string]string{"errMsg": "参数解析失败"}
		return
	}
	uid, _ := net.GetAuthInfo(req.Conn)
	if uid == "" {
		rsp.Body.Code = -1
		rsp.Body.Msg = map[string]string{"errMsg": "请先登录"}
		return
	}

	if body.RoomId == "" || body.TargetUid == "" {
		rsp.Body.Code = -1
		rsp.Body.Msg = map[string]string{"errMsg": "room_id 和 target_uid 不能为空"}
		return
	}

	if err := rediscache.KickFromRoom(body.RoomId, uid, body.TargetUid); err != nil {
		rsp.Body.Code = -1
		msg := "踢出失败"
		switch err {
		case rediscache.ErrRoomNotFound:
			msg = "房间不存在"
		case rediscache.ErrNotRoomOwner:
			msg = "只有房主能踢人"
		case rediscache.ErrCannotKickSelf:
			msg = "不能踢自己"
		}
		rsp.Body.Msg = map[string]string{"errMsg": msg}
		return
	}

	net.ConnManager.PushToUid(body.TargetUid, "room.notice", map[string]any{
		"type":   "kicked",
		"roomId": body.RoomId,
	})

	notice := map[string]any{
		"type":    "user_kicked",
		"roomId":  body.RoomId,
		"target":  body.TargetUid,
		"members": rediscache.GetRoomMembers(body.RoomId),
	}
	net.ConnManager.PushToLocalRoomMembers(body.RoomId, "room.notice", notice)
	rediscache.PublishBroadcast(rediscache.BroadcastMsg{
		Action: "room", Target: body.RoomId, Name: "room.notice", Data: notice,
	})

	rsp.Body.Msg = map[string]string{"msg": "ok"}
	log.Info("[WS] kicked", "room", body.RoomId, "from", uid, "target", body.TargetUid)
}

// wsDeleteRoom room.delete —— 房主删除房间
func (r *Room) wsDeleteRoom(req *net.WsMsgReq, rsp *net.WsMsgResp) {
	rsp.Body.Code = 0

	body := &proto.DeleteRoomReq{}
	if err := net.ParseMsg(req.Body.Msg, body); err != nil {
		rsp.Body.Code = -1
		rsp.Body.Msg = map[string]string{"errMsg": "参数解析失败"}
		return
	}
	uid, _ := net.GetAuthInfo(req.Conn)
	if uid == "" {
		rsp.Body.Code = -1
		rsp.Body.Msg = map[string]string{"errMsg": "请先登录"}
		return
	}

	if body.RoomId == "" {
		rsp.Body.Code = -1
		rsp.Body.Msg = map[string]string{"errMsg": "room_id 不能为空"}
		return
	}

	// ① 先删（Lua 原子，一次性清 room info + members + 所有用户索引）
	if err := rediscache.DeleteRoom(body.RoomId, uid); err != nil {
		rsp.Body.Code = -1
		msg := "删除失败"
		switch err {
		case rediscache.ErrRoomNotFound:
			msg = "房间不存在"
		case rediscache.ErrNotRoomOwner:
			msg = "只有房主能删除房间"
		}
		rsp.Body.Msg = map[string]string{"errMsg": msg}
		return
	}

	// ② 再广播 —— 此时 Redis 已删干净，广播里只带 roomId 让客户端自己清缓存
	notice := map[string]any{
		"type":   "room_deleted",
		"roomId": body.RoomId,
	}
	net.ConnManager.PushToLocalRoomMembers(body.RoomId, "room.notice", notice)
	rediscache.PublishBroadcast(rediscache.BroadcastMsg{
		Action: "room", Target: body.RoomId, Name: "room.notice", Data: notice,
	})

	rsp.Body.Msg = map[string]string{"msg": "ok"}
	log.Info("[WS] room deleted", "room", body.RoomId, "by", uid)
}

// wsRoomSend room.send —— 房间内发消息（跨服广播）
func (r *Room) wsRoomSend(req *net.WsMsgReq, rsp *net.WsMsgResp) {
	rsp.Body.Code = 0

	body := &proto.RoomSendReq{}
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

	if body.RoomId == "" {
		rsp.Body.Code = -1
		rsp.Body.Msg = map[string]string{"errMsg": "room_id 不能为空"}
		return
	}
	body.Content = strings.TrimSpace(body.Content)
	if body.Content == "" || len(body.Content) > 500 {
		rsp.Body.Code = -1
		rsp.Body.Msg = map[string]string{"errMsg": "内容不能为空且不超过 500 字"}
		return
	}

	if !rediscache.IsRoomMember(body.RoomId, uid) {
		rsp.Body.Code = -1
		rsp.Body.Msg = map[string]string{"errMsg": "你不在这个房间里"}
		return
	}

	msg := proto.RoomChatMsg{
		MsgId:     genMsgIdShort(),
		RoomId:    body.RoomId,
		FromUid:   uid,
		FromName:  username,
		Content:   body.Content,
		Timestamp: time.Now().Unix(),
	}

	log.Info("[WS] room.send", "room", body.RoomId, "from", username, "msgId", msg.MsgId)

	net.ConnManager.PushToLocalRoomMembers(body.RoomId, "room.chat.new", msg)
	rediscache.PublishBroadcast(rediscache.BroadcastMsg{
		Action: "room",
		Target: body.RoomId,
		Name:   "room.chat.new",
		Data:   msg,
	})

	rsp.Body.Msg = map[string]string{"msgId": msg.MsgId}
}
