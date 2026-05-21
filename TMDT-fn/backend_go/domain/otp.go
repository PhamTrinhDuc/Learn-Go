package domain

import "context"

type SendOTPRequest struct {
	Email string `json:"email" binding:"required,email"`
}

type VerifyOTPRequest struct {
	Email string `json:"email" binding:"required,email"`
	OTP   string `json:"otp" binding:"required"`
}

type VerifyOTPResponse struct {
	Success    bool   `json:"success"`
	Message    string `json:"message"`
	ResetToken string `json:"resetToken"`
}

type OTPUsecase interface {
	SendOTP(ctx context.Context, req *SendOTPRequest) error
	VerifyOTP(ctx context.Context, req *VerifyOTPRequest) (*VerifyOTPResponse, error)
}

type MailerService interface {
	SendOTP(ctx context.Context, toEmail string, otp string) error
}
