package user

import (
	"context"
	"time"

	"null-links/http_service/internal/svc"
	"null-links/http_service/internal/types"
	"null-links/internal"
	"null-links/rpc_service/user/pb/user"

	"github.com/zeromicro/go-zero/core/logc"
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

	respRpc, err := l.svcCtx.UserRpc.Login(l.ctx, &user.LoginReq{
		Username: req.Username,
		Email:    req.UserEmail,
		Password: req.Password,
	})
	if err != nil {
		resp = &types.LoginResp{
			StatusCode: internal.StatusRpcErr,
			StatusMsg:  "登录失败",
			UserID:     respRpc.UserId, // is -1
		}
		logc.Error(l.ctx, "call UserRpc failed, err: "+err.Error())
		err = nil
		return
	} else if respRpc.UserId == -1 {
		// the username does not exsit or the password is incorrect
		resp = &types.LoginResp{
			StatusCode: respRpc.StatusCode,
			StatusMsg:  respRpc.StatusMsg,
			UserID:     respRpc.UserId, // is -1
		}
		logc.Error(l.ctx, "call UserRpc failed, err: "+respRpc.StatusMsg)
		err = nil
		return
	}

	secretKey := l.svcCtx.Config.Auth.AccessSecret
	iat := time.Now().Unix()
	seconds := l.svcCtx.Config.Auth.AccessExpire
	payload := respRpc.UserId // save user id in payload

	token, err := internal.GetJwtToken(secretKey, iat, seconds, payload)
	if err != nil {
		resp = &types.LoginResp{
			StatusCode: internal.StatusGatewayErr,
			StatusMsg:  "登录失败",
			UserID:     respRpc.UserId, // is -1
		}
		logc.Error(l.ctx, "getJwtToken() "+err.Error())
		err = nil
		return
	}

	resp = &types.LoginResp{
		StatusCode: internal.StatusSuccess,
		StatusMsg:  "登录成功",
		UserID:     respRpc.UserId,
		Token:      token,
	}
	return
}
