package user

import (
	"context"
	"net/http"
	"time"

	"nulltv/http_service/internal/svc"
	"nulltv/http_service/internal/types"
	"nulltv/rpc_service/user/pb/user"

	"github.com/golang-jwt/jwt"
	"github.com/zeromicro/go-zero/core/logc"
	"github.com/zeromicro/go-zero/core/logx"
)

type LoginLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

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
			StatusCode: http.StatusOK,
			StatusMsg:  "请输入用户名或邮箱",
			UserID:     -1,
		}
		err = nil
		return
	}

	respRpc, err := l.svcCtx.UserRpc.Login(l.ctx, &user.LoginReq{
		Username:  req.Username,
		UserEmail: req.UserEmail,
		Password:  req.Password,
	})
	if err != nil {
		resp = &types.LoginResp{
			StatusCode: http.StatusOK,
			StatusMsg:  "登录失败",
			UserID:     respRpc.UserId, // is -1
		}
		logc.Alert(l.ctx, "call UserRpc failed"+err.Error())
		err = nil
		return
	} else if respRpc.UserId == -1 {
		// the username does not exsit or the password is incorrect
		resp = &types.LoginResp{
			StatusCode: http.StatusOK,
			StatusMsg:  respRpc.StatusMsg,
			UserID:     respRpc.UserId, // is -1
		}
		err = nil
		return
	}

	secretKey := l.svcCtx.Config.Auth.AccessSecret
	iat := time.Now().Unix()
	seconds := l.svcCtx.Config.Auth.AccessExpire
	payload := respRpc.UserId // save user id in payload

	token, err := getJwtToken(secretKey, iat, seconds, payload)
	if err != nil {
		resp = &types.LoginResp{
			StatusCode: http.StatusOK,
			StatusMsg:  "Login fail",
			UserID:     respRpc.UserId, // is -1
		}
		logc.Alert(l.ctx, "getJwtToken() "+err.Error())
		err = nil
		return
	}

	resp = &types.LoginResp{
		StatusCode: http.StatusOK,
		StatusMsg:  respRpc.StatusMsg,
		UserID:     respRpc.UserId,
		Token:      token,
	}
	return
}

// @secretKey: JWT secret key
// @iat: time stamp
// @seconds: expire time(second)
// @payload: data payload
func getJwtToken(secretKey string, iat, seconds, payload int64) (string, error) {
	claims := make(jwt.MapClaims)
	claims["exp"] = iat + seconds
	claims["iat"] = iat
	claims["payload"] = payload
	token := jwt.New(jwt.SigningMethodHS256)
	token.Claims = claims
	return token.SignedString([]byte(secretKey))
}
