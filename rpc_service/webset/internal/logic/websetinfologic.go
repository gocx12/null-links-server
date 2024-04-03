package logic

import (
	"context"

	"null-links/rpc_service/user/pb/user"
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
	websetInfoResp := webset.WebsetInfoResp{}

	WebsetDb, err := l.svcCtx.WebsetModel.FindOne(l.ctx, in.WebsetId)
	if err != nil {
		logx.Error("get webset info from db failed, err: ", err)
		websetInfoResp.StatusCode = internal.StatusRpcErr
		websetInfoResp.StatusMsg = "get webset info from db failed"
		return &websetInfoResp, err
	}

	// 获取是否点赞信息, 如果用户未登录，则默认为未点赞
	isLike := false
	getLikeInfoFromBF := false
	if in.UserId != -1 {
		// 检查该webset对应的布隆过滤器是否存在
		res, err := l.svcCtx.RedisClient.Exists(l.ctx, "BF_LIKE_"+gocast.ToString(in.WebsetId)).Result()
		if err != nil {
			logx.Error("check whether BF_LIKE_ exists from redis error: ", err)
		} else if res == 1 {
			// 布隆过滤器存在
			res, err := l.svcCtx.RedisClient.BFExists(l.ctx, "BF_LIKE_"+gocast.ToString(in.WebsetId), in.UserId).Result()
			if err != nil {
				logx.Error("get like info from redis bloom filter error: ", err)
			} else if !res {
				// 用户未点赞
				logx.Debug("user not liked")
				getLikeInfoFromBF = true
			}
		}
		if !getLikeInfoFromBF {
			likeInfoDb, err := l.svcCtx.LikeModel.GetLikeWebsetUserInfo(l.ctx, in.WebsetId, in.UserId)
			if err != nil && err != sqlx.ErrNotFound {
				logx.Error("get like info error: ", err)
			} else if err == sqlx.ErrNotFound {
				isLike = false
			} else if likeInfoDb.Status == 1 {
				isLike = true
			} else {
				isLike = false
			}
		}
	}

	// 获取是否收藏信息
	isFavoriteResp := false
	// favoriteInfoDb, err = l.svcCtx.FavoriteModel.GetFavoriteWebsetUserInfos(l.ctx, websetIdList, in.UserId)
	// if err != nil {
	// 	l.Logger.Error("get favorite info failed, user, err: ", err)
	// }

	// 获取作者信息
	userInfoRpcRep, err := l.svcCtx.UserRpc.UserInfo(l.ctx, &user.UserInfoReq{
		UserId: WebsetDb.AuthorId,
	})
	if err != nil {
		logx.Error("get user info from rpc error: ", err)
		websetInfoResp.StatusCode = internal.StatusRpcErr
		websetInfoResp.StatusMsg = "get user info from rpc error"
		return &websetInfoResp, nil
	}

	// 获取weblink信息
	weblinksDb, err := l.svcCtx.WeblinkModel.FindByWebsetId(l.ctx, in.WebsetId)
	if err != nil {
		logx.Error("get user info from rpc error: ", err)
		websetInfoResp.StatusCode = internal.StatusRpcErr
		websetInfoResp.StatusMsg = "get weblink from rpc error"
		return &websetInfoResp, nil
	}

	weblinkListResp := make([]*webset.WebLink, 0, len(weblinksDb))
	for _, weblink := range weblinksDb {
		weblinkResp := &webset.WebLink{
			Id:       weblink.Id,
			Describe: weblink.Describe,
			Url:      weblink.Url,
			CoverUrl: weblink.CoverUrl,
		}
		weblinkListResp = append(weblinkListResp, weblinkResp)
	}

	websetInfoResp.StatusCode = internal.StatusSuccess
	websetInfoResp.StatusMsg = "success"
	websetInfoResp.Webset = &webset.Webset{
		Id:       WebsetDb.Id,
		Title:    WebsetDb.Title,
		Describe: WebsetDb.Describe,
		AuthorInfo: &webset.UserInfo{
			Id:        WebsetDb.AuthorId,
			Name:      userInfoRpcRep.UserInfo.Name,
			AvatarUrl: userInfoRpcRep.UserInfo.AvatarUrl,
		},
		CoverUrl:      WebsetDb.CoverUrl,
		ViewCount:     WebsetDb.ViewCnt,
		LikeCount:     WebsetDb.LikeCnt,
		FavoriteCount: WebsetDb.FavoriteCnt,
		IsLike:        isLike,
		IsFavorite:    isFavoriteResp,
		WebLinkList:   weblinkListResp,
	}

	return &websetInfoResp, nil
}
