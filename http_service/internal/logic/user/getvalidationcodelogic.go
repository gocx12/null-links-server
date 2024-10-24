package user

import (
	"bytes"
	"context"
	"math/rand"
	"net/smtp"
	"text/template"
	"time"

	"null-links/http_service/internal/common"
	"null-links/http_service/internal/svc"
	"null-links/http_service/internal/types"
	"null-links/internal"

	"github.com/jordan-wright/email"
	"github.com/zeromicro/go-zero/core/logx"
)

type GetValidationCodeLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetValidationCodeLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetValidationCodeLogic {
	return &GetValidationCodeLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetValidationCodeLogic) GetValidationCode(req *types.GetValidationCodeReq) (resp *types.GetValidationCodeResp, err error) {
	// respRpc, err := l.svcCtx.UserRpc.GetValidtaionCode(l.ctx, &user.GetValidtaionCodeReq{
	// 	Email: req.Email,
	// })
	// if err != nil {
	// 	logx.Error("call UserRpc failed, err: ", err)
	// 	resp = &types.GetValidationCodeResp{
	// 		StatusCode: internal.StatusRpcErr,
	// 		StatusMsg:  "获取验证码失败",
	// 	}
	// 	err = nil
	// 	return
	// } else if respRpc.StatusCode != internal.StatusSuccess {
	// 	logx.Error("call UserRpc failed, err: ", resp.StatusMsg)
	// 	resp = &types.GetValidationCodeResp{
	// 		StatusCode: internal.StatusRpcErr,
	// 		StatusMsg:  "获取验证码失败",
	// 	}
	// 	return
	// }

	logx.Debug("get validation code, email: ", req.Email)

	recipient := req.Email
	validationCode := l.generateValidationCode()

	// Redis存储验证码, 10分钟过期
	_, err = l.svcCtx.RedisClient.SetEx(l.ctx, common.RdsKeyEmailValidationPre+"_"+recipient, validationCode, 10*time.Minute).Result()
	if err != nil {
		logx.Error("failed to set email validation code to redis, error:", err, " email: ", recipient)
		resp = &types.GetValidationCodeResp{
			StatusCode: internal.StatusRpcErr,
			StatusMsg:  "set email validation code to redis failed",
		}
		return
	}

	// // kafka pusher
	// data := recipient + "::" + validationCode
	// if err := l.svcCtx.VdEmailKqPusher.Push(data); err != nil {
	// 	logx.Error("VdEmailKqPusher Push error:", err)
	// 	resp = &types.GetValidationCodeResp{
	// 		StatusCode: internal.StatusGatewayErr,
	// 		StatusMsg:  "push email validation code to kq failed, err: " + err.Error(),
	// 	}
	// 	return resp, nil
	// }

	go func() {
		l.sendEmail(recipient, validationCode)
	}()

	resp = &types.GetValidationCodeResp{
		StatusCode: internal.StatusSuccess,
		StatusMsg:  "成功",
	}
	return
}

func (l *GetValidationCodeLogic) generateValidationCode() string {
	// go 1.20 起弃用rand.Seed, 即使不用seed也可以产生不同的数字
	chars := []rune("ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz")
	length := 4
	buf := make([]rune, length)
	for i := range buf {
		buf[i] = chars[rand.Intn(len(chars))]
	}
	return string(buf)
}

// 验证码邮件模板
type TicketInfo struct {
	Picture        string
	Title          string
	Desc           string
	Warning        string
	ValidationCode string
}

func (l *GetValidationCodeLogic) sendEmail(recipient, validationCode string) error {
	logx.Info("send email to: ", recipient, ", validation code: ", validationCode)
	// 邮件服务器地址和端口 发件人邮箱和密码
	emailConf := l.svcCtx.Config.Email
	smtpDomain := emailConf.SmtpDomain
	sender := emailConf.Sender
	password := emailConf.Password
	auth := smtp.PlainAuth("", sender, password, smtpDomain)

	// 邮件主题和内容
	subject := "Null-Links 注册验证码"

	ticketInfo := TicketInfo{
		Picture:        "http://localhost:3000/static/logo.png",
		Title:          "Null-Links 注册验证码",
		Desc:           "",
		Warning:        "请不要回复本邮件",
		ValidationCode: validationCode,
	}
	tmpl, err := genHtml(ticketInfo)
	if err != nil {
		logx.Error("failed to render html:", err)
		return err
	}
	body := new(bytes.Buffer)
	tmpl.Execute(body, ticketInfo)

	// 发送邮件
	e := email.NewEmail()
	e.From = sender
	e.To = []string{recipient}
	e.Subject = subject
	e.HTML = body.Bytes()
	err = e.Send(smtpDomain+":25", auth)
	if err != nil {
		// TODO(chancyGao): 增加告警
		logx.Error("Failed to send email. err=", err, ", recipient=", recipient)
		return err
	}

	return nil
}

func genHtml(ticketInfo TicketInfo) (*template.Template, error) {
	// 解析指定文件生成模板对象
	tmpl, err := template.ParseFiles("./etc/validation_code_page.html")

	return tmpl, err
}
