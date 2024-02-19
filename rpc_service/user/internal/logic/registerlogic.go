package logic

import (
	"context"
	"encoding/base64"
	"regexp"

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
	// 检查验证码是否正确
	validationCode, err := l.svcCtx.RedisClient.Get(RdsKeyEmailValidationPre + in.Email)
	if err != nil {
		logx.Error("get validation code failed, err: ", err)
		return nil, err
	}

	if validationCode != in.ValidationCode {
		return &user.RegisterResp{
			StatusMsg: "error validation code",
			UserId:    -1,
		}, nil
	}

	hash, err := scrypt.Key([]byte(in.Password), SALT, 1<<15, 8, 1, PW_HASH_BYTES)
	if err != nil {
		return nil, err
	}
	encodedHash := base64.StdEncoding.EncodeToString(hash)
	data := &model.TUser{
		Username: in.Username,
		Password: encodedHash,
	}
	res, err := l.svcCtx.UserModel.Insert(l.ctx, data)

	switch err {
	case nil:
		id, err := res.LastInsertId()
		if err != nil {
			return nil, err
		}
		return &user.RegisterResp{
			StatusMsg: "success",
			UserId:    id,
		}, nil
	default:
		if match, _ := regexp.MatchString(".*(23000).*", err.Error()); match {
			return &user.RegisterResp{
				StatusMsg: "this username has already existed",
				UserId:    -1,
			}, nil
		}

		return nil, err
	}
}
