package webset

import (
	"context"
	"null-links/http_service/internal/svc"
	"null-links/http_service/internal/types"
	"null-links/internal"
	"null-links/rpc_service/webset/pb/webset"

	"github.com/zeromicro/go-zero/core/logx"
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

	websetInfoRpcReq, err := l.svcCtx.WebsetRpc.WebsetInfo(l.ctx, &webset.WebsetInfoReq{
		UserId:   req.UserID,
		WebsetId: req.WebsetID,
	})

	if err != nil {
		logx.Error("call WebsetRpc failed, err: ", err)
		resp.StatusCode = internal.StatusRpcErr
		resp.StatusMsg = "获取网页单失败"
		err = nil
		return
	} else if websetInfoRpcReq.StatusCode != internal.StatusSuccess {
		logx.Error("call WebsetRpc failed, err: ", websetInfoRpcReq.StatusMsg)
		resp.StatusCode = internal.StatusRpcErr
		resp.StatusMsg = "获取网页单失败"
		return
	}

	resp.StatusCode = internal.StatusSuccess
	resp.StatusMsg = "获取网页单成功"
	weblinkListResp := make([]types.WebLink, 0, len(websetInfoRpcReq.Webset.WebLinkList))
	for _, weblink := range websetInfoRpcReq.Webset.WebLinkList {
		weblinkResp := types.WebLink{
			ID:       weblink.Id,
			Describe: weblink.Describe,
			Url:      weblink.Url,
			CoverURL: weblink.CoverUrl,
		}
		weblinkListResp = append(weblinkListResp, weblinkResp)
	}
	resp.WebsetInfo = types.Webset{
		ID:            websetInfoRpcReq.Webset.Id,
		Title:         websetInfoRpcReq.Webset.Title,
		Describe:      websetInfoRpcReq.Webset.Describe,
		ViewCount:     websetInfoRpcReq.Webset.ViewCount,
		LikeCount:     websetInfoRpcReq.Webset.LikeCount,
		FavoriteCount: websetInfoRpcReq.Webset.FavoriteCount,
		IsLike:        websetInfoRpcReq.Webset.IsLike,
		IsFavorite:    websetInfoRpcReq.Webset.IsFavorite,
		WebLinkList:   weblinkListResp,
		AuthorInfo: types.User{
			Id:        websetInfoRpcReq.Webset.AuthorInfo.Id,
			Name:      websetInfoRpcReq.Webset.AuthorInfo.Name,
			AvatarUrl: websetInfoRpcReq.Webset.AuthorInfo.AvatarUrl,
		},
	}

	return
}
