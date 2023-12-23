package middleware

type MiddlewareResponse struct {
	Header ResponseHeader `json:"header"`
}

type ResponseHeader struct {
	ErrorCode string `json:"error_code"`
	Message   string `json:"message"`
}

type MiddlewareRequest struct {
	Header map[string]interface{} `json:"header"`
	Data   map[string]interface{} `json:"data"`
}
