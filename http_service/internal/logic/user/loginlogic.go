package user

import (
	"context"
	"encoding/base64"
	"time"

	"null-links/cron/model"
	"null-links/http_service/internal/svc"
	"null-links/http_service/internal/types"
	"null-links/internal"

	"github.com/zeromicro/go-zero/core/logx"
	"golang.org/x/crypto/scrypt"
)

type LoginLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

var PW_HASH_BYTES = 32
var SALT = []byte{126, 145, 58, 233, 153, 107, 4, 231}

func NewLoginLogic(ctx context.Context, svcCtx *svc.ServiceContext) *LoginLogic {
	return &LoginLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *LoginLogic) Login(req *types.LoginReq) (resp *types.LoginResp, err error) {
	// check whether both username and email are empty
	if req.Username == "" && req.UserEmail == "" {
		resp = &types.LoginResp{
			StatusCode: internal.StatusParamErr,
			StatusMsg:  "please input your name or your email address",
			UserID:     -1,
		}
		err = nil
		return
	}
	if req.Password == "" {
		resp = &types.LoginResp{
			StatusCode: internal.StatusParamErr,
			StatusMsg:  "please input your password",
			UserID:     -1,
		}
		err = nil
		return
	}

	resp = &types.LoginResp{}

	// 查询数据库，获取用户信息
	UserInfoDb, err := l.svcCtx.UserModel.FindPasswordByEmail(l.ctx, req.UserEmail)

	switch err {
	case nil:
		hash, err := scrypt.Key([]byte(req.Password), SALT, 1<<15, 8, 1, PW_HASH_BYTES)
		encodedHash := base64.StdEncoding.EncodeToString(hash)
		if err != nil {
			logx.Error("scrypt encode password error: " + err.Error())
			resp.StatusCode = internal.StatusRpcErr
			resp.StatusMsg = ""
			resp.UserID = -1
			return nil, err
		}
		if UserInfoDb.Password != encodedHash {
			resp.StatusCode = internal.StatusPasswordErr
			resp.StatusMsg = "password is incorrect"
			resp.UserID = -1
			return resp, nil
		}
		resp.StatusCode = internal.StatusSuccess
		resp.StatusMsg = "success"
		resp.UserID = UserInfoDb.Id
		resp.Username = UserInfoDb.Username
		resp.AvatarUrl = UserInfoDb.AvatarUrl

		secretKey := l.svcCtx.Config.Auth.AccessSecret
		iat := time.Now().Unix()
		seconds := l.svcCtx.Config.Auth.AccessExpire
		payload := resp.UserID // save user id in payload

		// 生成token
		token, err := internal.GenJwtToken(secretKey, iat, seconds, payload)
		if err != nil {
			logx.Error("generate token err:", err, " ,token:", token)
			resp.StatusCode = internal.StatusGatewayErr
			resp.StatusMsg = "login error"
			err = nil
			return resp, nil
		}
		resp.Token = token

		return resp, nil

	case model.ErrNotFound:
		resp.StatusCode = internal.StatusUserNotExist
		resp.StatusMsg = "the email does not exist"
		resp.UserID = -1
		return resp, nil
	default:
		logx.Error("get user info from db error: ", err)
		resp.StatusCode = internal.StatusRpcErr
		resp.StatusMsg = "get user info from db error"
		return resp, nil
	}
}
