package logic

import (
	"context"

	"null-links/rpc_service/webset/internal/svc"
	"null-links/rpc_service/webset/pb/webset"

	"github.com/demdxx/gocast"
	"github.com/redis/go-redis/v9"
	"github.com/zeromicro/go-zero/core/logx"
	"null-links/internal"
)

type LikeActionLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewLikeActionLogic(ctx context.Context, svcCtx *svc.ServiceContext) *LikeActionLogic {
	return &LikeActionLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

var (
	RdsKeyUserWebsetLiked = "HASH_WEBESET_USER_LIKED"
	RdsKeyWebsetLikedCnt  = "HASH_WEBSET_LIKED_CNT"
)

func (l *LikeActionLogic) LikeAction(in *webset.LikeActionReq) (*webset.LikeActionResp, error) {
	// hash key: webset_id::user_id value:status
	likeActionResp := webset.LikeActionResp{
		StatusCode: internal.StatusSuccess,
		StatusMsg:  "success",
	}

	if in.ActionType == 1 {
		// 点赞
		// redis事务
		_, err := l.svcCtx.RedisClient.TxPipelined(l.ctx, func(pipe redis.Pipeliner) error {
			key := gocast.ToString(in.WebsetId) + "::" + gocast.ToString(in.UserId)
			pipe.HSet(l.ctx, RdsKeyUserWebsetLiked, key, 1)
			// 点赞数+1
			pipe.HIncrBy(l.ctx, RdsKeyWebsetLikedCnt, gocast.ToString(in.WebsetId), 1)
			return nil
		})
		if err != nil {
			logx.Error("webset like, redis pipeline err: ", err)
			likeActionResp.StatusCode = internal.StatusRpcErr
			likeActionResp.StatusMsg = "redis pipeline err"
		}
	} else if in.ActionType == 2 {
		// 取消点赞
		// redis事务
		_, err := l.svcCtx.RedisClient.TxPipelined(l.ctx, func(pipe redis.Pipeliner) error {
			key := gocast.ToString(in.WebsetId) + "::" + gocast.ToString(in.UserId)
			pipe.HSet(l.ctx, RdsKeyUserWebsetLiked, key, 2)
			// 点赞数-1
			pipe.HIncrBy(l.ctx, RdsKeyWebsetLikedCnt, gocast.ToString(in.WebsetId), -1)
			return nil
		})
		if err != nil {
			logx.Error("cancel webset like, redis pipeline err: ", err)
			likeActionResp.StatusCode = internal.StatusRpcErr
			likeActionResp.StatusMsg = "redis pipeline err"
		}
	} else {
		// 未知操作类型
		logx.Error("unknown like action type")
		likeActionResp.StatusCode = internal.StatusRpcErr
		likeActionResp.StatusMsg = "unknown like action type"
	}

	return &likeActionResp, nil
}
