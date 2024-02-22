package webset

import (
	"context"

	"null-links/http_service/internal/svc"
	"null-links/http_service/internal/types"
	"null-links/internal"

	"null-links/rpc_service/webset/pb/webset"

	"github.com/demdxx/gocast"
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

	var userId int64 = -1
	if req.Token != "" {
		// 即使没有解析出token，也继续执行，减少对使用体验的影响
		claims, err := internal.ParseJwtToken(l.svcCtx.Config.Auth.AccessSecret, req.Token)
		if err != nil {
			logx.Error("parse jwt token err:", err)
		}
		userId = gocast.ToInt64(claims["user_id"])
	}

	respRpc, err := l.svcCtx.WebsetRpc.Feed(l.ctx, &webset.FeedReq{
		UserId: userId,
	})
	if err != nil {
		resp.StatusCode = internal.StatusRpcErr
		resp.StatusMsg = "获取网页单流失败"
		err = nil
		return
	}

	respRpc.StatusCode = internal.StatusSuccess
	respRpc.StatusMsg = "获取网页单流成功"
	resp.WebsetList = make([]types.WebsetShort, 0, len(respRpc.WebsetList))
	for _, webset := range respRpc.WebsetList {
		resp.WebsetList = append(resp.WebsetList, types.WebsetShort{
			ID:    webset.Id,
			Title: webset.Title,
			AuthorInfo: types.UserShort{
				Id:        webset.AuthorInfo.Id,
				Name:      webset.AuthorInfo.Name,
				AvatarUrl: webset.AuthorInfo.AvatarUrl,
			},
		})
	}

	return
}
