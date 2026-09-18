package proto

// ──── 世界频道 ────
type SendWorldReq struct {
	Content string `json:"content"`
}
type SendWorldResp struct {
	MsgId string `json:"msg_id"`
}
type WorldChatMsg struct {
	MsgId     string `json:"msg_id"`
	FromUid   string `json:"from_uid"`
	FromName  string `json:"from_name"`
	Content   string `json:"content"`
	Timestamp int64  `json:"timestamp"`
}

// ──── 私聊 ────
type SendPrivateReq struct {
	ToUid   string `json:"to_uid"`
	Content string `json:"content"`
}
type SendPrivateResp struct {
	MsgId   string `json:"msg_id"`
	Offline bool   `json:"offline"`
}
type PrivateChatMsg struct {
	MsgId     string `json:"msg_id"`
	FromUid   string `json:"from_uid"`
	FromName  string `json:"from_name"`
	Content   string `json:"content"`
	Timestamp int64  `json:"timestamp"`
}

// ──── 群组 ────
type SendGroupReq struct {
	GroupId string `json:"group_id"`
	Content string `json:"content"`
}
type GroupChatMsg struct {
	MsgId     string `json:"msg_id"`
	GroupId   string `json:"group_id"`
	FromUid   string `json:"from_uid"`
	FromName  string `json:"from_name"`
	Content   string `json:"content"`
	Timestamp int64  `json:"timestamp"`
}
type GroupMemberReq struct {
	GroupId string `json:"group_id"`
}
type GroupMemberResp struct {
	GroupId string   `json:"group_id"`
	Members []string `json:"members"`
	Count   int      `json:"count"`
}

// ──── 在线用户 ────
type OnlineUsersResp struct {
	Total int      `json:"total"`
	Uids  []string `json:"uids"`
}
