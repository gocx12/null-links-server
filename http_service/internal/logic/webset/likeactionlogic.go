package webset

import (
	"context"

	"null-links/http_service/internal/svc"
	"null-links/http_service/internal/types"
	"null-links/internal"
	"null-links/rpc_service/webset/pb/webset"

	"github.com/zeromicro/go-zero/core/logx"
)

type LikeActionLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewLikeActionLogic(ctx context.Context, svcCtx *svc.ServiceContext) *LikeActionLogic {
	return &LikeActionLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *LikeActionLogic) LikeAction(req *types.LikeActionReq) (resp *types.LikeActionResp, err error) {
	resp = &types.LikeActionResp{}

	LikeActionRpcReq := &webset.LikeActionReq{
		UserId:     req.UserId,
		ActionType: req.ActionType,
		WebsetId:   req.WebsetId,
	}
	likeActionRpcResp, err := l.svcCtx.WebsetRpc.LikeAction(l.ctx, LikeActionRpcReq)

	if likeActionRpcResp.StatusCode != internal.StatusSuccess {
		resp.StatusCode = internal.StatusRpcErr
		if req.ActionType == 1 {
			resp.StatusMsg = "点赞失败"

		} else if req.ActionType == 2 {
			resp.StatusMsg = "取消赞失败"
		}
	}

	resp.StatusCode = internal.StatusSuccess
	if req.ActionType == 1 {
		resp.StatusMsg = "点赞成功"

	} else if req.ActionType == 2 {
		resp.StatusMsg = "取消赞成功"
	}

	return
}
