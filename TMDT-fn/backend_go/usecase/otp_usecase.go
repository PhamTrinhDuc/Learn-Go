package usecase

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"errors"
	"fmt"
	"math/big"
	"time"

	"github.com/redis/go-redis/v9"
	"golang.org/x/crypto/bcrypt"
	"tmdt-backend/domain"
)

type otpUsecase struct {
	redis  *redis.Client
	mailer domain.MailerService
}

// NewOTPUsecase creates a new instance of OTPUsecase
func NewOTPUsecase(rClient *redis.Client, mailer domain.MailerService) domain.OTPUsecase {
	return &otpUsecase{
		redis:  rClient,
		mailer: mailer,
	}
}

func (u *otpUsecase) SendOTP(ctx context.Context, req *domain.SendOTPRequest) error {
	// 1. Generate 6-digit OTP
	otp, err := generateOTP()
	if err != nil {
		return fmt.Errorf("failed to generate OTP: %w", err)
	}

	// 2. Hash OTP
	hashedOtp, err := bcrypt.GenerateFromPassword([]byte(otp), 10)
	if err != nil {
		return fmt.Errorf("failed to hash OTP: %w", err)
	}

	// 3. Save hashed OTP to Redis (expiry 5 minutes)
	err = u.redis.Set(ctx, fmt.Sprintf("otp:%s", req.Email), hashedOtp, 300*time.Second).Err()
	if err != nil {
		return fmt.Errorf("failed to save OTP to redis: %w", err)
	}

	// 4. Send email
	err = u.mailer.SendOTP(ctx, req.Email, otp)
	if err != nil {
		return fmt.Errorf("failed to send OTP email: %w", err)
	}

	return nil
}

func (u *otpUsecase) VerifyOTP(ctx context.Context, req *domain.VerifyOTPRequest) (*domain.VerifyOTPResponse, error) {
	redisKey := fmt.Sprintf("otp:%s", req.Email)

	// 1. Fetch hashed OTP from Redis
	storedHash, err := u.redis.Get(ctx, redisKey).Result()
	if err == redis.Nil {
		return nil, errors.New("Mã OTP hết hạn")
	} else if err != nil {
		return nil, fmt.Errorf("failed to get OTP from redis: %w", err)
	}

	// 2. Compare OTP
	err = bcrypt.CompareHashAndPassword([]byte(storedHash), []byte(req.OTP))
	if err != nil {
		return nil, errors.New("Mã OTP không chính xác")
	}

	// 3. Generate secure reset token
	resetToken, err := generateResetToken()
	if err != nil {
		return nil, fmt.Errorf("failed to generate reset token: %w", err)
	}

	// 4. Store reset token in Redis (expiry 10 minutes)
	err = u.redis.Set(ctx, fmt.Sprintf("resetToken:%s", req.Email), resetToken, 600*time.Second).Err()
	if err != nil {
		return nil, fmt.Errorf("failed to save reset token to redis: %w", err)
	}

	// 5. Delete OTP from Redis
	_ = u.redis.Del(ctx, redisKey).Err()

	return &domain.VerifyOTPResponse{
		Success:    true,
		Message:    "Xác thực thành công!",
		ResetToken: resetToken,
	}, nil
}

func generateOTP() (string, error) {
	nBig, err := rand.Int(rand.Reader, big.NewInt(900000))
	if err != nil {
		return "", err
	}
	otp := nBig.Int64() + 100000
	return fmt.Sprintf("%d", otp), nil
}

func generateResetToken() (string, error) {
	b := make([]byte, 32)
	_, err := rand.Read(b)
	if err != nil {
		return "", err
	}
	return hex.EncodeToString(b), nil
}
