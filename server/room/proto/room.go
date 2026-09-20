package proto

// ──── 房间管理请求/响应 ────

type CreateRoomReq struct {
	Name    string `json:"name"`
	MaxSize int    `json:"max_size"` // 0 = 默认 10
}

type CreateRoomResp struct {
	RoomId string `json:"room_id"`
	Owner  string `json:"owner"`
	Name   string `json:"name"`
}

type JoinRoomReq struct {
	RoomId string `json:"room_id"`
}

type JoinRoomResp struct {
	RoomId  string   `json:"room_id"`
	Name    string   `json:"name"`
	Owner   string   `json:"owner"`
	Size    int      `json:"size"`
	Max     int      `json:"max"`
	Members []string `json:"members"`
}

type KickRoomReq struct {
	RoomId    string `json:"room_id"`
	TargetUid string `json:"target_uid"`
}

type DeleteRoomReq struct {
	RoomId string `json:"room_id"`
}

type LeaveRoomReq struct {
	RoomId string `json:"room_id"`
}

// ──── 房间聊天 ────

type RoomSendReq struct {
	RoomId  string `json:"room_id"`
	Content string `json:"content"`
}

type RoomChatMsg struct {
	MsgId     string `json:"msg_id"`
	RoomId    string `json:"room_id"`
	FromUid   string `json:"from_uid"`
	FromName  string `json:"from_name"`
	Content   string `json:"content"`
	Timestamp int64  `json:"timestamp"`
}
