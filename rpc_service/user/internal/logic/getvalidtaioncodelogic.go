package logic

import (
	"context"
	"time"

	"null-links/rpc_service/user/internal/svc"
	"null-links/rpc_service/user/pb/user"

	"null-links/internal"

	"github.com/zeromicro/go-zero/core/logx"

	"math/rand"
)

type GetValidtaionCodeLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetValidtaionCodeLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetValidtaionCodeLogic {
	return &GetValidtaionCodeLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

var (
	RdsKeyEmailValidationPre = "EMAIL_VALIDATION"
)

func (l *GetValidtaionCodeLogic) GetValidtaionCode(in *user.GetValidtaionCodeReq) (*user.GetValidtaionCodeResp, error) {
	resp := &user.GetValidtaionCodeResp{}
	logx.Debug("get validation cod, email: ", in.Email)

	recipient := in.Email
	validationCode := l.generateValidationCode()

	// Redis存储验证码, 10分钟过期
	_, err := l.svcCtx.RedisClient.SetEx(l.ctx, RdsKeyEmailValidationPre+"_"+recipient, validationCode, 10*time.Minute).Result()
	if err != nil {
		logx.Error("failed to set email validation code to redis, error:", err, " email: ", recipient)
		resp.StatusCode = internal.StatusRpcErr
		resp.StatusMsg = "set email validation code to redis failed"
		return resp, nil
	}

	// kafka pusher
	data := recipient + "::" + validationCode
	if err := l.svcCtx.VdEmailMqPusher.Push(data); err != nil {
		logx.Error("VdEmailMqPusher Push error:", err)
		resp.StatusCode = 0
		resp.StatusMsg = "push email validation code to kq failed, err: " + err.Error()
		return resp, nil
	}

	resp.StatusCode = internal.StatusSuccess
	resp.StatusMsg = "success"
	return resp, nil
}

func (l *GetValidtaionCodeLogic) generateValidationCode() string {
	// go 1.20 起弃用rand.Seed, 即使不用seed也可以产生不同的数字
	chars := []rune("ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz")
	length := 4
	buf := make([]rune, length)
	for i := range buf {
		buf[i] = chars[rand.Intn(len(chars))]
	}
	return string(buf)
}
