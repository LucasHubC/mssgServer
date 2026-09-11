package proto

type LoginReq struct {
	Username string `json:"username"`
	Password string `json:"password"`

	IP       string `json:"ip"`
	Hardware string `json:"hardware"`
}

type LoginResp struct {
	Username string `json:"username"`
	Session  string `json:"session"`
	UId      string `json:"uid"`
}

// ──── 注册 ────

type RegisterReq struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type RegisterResp struct {
	Username string `json:"username"`
	UId      string `json:"uid"`
}
