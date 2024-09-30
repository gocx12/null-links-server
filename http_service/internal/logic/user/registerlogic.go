package user

import (
	"context"
	"encoding/base64"
	"regexp"
	"time"

	"null-links/http_service/internal/common"
	"null-links/http_service/internal/infrastructure/model"
	"null-links/http_service/internal/svc"
	"null-links/http_service/internal/types"
	"null-links/internal"

	"github.com/zeromicro/go-zero/core/logx"
	"golang.org/x/crypto/scrypt"
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

	// 检查验证码是否正确。使用邮箱作为key，因此如果不是使用获取验证码时使用的邮箱，也会报错
	validationCode, err := l.svcCtx.RedisClient.Get(l.ctx, common.RdsKeyEmailValidationPre+"_"+req.UserEmail).Result()
	if err != nil {
		if err.Error() == "redis: nil" {
			logx.Debug("validation code do not exist in redis")
			resp.StatusCode = internal.StatusValidationCodeErr
			resp.StatusMsg = "the validation code is error"
			return resp, nil
		}
		logx.Error("get validation code from redis failed, err: ", err)
		resp.StatusCode = internal.StatusRpcErr
		resp.StatusMsg = "get validation code from redis failed"
		return resp, nil
	}
	logx.Debug("req valid:", req.ValidationCode, ", redis valid:", validationCode)
	if validationCode != req.ValidationCode {
		resp.StatusCode = internal.StatusValidationCodeErr
		resp.StatusMsg = "the validation code is error"
		return resp, nil
	}

	hash, err := scrypt.Key([]byte(req.Password), SALT, 1<<15, 8, 1, PW_HASH_BYTES)
	if err != nil {
		logx.Error("scrpyt error: ", err)
		return resp, nil
	}
	encodedHash := base64.StdEncoding.EncodeToString(hash)
	data := &model.TUser{
		Username: req.Username,
		Email:    req.UserEmail,
		Password: encodedHash,
	}

	resDB, err := l.svcCtx.UserModel.Insert(l.ctx, data)
	if err != nil {
		logx.Error("insert user into mysql error: ", err, ", data: ", data)
		if match, _ := regexp.MatchString(".*(23000).*uidx_email.*", err.Error()); match {
			resp.StatusCode = internal.StatusEmailExist
			resp.StatusMsg = "this email has already existed"
			resp.UserID = -1
			return resp, nil
		} else if match, _ := regexp.MatchString(".*(23000).*uidx_username.*", err.Error()); match {
			resp.StatusCode = internal.StatusUserNameExist
			resp.StatusMsg = "this username has already existed"
			resp.UserID = -1
			return resp, nil
		}
		resp.StatusCode = internal.StatusRpcErr
		resp.StatusMsg = "insert user info into mysql error: " + err.Error()
		return resp, nil
	}

	id, err := resDB.LastInsertId()
	logx.Debug("debug: register rpc response: ", resp)
	secretKey := l.svcCtx.Config.Auth.AccessSecret
	iat := time.Now().Unix()
	seconds := l.svcCtx.Config.Auth.AccessExpire
	payload := id
	token, err := internal.GenJwtToken(secretKey, iat, seconds, payload)
	if err != nil {
		logx.Error("get jwt token error:", err)
		resp = &types.RegisterResp{
			StatusCode: internal.StatusGatewayErr,
			StatusMsg:  "register error",
			UserID:     -1, // is -1
		}
		err = nil
		return
	}

	resp.StatusCode = internal.StatusSuccess
	resp.StatusMsg = "success"
	resp.UserID = id
	resp.Token = token
	return
}
