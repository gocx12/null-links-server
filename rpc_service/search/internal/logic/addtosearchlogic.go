package logic

import (
	"context"

	"null-links/rpc_service/search/internal/svc"
	"null-links/rpc_service/search/pb/search"

	"github.com/zeromicro/go-zero/core/logx"
)

type AddToSearchLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewAddToSearchLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AddToSearchLogic {
	return &AddToSearchLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *AddToSearchLogic) AddToSearch(in *search.AddToSearchReq) (*search.AddToSearchResp, error) {

	if in.DataType == 1 {

	} else if in.DataType == 2 {

	}
	return &search.AddToSearchResp{}, nil
}
