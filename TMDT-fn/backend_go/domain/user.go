package domain

import (
	"context"
	"time"
)

type User struct {
	ID         int        `json:"id"`
	Password   string     `json:"password,omitempty"` // omitempty to hide password in json responses
	FullName   string     `json:"full_name"`
	Email      string     `json:"email"`
	NumPhone   string     `json:"num_phone"`
	Role       string     `json:"role"`
	IsLock     bool       `json:"is_lock"`
	JoinedDate time.Time  `json:"joined_date"`
	Gender     string     `json:"gender"`
	DOB        *time.Time `json:"dob"` // pointer to allow null value
}

// UserLoginRequest input for Login
type UserLoginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

// UserRegisterRequest input for Register
type UserRegisterRequest struct {
	FullName string `json:"fullName" binding:"required"`
	Email    string `json:"email" binding:"required,email"`
	NumPhone string `json:"numPhone" binding:"required"`
	Password string `json:"password" binding:"required,min=6"`
}

// UserUpdateRequest input for updating profile
type UserUpdateRequest struct {
	FullName string  `json:"fullName" binding:"required"`
	Email    string  `json:"email" binding:"required,email"`
	NumPhone string  `json:"numPhone" binding:"required"`
	DOB      *string `json:"dob"`
	Gender   string  `json:"gender" binding:"required"`
}

// CheckPasswordValidityRequest input
type CheckPasswordValidityRequest struct {
	Email       string `json:"email" binding:"required"`
	NewPassword string `json:"newPassword" binding:"required"`
}

// ResetPasswordRequest input
type ResetPasswordRequest struct {
	Email       string `json:"email" binding:"required"`
	ResetToken  string `json:"resetToken" binding:"required"`
	NewPassword string `json:"newPassword" binding:"required,min=6"`
}

// CheckEmailRequest input
type CheckEmailRequest struct {
	Email string `json:"email" binding:"required,email"`
}

// GoogleLoginRequest input
type GoogleLoginRequest struct {
	IDToken string `json:"idToken" binding:"required"`
}

// AdminResetPasswordRequest input
type AdminResetPasswordRequest struct {
	Password string `json:"password" binding:"required,min=6"`
}

// UserRepository defines the methods a repository must implement
type UserRepository interface {
	FetchAll(ctx context.Context) ([]User, error)
	GetByID(ctx context.Context, id int) (*User, error)
	GetByEmailOrPhone(ctx context.Context, username string) (*User, error)
	GetByEmail(ctx context.Context, email string) (*User, error)
	GetByPhone(ctx context.Context, phone string) (*User, error)
	Create(ctx context.Context, user *User) error
	Update(ctx context.Context, user *User) error
	UpdateRole(ctx context.Context, id int, role string) error
	UpdateStatus(ctx context.Context, id int, isLock bool) error
	UpdatePassword(ctx context.Context, email string, hashedPassword string) error
	GetDistinctRoles(ctx context.Context) ([]string, int, error)
}

// UserUsecase defines the use cases for User module
type UserUsecase interface {
	Register(ctx context.Context, req *UserRegisterRequest) (*User, error)
	Login(ctx context.Context, req *UserLoginRequest) (*User, error)
	GoogleLogin(ctx context.Context, idToken string) (*User, error)
	GetProfile(ctx context.Context, id int) (*User, error)
	GetAllUsers(ctx context.Context) ([]User, error)
	GetRoles(ctx context.Context) ([]string, int, error)
	UpdateProfile(ctx context.Context, id int, req *UserUpdateRequest) (*User, error)
	UpdateRole(ctx context.Context, id int, role string) (*User, error)
	UpdateStatus(ctx context.Context, id int, isLock bool) (*User, error)
	CheckPasswordValidity(ctx context.Context, email string, newPassword string) (bool, error)
	ResetPassword(ctx context.Context, req *ResetPasswordRequest) error
	CheckEmail(ctx context.Context, email string) (*User, error)
	AdminResetPassword(ctx context.Context, id int, newPassword string) (*User, error)
}
