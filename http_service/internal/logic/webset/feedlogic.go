package webset

import (
	"context"

	"null-links/http_service/internal/svc"
	"null-links/http_service/internal/types"
	"null-links/internal"

	"null-links/rpc_service/webset/pb/webset"

	"github.com/zeromicro/go-zero/core/logx"
)

type FeedLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewFeedLogic(ctx context.Context, svcCtx *svc.ServiceContext) *FeedLogic {
	return &FeedLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *FeedLogic) Feed(req *types.FeedReq) (resp *types.FeedResp, err error) {
	resp = &types.FeedResp{}

	respRpc, err := l.svcCtx.WebsetRpc.Feed(l.ctx, &webset.FeedReq{
		UserId:   req.UserId,
		Page:     req.Page,
		PageSize: req.PageSize,
	})
	if err != nil {
		resp.StatusCode = internal.StatusRpcErr
		resp.StatusMsg = "获取网页单流失败"
		err = nil
		return
	}

	resp.StatusCode = internal.StatusSuccess
	resp.StatusMsg = "获取网页单流成功"
	resp.WebsetList = make([]types.WebsetShort, 0, len(respRpc.WebsetList))
	for _, webset := range respRpc.WebsetList {
		resp.WebsetList = append(resp.WebsetList, types.WebsetShort{
			ID:        webset.Id,
			Title:     webset.Title,
			CoverUrl:  webset.CoverUrl,
			LikeCount: webset.LikeCount,
			ViewCount: webset.ViewCount,
			AuthorInfo: types.UserShort{
				Id:        webset.AuthorInfo.Id,
				Name:      webset.AuthorInfo.Name,
				AvatarUrl: webset.AuthorInfo.AvatarUrl,
			},
		})
	}

	return
}
