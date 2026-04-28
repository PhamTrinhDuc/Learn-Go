package domain

import (
	"backend/internal/utils"
	"context"
	"time"
)

type Branch struct {
	ID           string    `json:"id" db:"id"`
	Name         string    `json:"name" db:"name"`
	Address      string    `json:"address" db:"address"`
	Phone        string    `json:"phone" db:"phone"`
	OpeningHours string    `json:"opening_hours" db:"opening_hours"`
	IsActive     bool      `json:"is_active" db:"is_active"`
	CreatedAt    time.Time `json:"created_at" db:"created_at"`
	UpdatedAt    time.Time `json:"updated_at" db:"updated_at"`
}

type BranchRepository interface {
	Create(ctx context.Context, branch *Branch) (*Branch, error)
	GetByID(ctx context.Context, id string) (*Branch, error)
	Update(ctx context.Context, id string, branch *Branch) (*Branch, error)
	Delete(ctx context.Context, id string) error
	List(ctx context.Context, page, limit int) ([]*Branch, int64, error)
}

type BranchUsecase interface {
	Create(ctx context.Context, branch *Branch) (*Branch, error)
	GetByID(ctx context.Context, id string) (*Branch, error)
	Update(ctx context.Context, id string, branch *Branch) (*Branch, error)
	Delete(ctx context.Context, id string) error
	List(ctx context.Context, page, limit int) (*utils.PaginatedResponse[Branch], error)
}
