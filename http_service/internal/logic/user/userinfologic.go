package user

import (
	"context"
	"log"
	"time"

	"nulltv/http_service/internal/svc"
	"nulltv/http_service/internal/types"
	"nulltv/internal"
	"nulltv/rpc_service/user/pb/user"

	"github.com/zeromicro/go-zero/core/logx"
)

type UserInfoLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewUserInfoLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UserInfoLogic {
	return &UserInfoLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UserInfoLogic) UserInfo(req *types.UserInfoReq) (resp *types.UserInfoResp, err error) {
	respRpc, err := l.svcCtx.UserRpc.UserInfo(l.ctx, &user.UserInfoReq{
		UserId: req.UserID,
	})
	if err != nil {
		resp = &types.UserInfoResp{
			StatusCode: internal.StatusRpcErr,
			StatusMsg:  "获取信息失败",
			UserID:     respRpc.UserId, // is -1
		}
		log.Fatal(err)
		err = nil
		return
	} else if respRpc.UserId == -1 {
		// the username does not exsit or the password is incorrect
		resp = &types.UserInfoResp{
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
		resp = &types.UserInfoResp{
			StatusCode: internal.StatusGatewayErr,
			StatusMsg:  "获取信息失败",
			UserID:     respRpc.UserId, // is -1
		}
		log.Fatal(err)
		err = nil
		return
	}

	resp = &types.UserInfoResp{
		StatusCode: internal.StatusSuccess,
		StatusMsg:  respRpc.StatusMsg,
		UserID:     respRpc.UserId,
		Token:      token,
	}

	return
}
