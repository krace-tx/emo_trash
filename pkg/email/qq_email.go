package email

import (
	"bytes"
	"crypto/tls"
	"fmt"
	"html/template"
	"net/smtp"
	"regexp"
)

type Pop3 struct {
	Id   string
	Auth string
	Host string
	Port string
}

type codeMailData struct {
	Title       string
	SceneText   string
	Code        string
	ExpireText  string
	SupportText string
}

// 发送验证码
func SendCode(smt Pop3, to string, vcode string, scene string) error {
	// 发件人邮箱地址和授权码
	from := smt.Id
	password := smt.Auth
	smtpHost := smt.Host
	smtpPort := smt.Port

	// 收件人
	toEmail := to

	subjectText, sceneText := buildSubjectAndScene(scene)
	// 构建邮件内容
	subject := fmt.Sprintf("Subject: %s\r\n", subjectText)
	fromHeader := fmt.Sprintf("From: %s\r\n", from)
	toHeader := fmt.Sprintf("To: %s\r\n", toEmail)
	body, err := renderCodeMailHTML(codeMailData{
		Title:       "Emo Trash 账号安全验证",
		SceneText:   sceneText,
		Code:        vcode,
		ExpireText:  "验证码 5 分钟内有效，请勿泄露给他人。",
		SupportText: "若非本人操作，请忽略本邮件。",
	})
	if err != nil {
		return err
	}
	message := fromHeader + toHeader + subject + "Content-Type: text/html; charset=UTF-8\r\n\r\n" + body

	// 设置 SMTP 服务器配置
	auth := smtp.PlainAuth("", from, password, smtpHost)

	// TLS 配置
	tlsconfig := &tls.Config{
		InsecureSkipVerify: true,
		ServerName:         smtpHost,
	}

	// 连接到 SMTP 服务器
	conn, err := tls.Dial("tcp", smtpHost+":"+smtpPort, tlsconfig)
	if err != nil {
		return err
	}
	defer conn.Close()

	client, err := smtp.NewClient(conn, smtpHost)
	if err != nil {
		return err
	}
	defer client.Quit()

	// 认证
	if err = client.Auth(auth); err != nil {
		return err
	}

	// 设置发件人和收件人
	if err = client.Mail(from); err != nil {
		return err
	}

	if err = client.Rcpt(toEmail); err != nil {
		return err
	}

	// 发送邮件内容
	writer, err := client.Data()
	if err != nil {
		return err
	}

	_, err = writer.Write([]byte(message))
	if err != nil {
		return err
	}

	err = writer.Close()
	if err != nil {
		return err
	}

	return nil
}

func buildSubjectAndScene(scene string) (subject string, sceneText string) {
	switch scene {
	case "register":
		return "Emo Trash 注册验证码", "您正在进行账号注册，请使用以下验证码完成验证："
	case "login":
		return "Emo Trash 登录验证码", "您正在进行账号登录，请使用以下验证码完成验证："
	case "reset_pwd":
		return "Emo Trash 重置密码验证码", "您正在进行密码重置，请使用以下验证码完成验证："
	default:
		return "Emo Trash 安全验证码", "您正在进行安全验证，请使用以下验证码完成操作："
	}
}

func renderCodeMailHTML(data codeMailData) (string, error) {
	const mailTpl = `
<div style="font-family: -apple-system,BlinkMacSystemFont,'Segoe UI',Roboto,Arial,sans-serif; max-width: 640px; margin: 0 auto; padding: 24px; border: 1px solid #E5E7EB; border-radius: 12px; background: #FFFFFF;">
  <h2 style="margin: 0 0 16px; color: #111827;">{{.Title}}</h2>
  <p style="margin: 0 0 12px; color: #374151; font-size: 15px;">您好，</p>
  <p style="margin: 0 0 20px; color: #374151; font-size: 15px;">{{.SceneText}}</p>
  <div style="margin: 0 0 20px; text-align: center;">
    <span style="display: inline-block; letter-spacing: 6px; font-size: 30px; font-weight: 700; color: #2563EB; background: #EFF6FF; border: 1px dashed #93C5FD; border-radius: 10px; padding: 10px 20px;">{{.Code}}</span>
  </div>
  <p style="margin: 0 0 8px; color: #6B7280; font-size: 13px;">{{.ExpireText}}</p>
  <p style="margin: 0; color: #9CA3AF; font-size: 13px;">{{.SupportText}}</p>
</div>`

	tpl, err := template.New("code-mail").Parse(mailTpl)
	if err != nil {
		return "", err
	}

	var buf bytes.Buffer
	if err = tpl.Execute(&buf, data); err != nil {
		return "", err
	}
	return buf.String(), nil
}

// 验证电子邮件地址的合法性
func IsValidEmail(email string) bool {
	// 电子邮件地址的正则表达式
	emailRegex := `^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`
	re := regexp.MustCompile(emailRegex)
	return re.MatchString(email)
}
