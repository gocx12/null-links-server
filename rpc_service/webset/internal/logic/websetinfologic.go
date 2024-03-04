package logic

import (
	"context"

	"null-links/rpc_service/webset/internal/svc"
	"null-links/rpc_service/webset/pb/webset"

	"null-links/internal"

	"github.com/demdxx/gocast"
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

type WebsetInfoLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewWebsetInfoLogic(ctx context.Context, svcCtx *svc.ServiceContext) *WebsetInfoLogic {
	return &WebsetInfoLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *WebsetInfoLogic) WebsetInfo(in *webset.WebsetInfoReq) (*webset.WebsetInfoResp, error) {
	websetInfoResp := webset.WebsetInfoResp{
		StatusCode: 1,
		StatusMsg:  "success",
		Webset:     nil,
	}

	WebsetDb, err := l.svcCtx.WebsetModel.FindOne(l.ctx, in.WebsetId)
	if err != nil {
		logx.Error("get webset info from db failed, err: ", err)
		websetInfoResp.StatusCode = internal.StatusRpcErr
		websetInfoResp.StatusMsg = "get webset info from db failed"
		return &websetInfoResp, err
	}

	// 点赞信息, redis优先
	isLikeResp := false
	var likeCnt int64 = 0
	isLikeRds, err := l.svcCtx.RedisClient.HGet(l.ctx, RdsKeyUserWebsetLiked, gocast.ToString(in.WebsetId)+"::"+gocast.ToString(in.UserId)).Result()
	if err != nil {
		logx.Error("get like info from redis failed, err: ", err)
		// TODO(chancyGao): 验证当redis不存在这个key时
		likeInfoDb, err := l.svcCtx.LikeModel.GetLikeWebsetUserInfo(l.ctx, in.WebsetId, in.UserId)
		if err != nil && err != sqlx.ErrNotFound {
			logx.Error("get like info failed, err: ", err)
		} else if err == nil && likeInfoDb.Status == 1 {
			isLikeResp = true
		}
	} else if isLikeRds == "1" {
		isLikeResp = true
	}

	likeCntRds, err := l.svcCtx.RedisClient.HGet(l.ctx, RdsKeyWebsetLikedCnt, gocast.ToString(in.WebsetId)).Result()
	if err != nil {
		logx.Error("get like count from redis failed, err: ", err)
	} else {
		likeCnt = gocast.ToInt64(likeCntRds) + WebsetDb.LikeCnt
	}

	// 收藏信息， redis优先
	isFavoriteResp := false
	// favoriteInfoDb, err = l.svcCtx.FavoriteModel.GetFavoriteWebsetUserInfos(l.ctx, websetIdList, in.UserId)
	// if err != nil {
	// 	l.Logger.Error("get favorite info failed, user, err: ", err)
	// }

	websetInfoResp.Webset = &webset.Webset{
		Id:       WebsetDb.Id,
		Title:    WebsetDb.Title,
		Describe: WebsetDb.Describe,
		AuthorInfo: &webset.UserInfo{
			Id: WebsetDb.AuthorId,
		},
		CoverUrl:      WebsetDb.CoverUrl,
		ViewCount:     WebsetDb.ViewCnt,
		LikeCount:     likeCnt,
		FavoriteCount: WebsetDb.FavoriteCnt,
		IsLike:        isLikeResp,
		IsFavorite:    isFavoriteResp,
	}

	return &websetInfoResp, nil
}
