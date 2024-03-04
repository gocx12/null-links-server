package logic

import (
	"context"
	"encoding/base64"

	"null-links/internal"
	"null-links/rpc_service/user/internal/model"
	"null-links/rpc_service/user/internal/svc"
	"null-links/rpc_service/user/pb/user"

	"github.com/zeromicro/go-zero/core/logx"
	"golang.org/x/crypto/scrypt"
)

type LoginLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

var PW_HASH_BYTES = 32
var SALT = []byte{126, 145, 58, 233, 153, 107, 4, 231}

func NewLoginLogic(ctx context.Context, svcCtx *svc.ServiceContext) *LoginLogic {
	return &LoginLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *LoginLogic) Login(in *user.LoginReq) (*user.LoginResp, error) {
	resp := &user.LoginResp{}

	// UserInfoDb, err := l.svcCtx.UserModel.FindOneByName(l.ctx, in.Username)
	UserInfoDb, err := l.svcCtx.UserModel.FindPasswordByEmail(l.ctx, in.Email)

	switch err {
	case nil:
		hash, err := scrypt.Key([]byte(in.Password), SALT, 1<<15, 8, 1, PW_HASH_BYTES)
		encodedHash := base64.StdEncoding.EncodeToString(hash)
		if err != nil {
			logx.Error()
			resp.StatusCode = internal.StatusRpcErr
			resp.StatusMsg = "scrypt encode password error: " + err.Error()
			resp.UserId = -1
			return nil, err
		}
		if UserInfoDb.Password != encodedHash {
			resp.StatusCode = internal.StatusPasswordErr
			resp.StatusMsg = "password is incorrect"
			resp.UserId = -1
			return resp, nil
		}

		resp.StatusCode = internal.StatusSuccess
		resp.StatusMsg = "success"
		resp.UserId = UserInfoDb.Id
		return resp, nil
	case model.ErrNotFound:
		resp.StatusCode = internal.StatusUserNotExist
		resp.StatusMsg = "the email does not exist"
		resp.UserId = -1
		return resp, nil
	default:
		logx.Error("get user info from db error: ", err)
		resp.StatusCode = internal.StatusRpcErr
		resp.StatusMsg = "get user info from db error"
		return resp, nil
	}
}
