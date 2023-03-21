package helper

import (
	"OJPlatform/define"
	"crypto/tls"
	"fmt"
	"github.com/jordan-wright/email"
	"math/rand"
	"net/smtp"
	"strconv"
	"time"
)

// SendCodeMail 发送邮件
func SendCodeMail(toAddress, code string) error {
	e := email.NewEmail()
	e.From = "OJPlatform <343921998@qq.com>"
	e.To = []string{toAddress}
	e.Subject = "验证码发送测试"
	e.HTML = []byte("您的验证码：<b>" + code + "</b>")
	// 返回 EOF 时，关闭SSL重试
	return e.SendWithTLS("smtp.qq.com:465",
		smtp.PlainAuth("", "343921998@qq.com", define.MailPassword, "smtp.qq.com"),
		&tls.Config{InsecureSkipVerify: true, ServerName: "smtp.qq.com"})
}

// GenerateRandomCode 生成6位验证码
func GenerateRandomCode() string {
	rand.Seed(time.Now().UnixNano())
	code := ""
	for i := 0; i < 6; i++ {
		code += strconv.Itoa(rand.Intn(10))
	}
	fmt.Println(code)
	return code
}
