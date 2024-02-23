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

func FavoriteJob(ctx context.Context) (int, error) {
	rdb := nl_redis.NewRedisClient(redisConf)
	if rdb == nil {
		return 0, fmt.Errorf("cannot create redis client")
	}

	newTFavoriteModel := model.NewTFavoriteModel(sqlx.NewMysql(mysqlDataSource))
	if newTFavoriteModel == nil {
		return 0, fmt.Errorf("cannot create TFavoriteModel mysql")
	}
	newTWebsetModel := model.NewTWebsetModel(sqlx.NewMysql(mysqlDataSource))
	if newTWebsetModel == nil {
		return 0, fmt.Errorf("cannot create TWebsetModel mysql")
	}

	// 1. 获取当前所有redis收藏记录
	rdsKeyWebsetUserFavorited := "HASH_WEBESET_USER_FAVORITED"
	rdsKeyWebsetFavoritedCnt := "HASH_WEBSET_FAVORITED_CNT"

	websetLikeRdsKeys, err := rdb.HKeys(ctx, rdsKeyWebsetUserFavorited).Result()
	if err != nil {
		logx.Error("redis Hkeys error: ", err)
		return 0, nil
	}

	for _, key := range websetLikeRdsKeys {
		logx.Debug("webset favorite key: ", key)
		// format: webset_id::user_id
		// 开启redis事务
		_, err := rdb.Pipelined(ctx, func(pipe redis.Pipeliner) error {
			value, err := pipe.HGet(ctx, rdsKeyWebsetUserFavorited, key).Result()
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
			favoriteStatus := value

			// 2. 插入mysql收藏记录
			_, err = newTFavoriteModel.Insert(ctx, &model.TFavorite{
				WebsetId: gocast.ToInt64(websetId),
				UserId:   gocast.ToInt64(userId),
				Status:   gocast.ToInt64(favoriteStatus),
			})
			if err != nil {
				logx.Error("insert favorite, error: ", err, ", websetId: ", websetId, ", userId: ", userId)
				return err
			}

			// 3. 更新mysql收藏数
			favoriteCnt, err := pipe.HGet(ctx, rdsKeyWebsetFavoritedCnt, websetId).Result()
			if err != nil {
				logx.Error("HGet error: ", err)
				return err
			}
			err = newTWebsetModel.UpdateFavoriteCnt(ctx, gocast.ToInt64(favoriteCnt), gocast.ToInt64(websetId))
			if err != nil {
				logx.Error("update favorite count, error: ", err, ", websetId: ", websetId)
				return err
			}

			// 4. 删除redis中的收藏记录
			_, err = pipe.HDel(ctx, rdsKeyWebsetUserFavorited, key).Result()
			if err != nil {
				// TODO(chancyGao): 告警
				logx.Error("Hdel error: ", err)
				return err
			}

			// 5. 清空收藏数
			pipe.HSet(ctx, rdsKeyWebsetFavoritedCnt, websetId, 0)
			return nil

		})
		if err != nil {
			logx.Error("redis pipeline error: ", err)
			return 0, err
		}
	}
	return 0, nil
}
