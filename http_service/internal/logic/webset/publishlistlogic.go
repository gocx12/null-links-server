package webset

import (
	"context"

	"null-links/http_service/internal/infrastructure/model"
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

	publishListDb := make([]*model.TWebset, 0)
	if req.Tag == 1 {
		publishListDb, err = l.svcCtx.WebsetModel.FindPublishList(l.ctx, userId, req.Page, req.PageSize)
	} else {
		status := internal.WebsetPendReview
		switch req.Tag {
		case 2:
			status = internal.WebsetPublished
		case 3:
			status = internal.WebsetRejected
		default:
			status = internal.WebsetPendReview
		}
		publishListDb, err = l.svcCtx.WebsetModel.FindPublishListWithStatus(l.ctx, userId, req.Page, req.PageSize, status.Code())
	}

	if err != nil {
		logx.Error("get publish list from db error=", err, ", userId=", userId)
		resp.StatusCode = internal.StatusRpcErr
		resp.StatusMsg = "get publish list failed from db error"
		return
	}

	WebsetListResp := make([]types.PublishWebset, 0, len(publishListDb))
	// TODO(chancy): 增加在线状态
	for _, item := range publishListDb {
		WebsetListResp = append(WebsetListResp, types.PublishWebset{
			Id:            item.Id,
			Title:         item.Title,
			CoverUrl:      item.CoverUrl,
			ViewCount:     item.ViewCnt,
			LikeCount:     item.LikeCnt,
			FavoriteCount: item.FavoriteCnt,
			Status:        gocast.ToInt32(item.Status),
		})
	}

	resp.StatusCode = internal.StatusSuccess
	resp.StatusMsg = "success"
	resp.WebsetList = WebsetListResp

	return
}
