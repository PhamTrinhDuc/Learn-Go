package controller

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"tmdt-backend/domain"
)

type OTPController struct {
	otpUsecase domain.OTPUsecase
}

// NewOTPController creates a new instance of OTPController
func NewOTPController(otpUsecase domain.OTPUsecase) *OTPController {
	return &OTPController{
		otpUsecase: otpUsecase,
	}
}

// SendOTP handles POST /api/otp/send-otp
func (ctrl *OTPController) SendOTP(c *gin.Context) {
	var req domain.SendOTPRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Vui lòng nhập email"})
		return
	}

	err := ctrl.otpUsecase.SendOTP(c.Request.Context(), &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Lỗi gửi OTP"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Mã OTP đã được gửi!"})
}

// VerifyOTP handles POST /api/otp/verify-otp
func (ctrl *OTPController) VerifyOTP(c *gin.Context) {
	var req domain.VerifyOTPRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Thiếu email hoặc OTP"})
		return
	}

	resp, err := ctrl.otpUsecase.VerifyOTP(c.Request.Context(), &req)
	if err != nil {
		errMsg := err.Error()
		if errMsg == "Mã OTP hết hạn" || errMsg == "Mã OTP không chính xác" {
			c.JSON(http.StatusBadRequest, gin.H{"message": errMsg})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Lỗi xác thực"})
		return
	}

	c.JSON(http.StatusOK, resp)
}
