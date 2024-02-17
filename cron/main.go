package main

import (
	"context"
	"fmt"
	"null-links/cron/model"
	"null-links/cron/nl_redis"
	"strings"

	"github.com/demdxx/gocast"
	"github.com/reugn/go-quartz/job"
	"github.com/reugn/go-quartz/quartz"
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/core/stores/redis"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

var redisConf = redis.RedisConf{
	Host: "127.0.0.1:6379",
	Type: "node",
}
var mysqlDataSource = "root:123456@tcp(127.0.0.1:3306)/db_null_links?charset=utf8mb4&parseTime=true"

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// create scheduler
	sched := quartz.NewStdScheduler()

	// async start scheduler
	sched.Start(ctx)

	// create jobs
	cronTrigger, _ := quartz.NewCronTrigger("1/5 * * * * *")

	likeFuncJob := job.NewFunctionJob(func(ctx context.Context) (int, error) {
		// redisConf := redis.RedisConf{
		// 	Host: "127.0.0.1:6379",
		// 	Type: "node",
		// }
		rds := nl_redis.NewRedisClient(redisConf)
		if rds == nil {
			logx.Error("rds is nil")
			return 0, nil
		}

		userLikedRdsKey := "MAP_USER_LIKED"
		WebsetUserLIkeMap, err := rds.Hgetall(userLikedRdsKey)
		if err != nil {
			logx.Error("Hgetall error: ", err)
			return 0, nil
		}

		for k, v := range WebsetUserLIkeMap {
			fmt.Printf("k: %v, v: %v\n", k, v)
			// format: webset_id::user_id
			arr := strings.Split(k, "::")
			websetId, userId := arr[0], arr[1]
			likeStatus := v
			if len(arr) != 2 {
				logx.Error("invalid key: ", k)
				continue
			}
			newTLikeModel := model.NewTLikeModel(sqlx.NewMysql(mysqlDataSource))
			_, err := newTLikeModel.Insert(ctx, &model.TLike{
				WebsetId: gocast.ToInt64(websetId),
				UserId:   0,
				Status:   gocast.ToInt64(likeStatus),
				// CreatedAt: nil,
				// UpdatedAt: nil,
				// DeletedAt: nil,
			})
			if err != nil {
				logx.Error("insert like, error: ", err, " ,websetId: ", websetId, " ,userId: ", userId)
				continue
			}

			_, err = rds.Hdel(userLikedRdsKey, k)
			if err != nil {
				logx.Error("Hdel error: ", err)
				continue
			}

		}

		return 0, nil
	})

	// register jobs to scheduler
	sched.ScheduleJob(quartz.NewJobDetail(likeFuncJob, quartz.NewJobKey("likeFuncJob")),
		cronTrigger)

	// // stop scheduler
	// sched.Stop()

	// wait for all workers to exit
	sched.Wait(ctx)
}
