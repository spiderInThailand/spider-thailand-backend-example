package model

var (
	ACCOUNT_ROLE_MASTER  = "master"
	ACCOUNT_ROLE_ADMIN   = "admin"
	ACCOUNT_ROLE_GENERAL = "general"
)

type Account struct {
	Username     string `json:"username" bson:"username"`
	HashPassword string `json:"hash_password" bson:"hash_password"`
	Title        string `json:"title" bson:"title"`
	FirstName    string `json:"first_name" bson:"first_name"`
	LastName     string `json:"last_name" bson:"last_name"`
	Age          int    `json:"age" bson:"age"`
	MobileNO     string `json:"mobile_no" bson:"mobile_no"`
	Role         string `json:"role" bson:"role"`
}

type LoginUser struct {
	Username string `json:"username" bson:"username"`
	Role     string `json:"role" bson:"role"`
}
