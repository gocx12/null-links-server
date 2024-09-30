package cache

import (
	"github.com/redis/go-redis/v9"
)

func NewClient(redisConfig *redis.Options) *redis.Client {
	rdb := redis.NewClient(redisConfig)
	return rdb
}
