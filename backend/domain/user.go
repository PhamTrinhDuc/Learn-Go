package domain

import (
	"time"

	"github.com/google/uuid"
)

// UserRole defines the allowed roles for a user
type UserRole string

const (
	RoleCustomer UserRole = "customer"
	RoleManager  UserRole = "manager"
)

// User represents the users table in the database
type User struct {
	ID                uuid.UUID  `json:"id" db:"id"`
	Name              string     `json:"name" db:"name"`
	Phone             string     `json:"phone" db:"phone"`
	Birthday          *time.Time `json:"birthday,omitempty" db:"birthday"`
	Address           *string    `json:"address,omitempty" db:"address"`
	Role              UserRole   `json:"role" db:"role"`
	LoyaltyPoints     int        `json:"loyalty_points" db:"loyalty_points"`
	PreferredBranchID *uuid.UUID `json:"preferred_branch_id,omitempty" db:"preferred_branch_id"`
	LastVisitAt       *time.Time `json:"last_visit_at,omitempty" db:"last_visit_at"`
	CreatedAt         time.Time  `json:"created_at" db:"created_at"`
	UpdatedAt         time.Time  `json:"updated_at" db:"updated_at"`
}

// UserRepo defines the interface for user data operations
type UserRepo interface {
	GetByID(id uuid.UUID) (*User, error)
	GetByPhone(phone string) (*User, error)
	Create(user *User) error
	Update(user *User) error
	Delete(id uuid.UUID) error
	List() ([]*User, error)
}
