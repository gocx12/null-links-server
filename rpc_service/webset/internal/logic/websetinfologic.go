package logic

import (
	"context"

	"null-links/rpc_service/webset/internal/svc"
	"null-links/rpc_service/webset/pb/webset"

	"github.com/demdxx/gocast"
	"github.com/zeromicro/go-zero/core/logx"
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
		l.Logger.Error("find webset failed, err: ", err)
		return &websetInfoResp, err
	}

	// 点赞信息, redis优先
	websetIdList := []int64{in.WebsetId}
	likeInfoRds, err := l.svcCtx.RedisClient.Hget(RdsKeyUserWebsetLiked, gocast.ToString(in.WebsetId)+"::"+gocast.ToString(in.UserId))
	if err != nil {
		l.Logger.Error("get like info from redis failed, err: ", err)
	}

	likeInfoDb, err = l.svcCtx.LikeModel.GetLikeWebsetUserInfos(l.ctx, websetIdList, in.UserId)
	if err != nil {
		l.Logger.Error("get like info failed, err: ", err)
	}

	// 收藏信息， redis优先
	favoriteInfoDb, err = l.svcCtx.FavoriteModel.GetFavoriteWebsetUserInfos(l.ctx, websetIdList, in.UserId)
	if err != nil {
		l.Logger.Error("get favorite info failed, user, err: ", err)
	}

	websetInfoResp.Webset = &webset.Webset{
		Id:       WebsetDb.Id,
		Title:    WebsetDb.Title,
		Describe: WebsetDb.Describe,
		AuthorInfo: &webset.UserInfo{
			Id: WebsetDb.AuthorId,
		},
		CoverUrl:      WebsetDb.CoverUrl,
		ViewCount:     WebsetDb.ViewCnt,
		LikeCount:     WebsetDb.LikeCnt,
		FavoriteCount: WebsetDb.FavoriteCnt,
		IsLike:        false,
		IsFavorite:    false,
	}

	return &websetInfoResp, nil
}
