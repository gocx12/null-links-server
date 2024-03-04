package webset

import (
	"context"

	"null-links/http_service/internal/svc"
	"null-links/http_service/internal/types"
	"null-links/rpc_service/webset/pb/webset"

	"null-links/internal"

	"github.com/zeromicro/go-zero/core/logx"
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

	var publishActionRpcReq webset.PublishActionReq
	if req.ActionType == 1 || req.ActionType == 2 {
		// 发布 或 修改
		weblinkListRpcReq := make([]*webset.WebLink, 0, len(req.WebLinkList))
		for _, weblink := range req.WebLinkList {
			weblinkListRpcReq = append(weblinkListRpcReq, &webset.WebLink{
				Url:      weblink.Url,
				Describe: weblink.Describe,
			})
		}
		publishActionRpcReq = webset.PublishActionReq{
			ActionType: req.ActionType,
			UserId:     req.AuthorId,
			Webset: &webset.Webset{
				Title:       req.Title,
				Describe:    req.Describe,
				CoverUrl:    req.CoverURL,
				WebLinkList: weblinkListRpcReq,
			},
		}
	} else {
		resp.StatusCode = internal.StatusParamErr
		resp.StatusMsg = "未知操作类型"
		return
	}
	publishActionRpcResp, err := l.svcCtx.WebsetRpc.PublishAction(l.ctx, &publishActionRpcReq)

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
