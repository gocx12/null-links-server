package logic

import (
	"context"

	"null-links/rpc_service/search/internal/svc"
	"null-links/rpc_service/search/pb/search"

	"github.com/zeromicro/go-zero/core/logx"
)

type SearchLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewSearchLogic(ctx context.Context, svcCtx *svc.ServiceContext) *SearchLogic {
	return &SearchLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *SearchLogic) Search(in *search.SearchReq) (*search.SearchResp, error) {
	// todo: add your logic here and delete this line

	return &search.SearchResp{}, nil
}
