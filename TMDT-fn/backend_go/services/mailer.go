package services

import (
	"context"
	"fmt"
	"net/smtp"
	"os"
	"tmdt-backend/domain"
)

type mailerService struct{}

// NewMailerService creates a new instance of MailerService
func NewMailerService() domain.MailerService {
	return &mailerService{}
}

func (s *mailerService) SendOTP(ctx context.Context, toEmail string, otp string) error {
	from := os.Getenv("OTP_GMAIL")
	password := os.Getenv("OTP_PASSWORD")
	host := "smtp.gmail.com"
	port := "587"

	if from == "" || password == "" {
		return fmt.Errorf("OTP_GMAIL or OTP_PASSWORD environment variable is not configured")
	}

	subject := fmt.Sprintf("[OTP] %s là mã xác nhận của bạn", otp)
	body := getOTPTemplate(otp)

	// Set headers
	msg := []byte("To: " + toEmail + "\r\n" +
		"From: \"Hệ thống CellphoneX\" <" + from + ">\r\n" +
		"Subject: " + subject + "\r\n" +
		"MIME-Version: 1.0\r\n" +
		"Content-Type: text/html; charset=UTF-8\r\n" +
		"\r\n" +
		body + "\r\n")

	auth := smtp.PlainAuth("", from, password, host)
	err := smtp.SendMail(host+":"+port, auth, from, []string{toEmail}, msg)
	if err != nil {
		return fmt.Errorf("failed to send smtp email: %w", err)
	}

	return nil
}

func getOTPTemplate(otp string) string {
	return fmt.Sprintf(`<!DOCTYPE html>
<html>
<head>
    <meta charset="utf-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>Xác thực mã OTP</title>
    <style>
        body { font-family: 'Helvetica Neue', Helvetica, Arial, sans-serif; background-color: #f4f7f6; margin: 0; padding: 0; }
        .container { max-width: 600px; margin: 20px auto; background-color: #ffffff; border-radius: 8px; overflow: hidden; box-shadow: 0 4px 10px rgba(0,0,0,0.05); }
        .header { background-color: #4A90E2; padding: 30px; text-align: center; color: white; }
        .content { padding: 40px; text-align: center; line-height: 1.6; color: #333; }
        .otp-code { font-size: 36px; font-weight: bold; letter-spacing: 8px; color: #4A90E2; margin: 30px 0; padding: 15px; border: 2px dashed #4A90E2; border-radius: 10px; display: inline-block; background-color: #f0f7ff; }
        .footer { background-color: #f9fafb; padding: 20px; text-align: center; font-size: 12px; color: #999; }
        .warning { color: #e74c3c; font-size: 13px; margin-top: 20px; }
    </style>
</head>
<body>
<div class="container">
    <div class="header">
        <h1 style="margin:0;">Xác Thực Tài Khoản</h1>
    </div>
    <div class="content">
        <p style="font-size: 18px;">Chào bạn,</p>
        <p>Bạn đang thực hiện yêu cầu xác thực. Vui lòng sử dụng mã OTP dưới đây để hoàn tất quy trình:</p>

        <div class="otp-code">%s</div>

        <p>Mã này có hiệu lực trong <strong>5 phút</strong>.</p>
        <p class="warning">⚠️ Tuyệt đối không chia sẻ mã này với bất kỳ ai để bảo vệ tài khoản của bạn.</p>
    </div>
    <div class="footer">
        <p>© 2026 Tên Công Ty Của Bạn. All rights reserved.</p>
        <p>Hà Nội, Việt Nam</p>
    </div>
</div>
</body>
</html>`, otp)
}
