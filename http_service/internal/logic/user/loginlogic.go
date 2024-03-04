package user

import (
	"context"
	"time"

	"null-links/http_service/internal/svc"
	"null-links/http_service/internal/types"
	"null-links/internal"
	"null-links/rpc_service/user/pb/user"

	"github.com/zeromicro/go-zero/core/logx"
)

type LoginLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewLoginLogic(ctx context.Context, svcCtx *svc.ServiceContext) *LoginLogic {
	return &LoginLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *LoginLogic) Login(req *types.LoginReq) (resp *types.LoginResp, err error) {
	// check whether both username and email are empty
	if req.Username == "" && req.UserEmail == "" {
		resp = &types.LoginResp{
			StatusCode: internal.StatusParamErr,
			StatusMsg:  "请输入用户名或邮箱",
			UserID:     -1,
		}
		err = nil
		return
	}
	if req.Password == "" {
		resp = &types.LoginResp{
			StatusCode: internal.StatusParamErr,
			StatusMsg:  "请输入密码",
			UserID:     -1,
		}
		err = nil
		return
	}

	resp = &types.LoginResp{}
	respRpc, err := l.svcCtx.UserRpc.Login(l.ctx, &user.LoginReq{
		Username: req.Username,
		Email:    req.UserEmail,
		Password: req.Password,
	})
	if err != nil {
		logx.Error("call UserRpc failed, err: " + err.Error())
		resp.StatusCode = internal.StatusRpcErr
		resp.StatusMsg = "登录失败"
		err = nil
		return
	} else if respRpc.StatusCode != internal.StatusSuccess {
		// the username does not exsit or the password is incorrect
		logx.Error("call UserRpc failed, err: " + respRpc.StatusMsg)
		if respRpc.StatusCode == internal.StatusUserNotExist {
			resp.StatusCode = internal.StatusUserNotExist
			resp.StatusMsg = "该邮箱不存在"
		} else if resp.StatusCode == internal.StatusPasswordErr {
			resp.StatusCode = internal.StatusPasswordErr
			resp.StatusMsg = "密码错误"
		} else {
			resp.StatusCode = internal.StatusRpcErr
			resp.StatusMsg = "登录失败"
		}
		resp.UserID = respRpc.UserId // is -1
		err = nil
		return
	}

	secretKey := l.svcCtx.Config.Auth.AccessSecret
	iat := time.Now().Unix()
	seconds := l.svcCtx.Config.Auth.AccessExpire
	payload := respRpc.UserId // save user id in payload

	token, err := internal.GenJwtToken(secretKey, iat, seconds, payload)
	if err != nil {
		logx.Error("generate token err:", err, " ,token:", token)
		resp.StatusCode = internal.StatusGatewayErr
		resp.StatusMsg = "登录失败"
		resp.UserID = respRpc.UserId // is -1
		err = nil
		return
	}

	resp.StatusCode = internal.StatusSuccess
	resp.StatusMsg = "登录成功"
	resp.UserID = respRpc.UserId
	resp.Token = token
	return
}
