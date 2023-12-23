package model

// request
type CreateAccoutReq struct {
	Data CreateAccountReqData `json:"data"`
}

type CreateAccountReqData struct {
	Username        string `json:"username"`
	Password        string `json:"password"`
	ConfirmPassword string `json:"confirm_password"`
	Title           string `json:"title"`
	FirstName       string `json:"first_name"`
	LastName        string `json:"last_name"`
	Age             int    `json:"age"`
	MobileNO        string `json:"mobile_no"`
	Role            string `json:"role"`
}

// response
type CreateAccoutResp struct {
	Header ResponseHeader        `json:"header"`
	Data   CreateAccountRespData `json:"data"`
}

type CreateAccountRespData struct{}
