package test

import (
	"OJPlatform/define"
	"crypto/tls"
	"fmt"
	"github.com/jordan-wright/email"
	"math/rand"
	"net/smtp"
	"strconv"
	"testing"
	"time"
)

// TestSendEmail 测试邮件发送通路
func TestSendEmail(t *testing.T) {
	e := email.NewEmail()
	e.From = "Get <343921998@qq.com>"
	e.To = []string{"21S151155@stu.hit.edu.cn"}
	e.Subject = "验证码发送测试"
	e.HTML = []byte("您的验证码：<b>123456</b>")
	// 返回 EOF 时，关闭SSL重试
	err := e.SendWithTLS("smtp.qq.com:465",
		smtp.PlainAuth("", "343921998@qq.com", define.MailPassword, "smtp.qq.com"),
		&tls.Config{InsecureSkipVerify: true, ServerName: "smtp.qq.com"})
	if err != nil {
		t.Fatal(err)
	}
}

// TestGenerateRandomCode 测试生成六位随机验证码
func TestGenerateRandomCode(t *testing.T) {
	rand.Seed(time.Now().UnixNano())
	code := ""
	for i := 0; i < 6; i++ {
		code += strconv.Itoa(rand.Intn(10))
	}
	fmt.Println(code)
}
