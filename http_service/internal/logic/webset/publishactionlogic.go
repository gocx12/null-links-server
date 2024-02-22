package webset

import (
	"context"

	"null-links/http_service/internal/svc"
	"null-links/http_service/internal/types"
	"null-links/rpc_service/webset/pb/webset"

	"github.com/zeromicro/go-zero/core/logx"
	"null-links/internal"
)

type PublishActionLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewPublishActionLogic(ctx context.Context, svcCtx *svc.ServiceContext) *PublishActionLogic {
	return &PublishActionLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *PublishActionLogic) PublishAction(req *types.PublishActionReq) (resp *types.PublishActionResp, err error) {
	resp = &types.PublishActionResp{}

	publishActionRpcResp, err := l.svcCtx.WebsetRpc.PublishAction(l.ctx, &webset.PublishActionReq{
		ActionType: req.ActionType,
		UserId:     req.AuthorId,
	})

	if err != nil || publishActionRpcResp.StatusCode != internal.StatusSuccess {
		if err != nil {
			logx.Error("call WebsetRpc failed, err: ", err)
			err = nil
		} else if publishActionRpcResp.StatusCode != internal.StatusSuccess {
			logx.Error("call WebsetRpc failed, err: ", publishActionRpcResp.StatusMsg)
		}
		resp.StatusCode = internal.StatusRpcErr
		if req.ActionType == 1 {
			resp.StatusMsg = "发布失败"
		} else if req.ActionType == 2 {
			resp.StatusMsg = "修改失败"
		} else if req.ActionType == 3 {
			resp.StatusMsg = "删除失败"
		} else {
			resp.StatusMsg = "操作失败"
		}
		return
	}

	resp.StatusCode = internal.StatusSuccess
	if req.ActionType == 1 {
		resp.StatusMsg = "发布成功"
	} else if req.ActionType == 2 {
		resp.StatusMsg = "修改成功"
	} else if req.ActionType == 3 {
		resp.StatusMsg = "删除成功"
	} else {
		resp.StatusMsg = "操作成功"
	}

	return
}
