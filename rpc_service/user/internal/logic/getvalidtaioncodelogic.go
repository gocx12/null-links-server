package logic

import (
	"context"

	"null-links/rpc_service/user/internal/svc"
	"null-links/rpc_service/user/pb/user"

	"github.com/zeromicro/go-zero/core/logx"

	"math/rand"
	"time"
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
	logx.Debug("GetValidtaionCodeLogic.GetValidtaionCode", "email: ", in.Email)

	recipient := in.Email
	validationCode := l.generateValidationCode()

	// Redis存储验证码, 5分钟过期
	err := l.svcCtx.RedisClient.Setex(RdsKeyEmailValidationPre+recipient, validationCode, 300)
	if err != nil {
		logx.Error("failed to set email validation code to redis:", err, " email: ", recipient)
		return &user.GetValidtaionCodeResp{
			StatusCode: 0,
			StatusMsg:  "set email validation code to redis failed, err: " + err.Error(),
		}, err
	}

	data := validationCode + "::" + recipient
	if err := l.svcCtx.ValidationKqPusherClient.Push(data); err != nil {
		logx.Errorf("ValidationKqPusherClient Push Error , err :%v", err)
		return &user.GetValidtaionCodeResp{
			StatusCode: 0,
			StatusMsg:  "push email validation code to kq failed, err: " + err.Error(),
		}, err
	}

	return &user.GetValidtaionCodeResp{
		StatusCode: 1,
		StatusMsg:  "success",
	}, nil
}

func (l *GetValidtaionCodeLogic) generateValidationCode() string {
	rand.Seed(time.Now().UnixNano())
	chars := []rune("ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz")
	length := 4
	buf := make([]rune, length)
	for i := range buf {
		buf[i] = chars[rand.Intn(len(chars))]
	}
	return string(buf)
}
