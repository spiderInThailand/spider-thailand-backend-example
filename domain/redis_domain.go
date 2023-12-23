package domain

import (
	"context"
	"time"
)

//go:generate mockgen -source=redis_domain.go -destination=./mock/redis_domain.go

type RedisRepository interface {
	SetDataToRedisWithTTL(ctx context.Context, key string, data interface{}, ttl time.Duration) (err error)
	GetDataFromRedis(ctx context.Context, key string, data interface{}) (err error)
}
