package proto

type LoginReq struct {
	Username string `json:"username"`
	Password string `json:"password"`

	IP       string `json:"ip"`
	Hardware string `json:"hardware"`
}

type LoginResp struct {
	Username string `json:"username"`
	Password string `json:"password"`

	Session string `json:"session"`
	UId     string `json:"uid"`
}
