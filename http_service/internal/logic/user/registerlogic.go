package user

import (
	"context"
	"log"
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
	respRpc, err := l.svcCtx.UserRpc.Register(l.ctx, &user.RegisterReq{
		Username: req.Username,
		Password: req.Password,
	})
	if err != nil {
		resp = &types.RegisterResp{
			StatusCode: internal.StatusRpcErr,
			StatusMsg:  "注册失败",
			UserID:     respRpc.UserId, // is -1
		}
		log.Fatal(err)
		err = nil
		return
	} else if respRpc.UserId == -1 {
		// the username does not exsit or the password is incorrect
		resp = &types.RegisterResp{
			StatusCode: internal.StatusRpcErr,
			StatusMsg:  respRpc.StatusMsg,
			UserID:     respRpc.UserId, // is -1
		}
		err = nil
		return
	}

	secretKey := l.svcCtx.Config.Auth.AccessSecret
	iat := time.Now().Unix()
	seconds := l.svcCtx.Config.Auth.AccessExpire
	payload := respRpc.UserId
	token, err := getJwtToken(secretKey, iat, seconds, payload)
	if err != nil {
		resp = &types.RegisterResp{
			StatusCode: internal.StatusGatewayErr,
			StatusMsg:  "注册失败",
			UserID:     respRpc.UserId, // is -1
		}
		log.Fatal(err)
		err = nil
		return
	}

	resp = &types.RegisterResp{
		StatusCode: internal.StatusSuccess,
		StatusMsg:  respRpc.StatusMsg,
		UserID:     respRpc.UserId,
		Token:      token,
	}

	return
}
