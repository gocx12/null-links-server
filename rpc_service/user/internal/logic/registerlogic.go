package logic

import (
	"context"
	"encoding/base64"
	"regexp"

	"nulltv/rpc_service/user/internal/model"
	"nulltv/rpc_service/user/internal/svc"
	"nulltv/rpc_service/user/pb/user"

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
	hash, err := scrypt.Key([]byte(in.Password), SALT, 1<<15, 8, 1, PW_HASH_BYTES)
	if err != nil {
		return nil, err
	}
	encodedHash := base64.StdEncoding.EncodeToString(hash)
	data := &model.User{
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
			StatusMsg: "注册成功",
			UserId:    id,
		}, nil
	default:
		if match, _ := regexp.MatchString(".*(23000).*", err.Error()); match {
			return &user.RegisterResp{
				StatusMsg: "用户名已经存在",
				UserId:    -1,
			}, nil
		}

		return nil, err
	}
}
