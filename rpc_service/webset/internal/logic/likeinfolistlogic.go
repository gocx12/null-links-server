package logic

import (
	"context"

	"null-links/rpc_service/webset/internal/svc"
	"null-links/rpc_service/webset/pb/webset"

	"github.com/zeromicro/go-zero/core/logx"
)

type LikeInfoListLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewLikeInfoListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *LikeInfoListLogic {
	return &LikeInfoListLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *LikeInfoListLogic) LikeInfoList(in *webset.LikeInfoListReq) (*webset.LikeInfoListResp, error) {
	l.svcCtx.LikeModel.GetLikeWebsetUserInfos(l.ctx, in.WebsetIdList, in.UserId)
	return &webset.LikeInfoListResp{}, nil
}
