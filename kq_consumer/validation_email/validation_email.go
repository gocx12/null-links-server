package main

import (
	"bytes"
	"fmt"
	"github.com/jordan-wright/email"
	"github.com/zeromicro/go-queue/kq"
	"github.com/zeromicro/go-zero/core/conf"
	"net/http"
	"net/smtp"
	"strings"
	"text/template"

	"github.com/zeromicro/go-zero/core/logx"
)

func main() {
	var c kq.KqConf
	conf.MustLoad("./kq_consumer/validation_email/config.yaml", &c)

	RunValidationPageServer()

	q := kq.MustNewQueue(c, kq.WithHandle(sendEmail))
	defer q.Stop()
	q.Start()
}

func sendEmail(k, v string) error {
	// 邮件服务器地址和端口 发件人邮箱和密码
	smtpDomain := "smtp.163.com"
	sender := "null_links@163.com"
	password := "1Q2W3E4r5t"
	auth := smtp.PlainAuth("GWRLOVSZYVLRYQCH", sender, password, smtpDomain)

	// 收件人邮箱
	kqValue := v
	kqVaules := strings.Split(kqValue, "::")
	validationCode := kqVaules[0]
	recipient := kqVaules[1]

	// 邮件主题和内容
	subject := "Null-Links 注册验证码"
	// body := fmt.Sprintf("本邮件来自NULL Links, 您的验证码是: %s。请在5分钟内完成验证。", validationCode)
	// message := "From: " + sender + "\n\r" +
	// 	"To: " + recipient + "\n\r" +
	// 	"Subject: " + subject + "\n\r\n\r" +
	// 	body

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
	// err := smtp.SendMail(smtpDomain, auth, sender, []string{recipient}, []byte(message))
	// if err != nil {
	// 	logx.Error("Failed to send email:", err)
	// 	return err
	// }

	e := email.NewEmail()
	e.From = sender
	e.To = []string{recipient}
	e.Subject = subject
	e.HTML = body.Bytes()
	err = e.Send(smtpDomain+":25", auth)
	if err != nil {
		// TODO(chancyGao): 增加告警
		logx.Error("Failed to send email:", err, " recipient: ", recipient)
		return err
	}

	return nil
}

type TicketInfo struct {
	Picture        string
	Title          string
	Desc           string
	Warning        string
	ValidationCode string
}

func genHtml(ticketInfo TicketInfo) (*template.Template, error) {
	// 解析指定文件生成模板对象
	tmpl, err := template.ParseFiles("./kq_consumer/validation_email/validation_code_page.html")

	return tmpl, err
}

func renderHtml(responseWriter http.ResponseWriter, request *http.Request) {
	// 解析指定文件生成模板对象
	ticketInfo := TicketInfo{
		Picture:        "http://localhost:3000/static/null_link_logo.ico",
		Title:          "Null-Links 注册验证码",
		Desc:           "欢迎注册Null-Links。",
		Warning:        "请不要回复本邮件",
		ValidationCode: "abca",
	}
	tmpl, err := genHtml(ticketInfo)
	if err != nil {
		logx.Error("failed to render html:", err)
		return
	}
	tmpl.Execute(responseWriter, ticketInfo)
}

func RunValidationPageServer() {
	fs := http.FileServer(http.Dir("./kq_consumer/validation_email/"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))

	http.HandleFunc("/", renderHtml)
	err := http.ListenAndServe("127.0.0.1:3000", nil)
	if err != nil {
		fmt.Println("HTTP server failed,err:", err)
		return
	}
	fmt.Print("starting http://localhost:3000")
}
