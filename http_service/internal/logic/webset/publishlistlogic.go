package webset

import (
	"context"

	"null-links/http_service/internal/svc"
	"null-links/http_service/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
	"null-links/internal"
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
	publishListRpcResp, err := l.svcCtx.WebsetRpc.PublishList(l.ctx, &types.PublishListReq{
		UserID:   req.UserID,
		Page:     req.Page,
		PageSize: req.PageSize,
	})
	if err != nil {
		logx.Error("call WebsetRpc failed, err: ", err)
		resp.StatusCode = internal.StatusRpcErr
		resp.StatusMsg = "获取发布列表失败"
		return
	}

	resp.StatusCode = internal.StatusSuccess
	resp.StatusMsg = "获取发布列表成功"
	resp.WebsetList = make([]types.Webset, 0, len(publishListRpcResp.WebsetList))
	for _, webset := range publishListRpcResp.WebsetList {
		resp.WebsetList = append(resp.WebsetList, types.Webset{
			ID:       webset.Id,
			Title:    webset.Title,
			CoverURL: webset.CoverUrl,
			AuthorInfo: types.User{
				Id:        webset.AuthorInfo.Id,
				Name:      webset.AuthorInfo.Name,
				AvatarUrl: webset.AuthorInfo.AvatarUrl,
			},
			CreatedAt: webset.CreatedAt,
		})
	}

	return
}
