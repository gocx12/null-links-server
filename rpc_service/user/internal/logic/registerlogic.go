package logic

import (
	"context"
	"encoding/base64"
	"regexp"

	"null-links/internal"
	"null-links/rpc_service/user/internal/model"
	"null-links/rpc_service/user/internal/svc"
	"null-links/rpc_service/user/pb/user"

	"github.com/zeromicro/go-zero/core/logx"
	"golang.org/x/crypto/scrypt"
)

type RegisterLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewRegisterLogic(ctx context.Context, svcCtx *svc.ServiceContext) *RegisterLogic {
	return &RegisterLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *RegisterLogic) Register(in *user.RegisterReq) (*user.RegisterResp, error) {
	resp := &user.RegisterResp{}

	// 检查验证码是否正确。使用邮箱作为key，因此如果不是使用获取验证码时使用的邮箱，也会报错
	validationCode, err := l.svcCtx.RedisClient.Get(l.ctx, RdsKeyEmailValidationPre+"_"+in.Email).Result()
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
	logx.Debug("req valid:", in.ValidationCode, ", redis valid:", validationCode)
	if validationCode != in.ValidationCode {
		resp.StatusCode = internal.StatusValidationCodeErr
		resp.StatusMsg = "the validation code is error"
		return resp, nil
	}

	hash, err := scrypt.Key([]byte(in.Password), SALT, 1<<15, 8, 1, PW_HASH_BYTES)
	if err != nil {
		logx.Error("scrpyt error: ", err)
		return resp, nil
	}
	encodedHash := base64.StdEncoding.EncodeToString(hash)
	data := &model.TUser{
		Username: in.Username,
		Email:    in.Email,
		Password: encodedHash,
	}

	res, err := l.svcCtx.UserModel.Insert(l.ctx, data)
	if err != nil {
		logx.Error("insert user into mysql error: ", err, ", data: ", data)
		if match, _ := regexp.MatchString(".*(23000).*", err.Error()); match {
			resp.StatusCode = internal.StatusUserNameExist
			resp.StatusMsg = "this username has already existed"
			resp.UserId = -1

			return resp, nil
		}
		resp.StatusCode = internal.StatusRpcErr
		resp.StatusMsg = "insert user info into mysql error: " + err.Error()
		return resp, nil
	}

	id, err := res.LastInsertId()
	if err != nil {
		logx.Error("get user id error: ", err.Error())
		resp.StatusCode = internal.StatusRpcErr
		resp.StatusMsg = "get user id error: " + err.Error()
		return resp, err
	}

	resp.StatusCode = internal.StatusSuccess
	resp.StatusMsg = "success"
	resp.UserId = id
	return resp, nil
}
