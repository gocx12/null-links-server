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

type RegisterLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewRegisterLogic(ctx context.Context, svcCtx *svc.ServiceContext) *RegisterLogic {
	return &RegisterLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *RegisterLogic) Register(req *types.RegisterReq) (resp *types.RegisterResp, err error) {
	resp = &types.RegisterResp{}
	respRpc, err := l.svcCtx.UserRpc.Register(l.ctx, &user.RegisterReq{
		Username:       req.Username,
		Email:          req.UserEmail,
		ValidationCode: req.ValidationCode,
		Password:       req.Password,
	})
	if err != nil {
		logx.Error("call register rpc error: ", err)
		resp.StatusCode = internal.StatusRpcErr
		resp.StatusMsg = "注册失败"
		err = nil
		return
	} else if respRpc.StatusCode != internal.StatusSuccess {
		if respRpc.StatusCode == internal.StatusEmailExist {
			resp.StatusCode = internal.StatusGatewayErr
			resp.StatusMsg = "该邮箱已注册，请更换邮箱，或直接登录"
		} else if respRpc.StatusCode == internal.StatusValidationCodeErr {
			resp.StatusCode = internal.StatusGatewayErr
			resp.StatusMsg = "验证码错误"
		} else {
			resp.StatusCode = internal.StatusRpcErr
			resp.StatusMsg = "注册失败"
		}
		err = nil
		return
	}
	logx.Debug("debug: register rpc response: ", respRpc)
	secretKey := l.svcCtx.Config.Auth.AccessSecret
	iat := time.Now().Unix()
	seconds := l.svcCtx.Config.Auth.AccessExpire
	payload := respRpc.UserId
	token, err := internal.GenJwtToken(secretKey, iat, seconds, payload)
	if err != nil {
		logx.Error("get jwt token error:", err)
		resp = &types.RegisterResp{
			StatusCode: internal.StatusGatewayErr,
			StatusMsg:  "注册失败",
			UserID:     respRpc.UserId, // is -1
		}
		err = nil
		return
	}

	resp.StatusCode = internal.StatusSuccess
	resp.StatusMsg = "注册成功"
	resp.UserID = respRpc.UserId
	resp.Token = token

	return
}
