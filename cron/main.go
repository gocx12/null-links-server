package main

import (
	"context"

	"github.com/reugn/go-quartz/job"
	"github.com/reugn/go-quartz/quartz"

	"github.com/redis/go-redis/v9"
)

var redisConf = redis.Options{
	Addr:     "10.63.180.57:6379",
	Password: "123456",
}
var mysqlDataSource = "root:123456@tcp(10.63.180.57:3306)/db_null_links?charset=utf8mb4&parseTime=true"

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// create scheduler
	sched := quartz.NewStdScheduler()

	// async start scheduler
	sched.Start(ctx)

	// create jobs
	cronTrigger, _ := quartz.NewCronTrigger("1/5 * * * * *")
	likeFuncJob := job.NewFunctionJob(LikeJob)
	favoriteFuncJob := job.NewFunctionJob(FavoriteJob)

	// register jobs to scheduler
	sched.ScheduleJob(quartz.NewJobDetail(likeFuncJob, quartz.NewJobKey("likeFuncJob")),
		cronTrigger)
	sched.ScheduleJob(quartz.NewJobDetail(favoriteFuncJob, quartz.NewJobKey("favoriteFuncJob")),
		cronTrigger)

	// stop scheduler
	// sched.Stop()

	// wait for all workers to exit
	sched.Wait(ctx)
}
