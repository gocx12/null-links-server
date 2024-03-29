// Code generated by goctl. DO NOT EDIT.
// Source: search.proto

package server

import (
	"context"

	"null-links/rpc_service/search/internal/logic"
	"null-links/rpc_service/search/internal/svc"
	"null-links/rpc_service/search/pb/search"
)

type SearchServiceServer struct {
	svcCtx *svc.ServiceContext
	search.UnimplementedSearchServiceServer
}

func NewSearchServiceServer(svcCtx *svc.ServiceContext) *SearchServiceServer {
	return &SearchServiceServer{
		svcCtx: svcCtx,
	}
}

func (s *SearchServiceServer) AddToSearch(ctx context.Context, in *search.AddToSearchReq) (*search.AddToSearchResp, error) {
	l := logic.NewAddToSearchLogic(ctx, s.svcCtx)
	return l.AddToSearch(in)
}

func (s *SearchServiceServer) Search(ctx context.Context, in *search.SearchReq) (*search.SearchResp, error) {
	l := logic.NewSearchLogic(ctx, s.svcCtx)
	return l.Search(in)
}
