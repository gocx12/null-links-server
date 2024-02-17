package main

import (
	"context"
	"fmt"
	"null-links/cron/nl_redis"

	"github.com/zeromicro/go-zero/core/stores/redis"
)

func LikeJob(_ context.Context) (int, error) {
	redisConf := redis.RedisConf{
		Host: "127.0.0.1:6438",
		Type: "node",
		Pass: "1334",
	}
	rds := nl_redis.NewRedisClient(redisConf)

	rds.Hget("MAP_USER_LIKED", "like")

	vals, err := rds.Hvals("MAP_USER_LIKED")
	fmt.Printf("vals: %v, err: %v\n", vals, err)

	var cursor uint64 = 0
	// keys []string, cur uint64, err error
	keys, cur, err := rds.Hscan("MAP_USER_LIKED", cursor, "", 10)
	fmt.Printf("keys: %v, cur: %v, err: %v\n", keys, cur, err)

	return 0, nil
}
