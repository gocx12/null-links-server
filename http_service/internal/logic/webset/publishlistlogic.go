package webset

import (
	"context"

	"null-links/cron/model"
	"null-links/http_service/internal/svc"
	"null-links/http_service/internal/types"
	"null-links/internal"

	"github.com/demdxx/gocast"
	"github.com/zeromicro/go-zero/core/logx"
)

type PublishListLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewPublishListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *PublishListLogic {
	return &PublishListLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *PublishListLogic) PublishList(req *types.PublishListReq) (resp *types.PublishListResp, err error) {
	resp = &types.PublishListResp{}
	userId := gocast.ToInt64(l.ctx.Value("userId"))

	publishListDb, err := l.svcCtx.WebsetModel.FindPublishList(l.ctx, userId, req.Page, req.PageSize)
	if err != nil {
		logx.Error("get publish list from db error=", err, " ,userId=", userId)
		resp.StatusCode = internal.StatusRpcErr
		resp.StatusMsg = "get publish list failed from db error"
		return
	}
	// 获取用户信息
	userInfoDb, err := l.svcCtx.UserModel.FindOne(l.ctx, userId)
	if err != nil && err != model.ErrNotFound {
		logx.Error("get user info from db error. err=", err)
		return nil, nil
	}

	WebsetListResp := make([]types.WebsetShort, 0, len(publishListDb))
	// TODO(chancy): 增加在线状态
	for _, item := range publishListDb {
		WebsetListResp = append(WebsetListResp, types.WebsetShort{
			Id:            item.Id,
			Title:         item.Title,
			CoverUrl:      item.CoverUrl,
			ViewCount:     item.ViewCnt,
			LikeCount:     item.LikeCnt,
			IsLike:        false,
			FavoriteCount: item.FavoriteCnt,
			AuthorInfo: types.UserShort{
				Id:        item.AuthorId,
				Name:      userInfoDb.Username,
				AvatarUrl: userInfoDb.AvatarUrl,
			},
		})
	}

	resp.StatusCode = internal.StatusSuccess
	resp.StatusMsg = "success"
	resp.WebsetList = WebsetListResp

	return
}
