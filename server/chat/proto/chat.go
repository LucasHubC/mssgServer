package proto

// ──── 世界频道 ────
type SendWorldReq struct {
	Content string `json:"content"`
}
type SendWorldResp struct {
	MsgId string `json:"msgId"`
}
type WorldChatMsg struct {
	MsgId     string `json:"msgId"`
	FromUid   string `json:"fromUid"`
	FromName  string `json:"fromName"`
	Content   string `json:"content"`
	Timestamp int64  `json:"timestamp"`
}

// ──── 私聊 ────
type SendPrivateReq struct {
	ToUid   string `json:"toUid"`
	Content string `json:"content"`
}
type SendPrivateResp struct {
	MsgId   string `json:"msgId"`
	Offline bool   `json:"offline"`
}
type PrivateChatMsg struct {
	MsgId     string `json:"msgId"`
	FromUid   string `json:"fromUid"`
	FromName  string `json:"fromName"`
	Content   string `json:"content"`
	Timestamp int64  `json:"timestamp"`
}

// ──── 群组 ────
type SendGroupReq struct {
	GroupId string `json:"groupId"`
	Content string `json:"content"`
}
type GroupChatMsg struct {
	MsgId     string `json:"msgId"`
	GroupId   string `json:"groupId"`
	FromUid   string `json:"fromUid"`
	FromName  string `json:"fromName"`
	Content   string `json:"content"`
	Timestamp int64  `json:"timestamp"`
}
type GroupMemberReq struct {
	GroupId string `json:"groupId"`
}
type GroupMemberResp struct {
	GroupId string   `json:"groupId"`
	Members []string `json:"members"`
	Count   int      `json:"count"`
}

// ──── 在线用户 ────
type OnlineUsersResp struct {
	Total int      `json:"total"`
	Uids  []string `json:"uids"`
}
