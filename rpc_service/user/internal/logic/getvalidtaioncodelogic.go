package logic

import (
	"context"
	"log"
	"net/smtp"

	"null-links/rpc_service/user/internal/svc"
	"null-links/rpc_service/user/pb/user"

	"github.com/zeromicro/go-zero/core/logx"
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

func (l *GetValidtaionCodeLogic) GetValidtaionCode(in *user.GetValidtaionCodeReq) (*user.GetValidtaionCodeResp, error) {
	email := in.Email
	logx.Info("GetValidtaionCodeLogic.GetValidtaionCode", "email: ", email)

	// 邮件服务器地址和端口
	smtpDomain := "smtp.163.com"

	// 发件人邮箱和密码
	sender := "null_links@163.com"
	password := "1Q2W3E4r5t"

	// 收件人邮箱
	recipient := in.Email

	// 邮件主题和内容
	// 随机生成四位字母验证码
	subject := "NULL Links 验证码"
	body := "本邮件来自NULL Links, 您的验证码是: ABCD。请在5分钟内完成验证。"

	// 构建邮件内容
	message := "From: " + sender + "\n" +
		"To: " + recipient + "\n" +
		"Subject: " + subject + "\n\n" +
		body

	// 认证信息
	auth := smtp.PlainAuth("GWRLOVSZYVLRYQCH", sender, password, smtpDomain)

	// 发送邮件
	err := smtp.SendMail(smtpDomain, auth, sender, []string{recipient}, []byte(message))
	if err != nil {
		log.Fatal("Failed to send email:", err)
		return &user.GetValidtaionCodeResp{
			StatusCode: 0,
			StatusMsg:  "send validation code email failed, err: " + err.Error(),
		}, nil
	}

	return &user.GetValidtaionCodeResp{
		StatusCode: 1,
		StatusMsg:  "success",
	}, nil
}
