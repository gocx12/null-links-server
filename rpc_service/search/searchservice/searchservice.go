// Code generated by goctl. DO NOT EDIT.
// Source: search.proto

package searchservice

import (
	"context"

	"null-links/rpc_service/search/pb/search"

	"github.com/zeromicro/go-zero/zrpc"
	"google.golang.org/grpc"
)

type (
	AddToSearchReq  = search.AddToSearchReq
	AddToSearchResp = search.AddToSearchResp
	SearchReq       = search.SearchReq
	SearchResp      = search.SearchResp
	UserInfoSearch  = search.UserInfoSearch
	UserInfoShort   = search.UserInfoShort
	WebsetSearch    = search.WebsetSearch
	WebsetShort     = search.WebsetShort

	SearchService interface {
		AddToSearch(ctx context.Context, in *AddToSearchReq, opts ...grpc.CallOption) (*AddToSearchResp, error)
		Search(ctx context.Context, in *SearchReq, opts ...grpc.CallOption) (*SearchResp, error)
	}

	defaultSearchService struct {
		cli zrpc.Client
	}
)

func NewSearchService(cli zrpc.Client) SearchService {
	return &defaultSearchService{
		cli: cli,
	}
}

func (m *defaultSearchService) AddToSearch(ctx context.Context, in *AddToSearchReq, opts ...grpc.CallOption) (*AddToSearchResp, error) {
	client := search.NewSearchServiceClient(m.cli.Conn())
	return client.AddToSearch(ctx, in, opts...)
}

func (m *defaultSearchService) Search(ctx context.Context, in *SearchReq, opts ...grpc.CallOption) (*SearchResp, error) {
	client := search.NewSearchServiceClient(m.cli.Conn())
	return client.Search(ctx, in, opts...)
}
