package nl_redis

import (
	"github.com/redis/go-redis/v9"
	"github.com/zeromicro/go-zero/core/logx"
)

func NewRedisClient(redisConfig redis.Options) *redis.Client {
	rdb := redis.NewClient(&redisConfig)

	if rdb == nil {
		// 告警
		logx.Error("cannot create redis client")
	}

	return rdb
}
