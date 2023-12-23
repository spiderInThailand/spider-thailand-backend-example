package usecase

import "fmt"

var (
	ErrorValidateUserFaild  = fmt.Errorf("validate user login failed")
	ErrorRedisNotFound      = fmt.Errorf("redis not found")
	ErrorRedisConnection    = fmt.Errorf("redis connection failed")
	ErrorMongoConnection    = fmt.Errorf("mongoDB connection failed")
	ErrorMongoTechnicalFail = fmt.Errorf("error mongo tech")
	ErrorTechnicalError     = fmt.Errorf("technical error")
)
