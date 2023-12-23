package database

import (
	"context"
	"spider-go/config"
	"spider-go/logger"

	"github.com/go-redis/redis/v9"
)

var (
	RedisClient *redis.Client
)

func NewRedisClient(conf *config.RedisConfig) {
	log := logger.L().Named("Redis")

	option := redis.Options{
		Addr:     conf.HostPort,
		Password: conf.Password,
		DB:       conf.Index,
	}

	r := redis.NewClient(&option)

	if err := r.Ping(context.TODO()).Err(); err != nil {
		log.Errorf("Addr: %v, passowrd: %v, db_index: %v", conf.HostPort, conf.Password, conf.Index)
	}

	log.Info("redis clinet successful")

	RedisClient = r
}
