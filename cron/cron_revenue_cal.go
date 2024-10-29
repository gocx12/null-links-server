package main

import (
	"context"
	"fmt"
	"null-links/cron/model"
	"null-links/cron/nl_redis"
	"strings"

	"github.com/demdxx/gocast"
	"github.com/redis/go-redis/v9"
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

func RevenueCalJob(ctx context.Context) (int, error) {
	rdb := nl_redis.NewRedisClient(redisConf)
	if rdb == nil {
		return 0, fmt.Errorf("cannot create redis client")
	}

	newTLikeModel := model.NewTLikeModel(sqlx.NewMysql(mysqlDataSource))
	if newTLikeModel == nil {
		return 0, fmt.Errorf("cannot create TLikeModel mysql")
	}
	newTWebsetModel := model.NewTWebsetModel(sqlx.NewMysql(mysqlDataSource))
	if newTWebsetModel == nil {
		return 0, fmt.Errorf("cannot create TWebsetModel mysql")
	}

	// 1. 获取当前所有redis点赞记录
	rdsKeyWebsetUserLiked := "HASH_WEBESET_USER_LIKED"
	rdsKeyWebsetLikedCnt := "HASH_WEBSET_LIKED_CNT"

	websetLikeRdsKeys, err := rdb.HKeys(ctx, rdsKeyWebsetUserLiked).Result()
	if err != nil {
		logx.Error("redis Hkeys error: ", err)
		return 0, nil
	}

	for _, key := range websetLikeRdsKeys {
		logx.Debug("webset like key: ", key)
		// format: webset_id::user_id
		// 开启redis事务
		_, err := rdb.Pipelined(ctx, func(pipe redis.Pipeliner) error {
			value, err := pipe.HGet(ctx, rdsKeyWebsetUserLiked, key).Result()
			if err != nil {
				logx.Error("HGet error: ", err)
				return err
			}
			arr := strings.Split(key, "::")
			if len(arr) != 2 {
				logx.Error("invalid key: ", key)
				return fmt.Errorf("invalid key: %s", key)
			}
			websetId, userId := arr[0], arr[1]
			likeStatus := value

			// 2. 插入mysql点赞记录
			_, err = newTLikeModel.Insert(ctx, &model.TLike{
				WebsetId: gocast.ToInt64(websetId),
				UserId:   gocast.ToInt64(userId),
				Status:   gocast.ToInt64(likeStatus),
			})
			if err != nil {
				logx.Error("insert like, error: ", err, " ,websetId: ", websetId, " ,userId: ", userId)
				return err
			}

			// 3. 更新mysql点赞数
			likeCnt, err := pipe.HGet(ctx, rdsKeyWebsetLikedCnt, websetId).Result()
			if err != nil {
				logx.Error("HGet error: ", err)
				return err
			}
			err = newTWebsetModel.UpdateLikeCnt(ctx, gocast.ToInt64(likeCnt), gocast.ToInt64(websetId))
			if err != nil {
				logx.Error("update like count, error: ", err, " ,websetId: ", websetId)
				return err
			}

			// 4. 删除redis中的点赞记录
			_, err = pipe.HDel(ctx, rdsKeyWebsetUserLiked, key).Result()
			if err != nil {
				// TODO(chancyGao): 告警
				logx.Error("Hdel error: ", err)
				return err
			}

			// 5. 清空点赞数
			pipe.HSet(ctx, rdsKeyWebsetLikedCnt, websetId, 0)
			return nil

		})
		if err != nil {
			logx.Error("redis pipeline error: ", err)
			return 0, err
		}
	}

	return 0, nil
}
