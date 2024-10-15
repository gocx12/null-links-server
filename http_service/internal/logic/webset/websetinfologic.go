package webset

import (
	"context"
	"null-links/http_service/internal/svc"
	"null-links/http_service/internal/types"
	"null-links/internal"

	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

type WebsetInfoLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewWebsetInfoLogic(ctx context.Context, svcCtx *svc.ServiceContext) *WebsetInfoLogic {
	return &WebsetInfoLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *WebsetInfoLogic) WebsetInfo(req *types.WebsetInfoReq) (resp *types.WebsetInfoResp, err error) {
	resp = &types.WebsetInfoResp{}

	WebsetDb, err := l.svcCtx.WebsetModel.FindOne(l.ctx, req.WebsetId)
	if err != nil {
		logx.Error("get webset info from db failed, err: ", err)
		resp.StatusCode = internal.StatusRpcErr
		resp.StatusMsg = "get webset info from db failed"
		return resp, err
	}

	// 获取是否点赞信息
	isLike := false
	likeInfoDb, err := l.svcCtx.LikeModel.GetLikeWebsetUserInfo(l.ctx, req.WebsetId, req.UserId)
	if err != nil && err != sqlx.ErrNotFound {
		logx.Error("get like info error: ", err)
	} else if err == sqlx.ErrNotFound {
		isLike = false
	} else if likeInfoDb.Status == 1 {
		isLike = true
	} else {
		isLike = false
	}

	// 获取是否收藏信息
	isFavoriteResp := false
	// favoriteInfoDb, err = l.svcCtx.FavoriteModel.GetFavoriteWebsetUserInfos(l.ctx, websetIdList, in.UserId)
	// if err != nil {
	// 	l.Logger.Error("get favorite info failed, user, err: ", err)
	// }

	// 获取作者信息
	userInfoDb, err := l.svcCtx.UserModel.FindOne(l.ctx, req.UserId)
	if err != nil {
		logx.Error("get user info from db error. err=", err)
		return nil, err
	}

	// 获取weblink信息
	weblinksDb, err := l.svcCtx.WeblinkModel.FindByWebsetId(l.ctx, req.WebsetId)
	if err != nil {
		logx.Error("get weblink info from db error=", err)
		resp.StatusCode = internal.StatusGatewayErr
		resp.StatusMsg = "fail"
		return resp, nil
	}

	weblinkListResp := make([]types.WebLink, 0, len(weblinksDb))
	for _, weblink := range weblinksDb {
		weblinkResp := types.WebLink{
			Id:       weblink.Id,
			Describe: weblink.Description,
			Url:      weblink.Url,
			CoverURL: weblink.CoverUrl,
		}
		weblinkListResp = append(weblinkListResp, weblinkResp)
	}

	resp.StatusCode = internal.StatusSuccess
	resp.StatusMsg = "success"
	resp.WebsetInfo = types.Webset{
		Id:       WebsetDb.Id,
		Title:    WebsetDb.Title,
		Describe: WebsetDb.Description,
		AuthorInfo: types.User{
			Id:        WebsetDb.AuthorId,
			Name:      userInfoDb.Username,
			AvatarUrl: userInfoDb.AvatarUrl,
		},
		CoverURL:      WebsetDb.CoverUrl,
		ViewCount:     WebsetDb.ViewCnt,
		LikeCount:     WebsetDb.LikeCnt,
		FavoriteCount: WebsetDb.FavoriteCnt,
		IsLike:        isLike,
		IsFavorite:    isFavoriteResp,
		WebLinkList:   weblinkListResp,
	}

	return
}
