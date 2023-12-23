package model

import "time"

type RedisRsaKey struct {
	PrivateKey string `json:"private_key" bson:"private_key"`
	PublicKey  string `json:"public_key" bson:"public_key"`
}

type Login struct {
	Username  string    `json:"username"`
	Role      string    `json:"role"`
	CreatedAt time.Time `json:"created_at"`
}
