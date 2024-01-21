// Code generated by goctl. DO NOT EDIT.
// Source: video.proto

package server

import (
	"context"

	"nulltv/rpc_service/video/internal/logic"
	"nulltv/rpc_service/video/internal/svc"
	"nulltv/rpc_service/video/pb/video"
)

type VideoServiceServer struct {
	svcCtx *svc.ServiceContext
	video.UnimplementedVideoServiceServer
}

func NewVideoServiceServer(svcCtx *svc.ServiceContext) *VideoServiceServer {
	return &VideoServiceServer{
		svcCtx: svcCtx,
	}
}

func (s *VideoServiceServer) Parse(ctx context.Context, in *video.ParseReq) (*video.ParseResp, error) {
	l := logic.NewParseLogic(ctx, s.svcCtx)
	return l.Parse(in)
}
