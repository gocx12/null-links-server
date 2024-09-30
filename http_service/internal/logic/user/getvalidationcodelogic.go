package user

import (
	"context"
	"math/rand"
	"time"

	"null-links/http_service/internal/common"
	"null-links/http_service/internal/svc"
	"null-links/http_service/internal/types"
	"null-links/internal"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetValidationCodeLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetValidationCodeLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetValidationCodeLogic {
	return &GetValidationCodeLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetValidationCodeLogic) GetValidationCode(req *types.GetValidationCodeReq) (resp *types.GetValidationCodeResp, err error) {
	// respRpc, err := l.svcCtx.UserRpc.GetValidtaionCode(l.ctx, &user.GetValidtaionCodeReq{
	// 	Email: req.Email,
	// })
	// if err != nil {
	// 	logx.Error("call UserRpc failed, err: ", err)
	// 	resp = &types.GetValidationCodeResp{
	// 		StatusCode: internal.StatusRpcErr,
	// 		StatusMsg:  "获取验证码失败",
	// 	}
	// 	err = nil
	// 	return
	// } else if respRpc.StatusCode != internal.StatusSuccess {
	// 	logx.Error("call UserRpc failed, err: ", resp.StatusMsg)
	// 	resp = &types.GetValidationCodeResp{
	// 		StatusCode: internal.StatusRpcErr,
	// 		StatusMsg:  "获取验证码失败",
	// 	}
	// 	return
	// }

	logx.Debug("get validation cod, email: ", req.Email)

	recipient := req.Email
	validationCode := l.generateValidationCode()

	// Redis存储验证码, 10分钟过期
	_, err = l.svcCtx.RedisClient.SetEx(l.ctx, common.RdsKeyEmailValidationPre+"_"+recipient, validationCode, 10*time.Minute).Result()
	if err != nil {
		logx.Error("failed to set email validation code to redis, error:", err, " email: ", recipient)
		resp = &types.GetValidationCodeResp{
			StatusCode: internal.StatusRpcErr,
			StatusMsg:  "set email validation code to redis failed",
		}
		return
	}

	// kafka pusher
	data := recipient + "::" + validationCode
	if err := l.svcCtx.VdEmailMqPusher.Push(data); err != nil {
		logx.Error("VdEmailMqPusher Push error:", err)
		resp = &types.GetValidationCodeResp{
			StatusCode: internal.StatusGatewayErr,
			StatusMsg:  "push email validation code to kq failed, err: " + err.Error(),
		}
		return resp, nil
	}

	resp = &types.GetValidationCodeResp{
		StatusCode: internal.StatusSuccess,
		StatusMsg:  "成功",
	}
	return
}

func (l *GetValidationCodeLogic) generateValidationCode() string {
	// go 1.20 起弃用rand.Seed, 即使不用seed也可以产生不同的数字
	chars := []rune("ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz")
	length := 4
	buf := make([]rune, length)
	for i := range buf {
		buf[i] = chars[rand.Intn(len(chars))]
	}
	return string(buf)
}
