package model

// request
type GetRsaKeyRequest struct {
	Data RsaKeyReqData `json:"data"`
}

type RsaKeyReqData struct{}

// response
type GetRsaKeyResponse struct {
	Header ResponseHeader `json:"header"`
	Data   RsaKeyRespData `json:"data"`
}

type RsaKeyRespData struct {
	PublicKey string `json:"public_key"`
	SearchKey string `json:"search_key"`
}
