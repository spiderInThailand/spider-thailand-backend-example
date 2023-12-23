package model

type RequestUserHeader struct {
	Username string `json:"username"`
	Token    string `json:"token"`
}
type ResponseHeader struct {
	ErrorCode string `json:"error_code"`
	Message   string `json:"message"`
}
