package main

import (
	"fmt"
	"strings"

	"github.com/zeromicro/go-queue/kq"
	"github.com/zeromicro/go-zero/core/conf"
	"net/smtp"

	"github.com/zeromicro/go-zero/core/logx"
)

func main() {
	var c kq.KqConf
	conf.MustLoad("config.yaml", &c)

	q := kq.MustNewQueue(c, kq.WithHandle(sendEmail))
	defer q.Stop()
	q.Start()
}

func sendEmail(k, v string) error {
	// 邮件服务器地址和端口
	smtpDomain := "smtp.163.com"

	// 发件人邮箱和密码
	sender := "null_links@163.com"
	password := "1Q2W3E4r5t"

	// 收件人邮箱
	kqValue := v
	kqVaules := strings.Split(kqValue, "::")
	validationCode := kqVaules[0]
	recipient := kqVaules[1]

	// 邮件主题和内容
	subject := "NULL Links 验证码"
	body := fmt.Sprint("本邮件来自NULL Links, 您的验证码是: %s。请在5分钟内完成验证。", validationCode)

	message := "From: " + sender + "\n\r" +
		"To: " + recipient + "\n\r" +
		"Subject: " + subject + "\n\r\n\r" +
		body

	// 认证信息
	auth := smtp.PlainAuth("GWRLOVSZYVLRYQCH", sender, password, smtpDomain)

	// 发送邮件
	err := smtp.SendMail(smtpDomain, auth, sender, []string{recipient}, []byte(message))
	if err != nil {
		logx.Error("Failed to send email:", err)
		return err
	}
	return nil
}
