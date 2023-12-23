package repository

import "fmt"

// error
var ErrorMongoNotFound = fmt.Errorf("mongo not found")
var ErrorRedisNotFound = fmt.Errorf("redis not found")
