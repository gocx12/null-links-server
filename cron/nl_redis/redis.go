package nl_redis

import (
	"github.com/zeromicro/go-zero/core/stores/redis"
)

type (
	NLRedis interface {
		Incr(key string) (int64, error)
	}
	NLRedisModel struct {
	}
)

func NewRedisClient(redisConfig redis.RedisConf) *redis.Redis {
	rds := redis.MustNewRedis(redisConfig)
	return rds
}
