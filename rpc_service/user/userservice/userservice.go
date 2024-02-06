// Code generated by goctl. DO NOT EDIT.
// Source: user.proto

package userservice

import (
	"context"

	"null-links/rpc_service/user/pb/user"

	"github.com/zeromicro/go-zero/zrpc"
	"google.golang.org/grpc"
)

type (
	LoginReq         = user.LoginReq
	LoginResp        = user.LoginResp
	RegisterReq      = user.RegisterReq
	RegisterResp     = user.RegisterResp
	UserInfo         = user.UserInfo
	UserInfoListReq  = user.UserInfoListReq
	UserInfoListResp = user.UserInfoListResp
	UserInfoReq      = user.UserInfoReq
	UserInfoResp     = user.UserInfoResp

	UserService interface {
		Register(ctx context.Context, in *RegisterReq, opts ...grpc.CallOption) (*RegisterResp, error)
		Login(ctx context.Context, in *LoginReq, opts ...grpc.CallOption) (*LoginResp, error)
		UserInfo(ctx context.Context, in *UserInfoReq, opts ...grpc.CallOption) (*UserInfoResp, error)
		UserInfoList(ctx context.Context, in *UserInfoListReq, opts ...grpc.CallOption) (*UserInfoListResp, error)
	}

	defaultUserService struct {
		cli zrpc.Client
	}
)

func NewUserService(cli zrpc.Client) UserService {
	return &defaultUserService{
		cli: cli,
	}
}

func (m *defaultUserService) Register(ctx context.Context, in *RegisterReq, opts ...grpc.CallOption) (*RegisterResp, error) {
	client := user.NewUserServiceClient(m.cli.Conn())
	return client.Register(ctx, in, opts...)
}

func (m *defaultUserService) Login(ctx context.Context, in *LoginReq, opts ...grpc.CallOption) (*LoginResp, error) {
	client := user.NewUserServiceClient(m.cli.Conn())
	return client.Login(ctx, in, opts...)
}

func (m *defaultUserService) UserInfo(ctx context.Context, in *UserInfoReq, opts ...grpc.CallOption) (*UserInfoResp, error) {
	client := user.NewUserServiceClient(m.cli.Conn())
	return client.UserInfo(ctx, in, opts...)
}

func (m *defaultUserService) UserInfoList(ctx context.Context, in *UserInfoListReq, opts ...grpc.CallOption) (*UserInfoListResp, error) {
	client := user.NewUserServiceClient(m.cli.Conn())
	return client.UserInfoList(ctx, in, opts...)
}
