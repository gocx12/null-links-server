package logic

import (
	"context"
	"fmt"

	"null-links/rpc_service/webset/internal/svc"
	"null-links/rpc_service/webset/pb/webset"

	"github.com/demdxx/gocast"
	"github.com/zeromicro/go-zero/core/logx"
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
	MapUserWebsetLiked = "MAP_USER_LIKED"
	MapWebsetLikedCnt  = "MAP_WEBSET_LIKED_CNT"
)

func (l *LikeActionLogic) LikeAction(in *webset.LikeActionReq) (*webset.LikeActionResp, error) {
	// hash key: webset_id::user_id value:status
	var err error = nil
	likeActionResp := webset.LikeActionResp{
		StatusCode: 1,
		StatusMsg:  "success",
	}

	if in.ActionType == 1 {
		// 点赞
		key := gocast.ToString(in.WebsetId) + "::" + gocast.ToString(in.UserId)
		l.svcCtx.RedisClient.Hset(MapUserWebsetLiked, key, "1")
		// 点赞数+1
		l.svcCtx.RedisClient.Hincrby(MapWebsetLikedCnt, gocast.ToString(in.WebsetId), 1)

	} else if in.ActionType == 2 {
		// 取消点赞
		key := gocast.ToString(in.WebsetId) + "::" + gocast.ToString(in.UserId)
		l.svcCtx.RedisClient.Hset(MapUserWebsetLiked, key, "2")
		// 点赞数-1
		l.svcCtx.RedisClient.Hincrby(MapWebsetLikedCnt, gocast.ToString(in.WebsetId), -1)
	} else {
		// 未知操作类型
		logx.Error("unknown like action type")
		likeActionResp.StatusCode = 0
		likeActionResp.StatusMsg = "unknown like action type"
		err = fmt.Errorf("unknown like action type, action type: %d", in.ActionType)
	}

	return &likeActionResp, err
}
