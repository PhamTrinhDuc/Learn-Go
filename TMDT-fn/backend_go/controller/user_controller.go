package controller

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"tmdt-backend/domain"
)

type UserController struct {
	userUsecase domain.UserUsecase
}

// NewUserController creates a new instance of UserController
func NewUserController(u domain.UserUsecase) *UserController {
	return &UserController{userUsecase: u}
}

// Register handles POST /api/user/register
func (ctrl *UserController) Register(c *gin.Context) {
	var req domain.UserRegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Vui lòng điền đầy đủ tất cả các trường hợp lệ!"})
		return
	}

	user, err := ctrl.userUsecase.Register(c.Request.Context(), &req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "Đăng ký thành công!",
		"user":    user,
	})
}

// Login handles POST /api/user/login
func (ctrl *UserController) Login(c *gin.Context) {
	var req domain.UserLoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Thiếu tài khoản hoặc mật khẩu!"})
		return
	}

	user, err := ctrl.userUsecase.Login(c.Request.Context(), &req)
	if err != nil {
		// Differentiate error status codes based on message to match Node.js
		errMsg := err.Error()
		if errMsg == "tài khoản không tồn tại" || errMsg == "mật khẩu không chính xác" {
			c.JSON(http.StatusUnauthorized, gin.H{"message": "Tài khoản không tồn tại hoặc mật khẩu không chính xác!"})
			return
		} else if errMsg == "tài khoản này đã bị khóa bởi Admin" {
			c.JSON(http.StatusForbidden, gin.H{"message": "Tài khoản này đã bị khóa bởi Admin!"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Lỗi hệ thống!"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Đăng nhập thành công!",
		"user":    user,
	})
}

// GoogleLogin handles POST /api/user/google/login
func (ctrl *UserController) GoogleLogin(c *gin.Context) {
	var req domain.GoogleLoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Thiếu mã xác thực Google (idToken)!"})
		return
	}

	user, err := ctrl.userUsecase.GoogleLogin(c.Request.Context(), req.IDToken)
	if err != nil {
		errMsg := err.Error()
		if errMsg == "tài khoản Google này đã bị khóa" {
			c.JSON(http.StatusForbidden, gin.H{"message": "Tài khoản Google này đã bị khóa!"})
			return
		}
		c.JSON(http.StatusUnauthorized, gin.H{"message": "Xác thực Google thất bại!"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Đăng nhập Google thành công!",
		"user":    user,
	})
}

// GetProfile handles GET /api/user/information/:id
func (ctrl *UserController) GetProfile(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "ID người dùng không hợp lệ!"})
		return
	}

	user, err := ctrl.userUsecase.GetProfile(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"message": "Không tìm thấy người dùng!"})
		return
	}

	c.JSON(http.StatusOK, user)
}

// ListUsers handles GET /api/user/
func (ctrl *UserController) ListUsers(c *gin.Context) {
	users, err := ctrl.userUsecase.GetAllUsers(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Lỗi máy chủ nội bộ!"})
		return
	}
	c.JSON(http.StatusOK, users)
}

// GetRoles handles GET /api/user/role
func (ctrl *UserController) GetRoles(c *gin.Context) {
	roles, count, err := ctrl.userUsecase.GetRoles(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"message": "Không tìm thấy quyền người dùng!"})
		return
	}

	// Format roleNum as an array of objects: [ { count: N } ] to match SQL output format in Node.js
	roleNum := []gin.H{
		{"count": strconv.Itoa(count)},
	}

	c.JSON(http.StatusOK, gin.H{
		"roles":   roles,
		"roleNum": roleNum,
	})
}

// UpdateProfile handles PUT /api/user/information/:id
func (ctrl *UserController) UpdateProfile(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "ID người dùng không hợp lệ!"})
		return
	}

	var req domain.UserUpdateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Dữ liệu cập nhật không hợp lệ!"})
		return
	}

	user, err := ctrl.userUsecase.UpdateProfile(c.Request.Context(), id, &req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Đã lưu thay đổi thành công!",
		"user":    user,
	})
}

// UpdateRole handles PUT /api/user/role/:id
func (ctrl *UserController) UpdateRole(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "ID người dùng không hợp lệ!"})
		return
	}

	var body struct {
		Role string `json:"role" binding:"required"`
	}
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Vui lòng cung cấp quyền (role) mới!"})
		return
	}

	user, err := ctrl.userUsecase.UpdateRole(c.Request.Context(), id, body.Role)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"message": "Không tìm thấy người dùng để cập nhật!"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Cập nhật quyền thành công!",
		"user":    user,
	})
}

// UpdateStatus handles PUT /api/user/status/:id
func (ctrl *UserController) UpdateStatus(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "ID người dùng không hợp lệ!"})
		return
	}

	var body struct {
		IsLock *bool `json:"is_lock" binding:"required"`
	}
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Vui lòng cung cấp trạng thái khóa (is_lock)!"})
		return
	}

	user, err := ctrl.userUsecase.UpdateStatus(c.Request.Context(), id, *body.IsLock)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"message": "Không tìm thấy người dùng để cập nhật trạng thái!"})
		return
	}

	actionStr := "mở khóa"
	if *body.IsLock {
		actionStr = "khóa"
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Tiến hành " + actionStr + " tài khoản thành công!",
		"user":    user,
	})
}

// CheckPasswordValidity handles POST /api/user/check-password-validity
func (ctrl *UserController) CheckPasswordValidity(c *gin.Context) {
	var req domain.CheckPasswordValidityRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Thiếu Email hoặc Mật khẩu mới!"})
		return
	}

	isValid, err := ctrl.userUsecase.CheckPasswordValidity(c.Request.Context(), req.Email, req.NewPassword)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"message": "Người dùng không tồn tại!"})
		return
	}

	if !isValid {
		c.JSON(http.StatusBadRequest, gin.H{
			"isValid": false,
			"message": "Mật khẩu mới không được trùng với mật khẩu hiện tại!",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"isValid": true,
		"message": "Mật khẩu hợp lệ.",
	})
}

// ResetPassword handles POST /api/user/reset-password
func (ctrl *UserController) ResetPassword(c *gin.Context) {
	var req domain.ResetPasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Vui lòng cung cấp đầy đủ Email, Reset Token và Mật khẩu mới!"})
		return
	}

	err := ctrl.userUsecase.ResetPassword(c.Request.Context(), &req)
	if err != nil {
		errMsg := err.Error()
		if errMsg == "yêu cầu không hợp lệ hoặc phiên làm việc đã hết hạn. Vui lòng xác thực lại OTP" {
			c.JSON(http.StatusForbidden, gin.H{"message": errMsg})
			return
		}
		c.JSON(http.StatusBadRequest, gin.H{"message": errMsg})
		return
	}

	// Retrieve user profile to return in success JSON
	user, err := ctrl.userUsecase.CheckEmail(c.Request.Context(), req.Email)
	var responseUser gin.H
	if err == nil {
		responseUser = gin.H{
			"id":        user.ID,
			"full_name": user.FullName,
			"email":     user.Email,
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Mật khẩu của bạn đã được đặt lại thành công!",
		"user":    responseUser,
	})
}

// CheckEmail handles POST /api/user/check-email
func (ctrl *UserController) CheckEmail(c *gin.Context) {
	var req domain.CheckEmailRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Vui lòng cung cấp địa chỉ email!"})
		return
	}

	user, err := ctrl.userUsecase.CheckEmail(c.Request.Context(), req.Email)
	if err != nil {
		errMsg := err.Error()
		if errMsg == "email này chưa được đăng ký trong hệ thống" {
			c.JSON(http.StatusNotFound, gin.H{
				"exists":  false,
				"message": "Email này chưa được đăng ký trong hệ thống!",
			})
			return
		} else if errMsg == "tài khoản này hiện đang bị khóa. Vui lòng liên hệ quản trị viên" {
			c.JSON(http.StatusForbidden, gin.H{
				"exists":  true,
				"is_lock": true,
				"message": "Tài khoản này hiện đang bị khóa. Vui lòng liên hệ quản trị viên!",
			})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Lỗi hệ thống khi kiểm tra email!"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"exists":  true,
		"is_lock": false,
		"message": "Email hợp lệ, chuẩn bị gửi mã xác thực!",
		"user": gin.H{
			"full_name": user.FullName,
			"email":     user.Email,
		},
	})
}

// AdminResetPassword handles PUT /api/user/admin-reset-password/:id
func (ctrl *UserController) AdminResetPassword(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "ID người dùng không hợp lệ!"})
		return
	}

	var req domain.AdminResetPasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Mật khẩu phải có ít nhất 6 ký tự!"})
		return
	}

	user, err := ctrl.userUsecase.AdminResetPassword(c.Request.Context(), id, req.Password)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Đặt lại mật khẩu người dùng thành công!",
		"user":    user,
	})
}
