package nl_redis

import (
	"github.com/redis/go-redis/v9"
)

func NewRedisClient(redisConfig redis.Options) *redis.Client {
	rdb := redis.NewClient(&redisConfig)
	return rdb
}
