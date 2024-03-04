package nl_redis

import (
	"sync"

	"github.com/redis/go-redis/v9"
)

var once sync.Once
var Rdb *redis.Client

func NewClient(redisConfig *redis.Options) *redis.Client {
	once.Do(func() {
		Rdb = redis.NewClient(redisConfig)
	})
	return Rdb
}
