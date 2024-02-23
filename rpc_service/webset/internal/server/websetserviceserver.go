// Code generated by goctl. DO NOT EDIT.
// Source: webset.proto

package server

import (
	"context"

	"null-links/rpc_service/webset/internal/logic"
	"null-links/rpc_service/webset/internal/svc"
	"null-links/rpc_service/webset/pb/webset"
)

type WebsetServiceServer struct {
	svcCtx *svc.ServiceContext
	webset.UnimplementedWebsetServiceServer
}

func NewWebsetServiceServer(svcCtx *svc.ServiceContext) *WebsetServiceServer {
	return &WebsetServiceServer{
		svcCtx: svcCtx,
	}
}

func (s *WebsetServiceServer) Feed(ctx context.Context, in *webset.FeedReq) (*webset.FeedResp, error) {
	l := logic.NewFeedLogic(ctx, s.svcCtx)
	return l.Feed(in)
}

func (s *WebsetServiceServer) PublishAction(ctx context.Context, in *webset.PublishActionReq) (*webset.PublishActionResp, error) {
	l := logic.NewPublishActionLogic(ctx, s.svcCtx)
	return l.PublishAction(in)
}

func (s *WebsetServiceServer) PublishList(ctx context.Context, in *webset.PublishListReq) (*webset.PublishListResp, error) {
	l := logic.NewPublishListLogic(ctx, s.svcCtx)
	return l.PublishList(in)
}

func (s *WebsetServiceServer) LikeAction(ctx context.Context, in *webset.LikeActionReq) (*webset.LikeActionResp, error) {
	l := logic.NewLikeActionLogic(ctx, s.svcCtx)
	return l.LikeAction(in)
}

func (s *WebsetServiceServer) LikeInfoList(ctx context.Context, in *webset.LikeInfoListReq) (*webset.LikeInfoListResp, error) {
	l := logic.NewLikeInfoListLogic(ctx, s.svcCtx)
	return l.LikeInfoList(in)
}

func (s *WebsetServiceServer) FavoriteAction(ctx context.Context, in *webset.FavoriteActionReq) (*webset.FavoriteActionResp, error) {
	l := logic.NewFavoriteActionLogic(ctx, s.svcCtx)
	return l.FavoriteAction(in)
}

func (s *WebsetServiceServer) FavoriteList(ctx context.Context, in *webset.FavoriteListReq) (*webset.FavoriteListResp, error) {
	l := logic.NewFavoriteListLogic(ctx, s.svcCtx)
	return l.FavoriteList(in)
}

func (s *WebsetServiceServer) WebsetInfo(ctx context.Context, in *webset.WebsetInfoReq) (*webset.WebsetInfoResp, error) {
	l := logic.NewWebsetInfoLogic(ctx, s.svcCtx)
	return l.WebsetInfo(in)
}