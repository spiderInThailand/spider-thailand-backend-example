package model

import "time"

// request
type LoginRequest struct {
	Data LoginRequestData `json:"data"`
}

type LoginRequestData struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

// response
type LoginResponse struct {
	Header ResponseHeader `json:"header"`
	Data   LoginInfo      `json:"data"`
}

type LoginInfo struct {
	User         User         `json:"user"`
	BackendToken BackendToken `json:"backendToken"`
}

type BackendToken struct {
	Token string `json:"token"`
}
type User struct {
	ID        string    `json:"id"`
	Username  string    `json:"username"`
	Name      string    `json:"name"`
	Role      string    `json:"role"`
	CreatedAt time.Time `json:"createdAt"`
}

// verify login
type VerifyLoginRequester struct {
	Header RequestUserHeader `json:"header"`
}

type VerifyLoginResponser struct {
	Header ResponseHeader `json:"header"`
}
