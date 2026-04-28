package domain

import (
	"backend/internal/utils"
	"context"
	"time"

	"github.com/google/uuid"
)

// UserRole defines the allowed roles for a user
type UserRole string

const (
	RoleCustomer UserRole = "customer"
	RoleManager  UserRole = "manager"
	RoleStylist  UserRole = "stylist"
	RoleAdmin    UserRole = "admin"
)

// User represents the users table in the database
type User struct {
	ID                uuid.UUID  `json:"id" db:"id"`
	Email             string     `json:"email" db:"email"`
	Password          string     `json:"-" db:"password"`
	Name              string     `json:"name" db:"name"`
	Phone             string     `json:"phone" db:"phone"`
	Birthday          *time.Time `json:"birthday,omitempty" db:"birthday"`
	Address           *string    `json:"address,omitempty" db:"address"`
	Role              UserRole   `json:"role" db:"role"`
	LoyaltyPoints     int        `json:"loyalty_points" db:"loyalty_points"`
	PreferredBranchID *uuid.UUID `json:"preferred_branch_id,omitempty" db:"preferred_branch_id"`
	LastVisitAt       *time.Time `json:"last_visit_at,omitempty" db:"last_visit_at"`
	IsActive          bool       `json:"is_active" db:"is_active"`
	CreatedAt         time.Time  `json:"created_at" db:"created_at"`
	UpdatedAt         time.Time  `json:"updated_at" db:"updated_at"`
}

// LoginRequest represents login payload
type LoginRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=6"`
}

// LoginResponse represents successful login response
type LoginResponse struct {
	User  *User  `json:"user"`
	Token string `json:"token"`
}

// RegisterRequest represents user registration payload
type RegisterRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=6"`
	Name     string `json:"name" validate:"required,min=2"`
	Phone    string `json:"phone" validate:"required"`
}

// UserRepository defines the interface for user data operations
type UserRepository interface {
	Create(ctx context.Context, user *User) (*User, error)
	GetByID(ctx context.Context, id uuid.UUID) (*User, error)
	GetByEmail(ctx context.Context, email string) (*User, error)
	Update(ctx context.Context, id uuid.UUID, user *User) (*User, error)
	Delete(ctx context.Context, id uuid.UUID) error
	List(ctx context.Context, page, limit int) ([]*User, int64, error)
	ListByRole(ctx context.Context, role UserRole, page, limit int) ([]*User, int64, error)
}

// UserUsecase defines usecase methods for User
type UserUsecase interface {
	Create(ctx context.Context, user *User) (*User, error)
	GetByID(ctx context.Context, id uuid.UUID) (*User, error)
	Update(ctx context.Context, id uuid.UUID, user *User) (*User, error)
	Delete(ctx context.Context, id uuid.UUID) error
	List(ctx context.Context, page, limit int) (*utils.PaginatedResponse[User], error)
	ListByRole(ctx context.Context, role UserRole, page, limit int) (*utils.PaginatedResponse[User], error)
	Authenticate(ctx context.Context, email, password string) (*User, string, error) // Returns user, token, error
}
