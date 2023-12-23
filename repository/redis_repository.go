package repository

import (
	"context"
	"encoding/json"
	"spider-go/domain"
	"spider-go/logger"
	"time"

	"github.com/go-redis/redis/v9"
)

type RedisRepository struct {
	client *redis.Client
	log    *logger.Logger
}

func NewRedisRepository(redisCline *redis.Client) domain.RedisRepository {
	return &RedisRepository{
		client: redisCline,
		log:    logger.L().Named("RedisRepository"),
	}
}

func (r *RedisRepository) SetDataToRedisWithTTL(ctx context.Context, key string, data interface{}, ttl time.Duration) (err error) {

	prepareJsonData, err := json.Marshal(data)
	if err != nil {
		return err
	}

	return r.client.Set(ctx, key, prepareJsonData, ttl).Err()
}

func (r *RedisRepository) GetDataFromRedis(ctx context.Context, key string, data interface{}) (err error) {

	result, err := r.client.Get(ctx, key).Bytes()
	if err != nil {
		if err == redis.Nil {
			return ErrorRedisNotFound
		}
		return err
	}

	json.Unmarshal(result, data)
	return nil
}
