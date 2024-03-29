// Code generated by goctl. DO NOT EDIT.
// Source: user.proto

package server

import (
	"context"

	"null-links/rpc_service/user/internal/logic"
	"null-links/rpc_service/user/internal/svc"
	"null-links/rpc_service/user/pb/user"
)

type UserServiceServer struct {
	svcCtx *svc.ServiceContext
	user.UnimplementedUserServiceServer
}

func NewUserServiceServer(svcCtx *svc.ServiceContext) *UserServiceServer {
	return &UserServiceServer{
		svcCtx: svcCtx,
	}
}

func (s *UserServiceServer) Register(ctx context.Context, in *user.RegisterReq) (*user.RegisterResp, error) {
	l := logic.NewRegisterLogic(ctx, s.svcCtx)
	return l.Register(in)
}

func (s *UserServiceServer) Login(ctx context.Context, in *user.LoginReq) (*user.LoginResp, error) {
	l := logic.NewLoginLogic(ctx, s.svcCtx)
	return l.Login(in)
}

func (s *UserServiceServer) CheckUsername(ctx context.Context, in *user.CheckUsernameReq) (*user.CheckUsernameResp, error) {
	l := logic.NewCheckUsernameLogic(ctx, s.svcCtx)
	return l.CheckUsername(in)
}

func (s *UserServiceServer) GetValidtaionCode(ctx context.Context, in *user.GetValidtaionCodeReq) (*user.GetValidtaionCodeResp, error) {
	l := logic.NewGetValidtaionCodeLogic(ctx, s.svcCtx)
	return l.GetValidtaionCode(in)
}

func (s *UserServiceServer) UserInfo(ctx context.Context, in *user.UserInfoReq) (*user.UserInfoResp, error) {
	l := logic.NewUserInfoLogic(ctx, s.svcCtx)
	return l.UserInfo(in)
}

func (s *UserServiceServer) UserInfoList(ctx context.Context, in *user.UserInfoListReq) (*user.UserInfoListResp, error) {
	l := logic.NewUserInfoListLogic(ctx, s.svcCtx)
	return l.UserInfoList(in)
}

func (s *UserServiceServer) Modify(ctx context.Context, in *user.ModifyReq) (*user.ModifyResp, error) {
	l := logic.NewModifyLogic(ctx, s.svcCtx)
	return l.Modify(in)
}
