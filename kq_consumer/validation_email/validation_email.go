package main

import (
	"bytes"
	"flag"
	"net/http"
	"net/smtp"
	"strings"
	"text/template"

	"github.com/jordan-wright/email"
	"github.com/zeromicro/go-queue/kq"
	"github.com/zeromicro/go-zero/core/conf"

	"github.com/zeromicro/go-zero/core/logx"
)

type Email struct {
	SmtpDomain string
	Sender     string
	Password   string
}
type Conf struct {
	Mode               string
	Email              Email
	VdEmailKqConsumser kq.KqConf
}

var c Conf

// 验证码邮件模板
type TicketInfo struct {
	Picture        string
	Title          string
	Desc           string
	Warning        string
	ValidationCode string
}

var configFile = flag.String("f", "kq_consumer/validation_email/config.yaml", "the config file")

func main() {
	flag.Parse()
	conf.MustLoad(*configFile, &c)

	// if c.Mode == "debug" {
	// 	logx.Debug("start server")
	// 	runVdEmailPageServer()
	// }

	logx.Info("validation email kq consumer starting")
	q := kq.MustNewQueue(c.VdEmailKqConsumser, kq.WithHandle(sendEmail))
	defer q.Stop()
	q.Start()
}

func sendEmail(k, v string) error {
	logx.Debug("sendEmail", "k: ", k, " v: ", v)
	// 邮件服务器地址和端口 发件人邮箱和密码
	smtpDomain := c.Email.SmtpDomain
	sender := c.Email.Sender
	password := c.Email.Password
	auth := smtp.PlainAuth("", sender, password, smtpDomain)

	// 收件人邮箱与验证码
	kqValue := v
	kqVaules := strings.Split(kqValue, "::")
	recipient := kqVaules[0]
	validationCode := kqVaules[1]

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
		logx.Error("Failed to send email:", err, " ,recipient: ", recipient)
		return err
	}

	return nil
}

func genHtml(ticketInfo TicketInfo) (*template.Template, error) {
	// 解析指定文件生成模板对象
	tmpl, err := template.ParseFiles("./kq_consumer/validation_email/validation_code_page.html")

	return tmpl, err
}

func runVdEmailPageServer() {
	http.HandleFunc("/", renderHtml)
	err := http.ListenAndServe("10.63.180.57:3002", nil)

	if err != nil {
		logx.Error("HTTP server failed, err: ", err)
		return
	}
	logx.Info("starting http://localhost:3002")
}

func renderHtml(responseWriter http.ResponseWriter, request *http.Request) {
	// 解析指定文件生成模板对象
	ticketInfo := TicketInfo{
		ValidationCode: "abca",
	}
	tmpl, err := genHtml(ticketInfo)
	if err != nil {
		logx.Error("failed to render html:", err)
		return
	}
	tmpl.Execute(responseWriter, ticketInfo)
}
