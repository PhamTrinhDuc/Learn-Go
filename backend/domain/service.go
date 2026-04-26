package domain

import (
	"backend/internal/utils"
	"context"
)

type CategoryType string

const (
	CategoryCut       CategoryType = "cut"
	CategoryColor     CategoryType = "color"
	CategoryTreatment CategoryType = "treatment"
)

type Service struct {
	ID                string       `json:"id" db:"id"`
	Name              string       `json:"name" db:"name"`
	Category          CategoryType `json:"category" db:"category"`
	Description       *string      `json:"description,omitempty" db:"description"`
	DurationMinutes   int          `json:"duration_minutes" db:"duration_minutes"`
	EstimatedDuration int          `json:"estimated_duration" db:"estimated_duration"`
	IsActive          bool         `json:"is_active" db:"is_active"`
	CreatedAt         string       `json:"created_at" db:"created_at"`
	UpdatedAt         string       `json:"updated_at" db:"updated_at"`
}

type ServiceRepository interface {
	Create(ctx context.Context, service *Service) (*Service, error)
	GetByID(ctx context.Context, id string) (*Service, error)
	Update(ctx context.Context, id string, service *Service) (*Service, error)
	Delete(ctx context.Context, id string) error
	List(ctx context.Context, page, limit int) ([]*Service, int64, error)
	ListByCategory(ctx context.Context, category CategoryType, page, limit int) ([]*Service, int64, error)
}

type ServiceUsecase interface {
	Create(ctx context.Context, service *Service) (*Service, error)
	GetByID(ctx context.Context, id string) (*Service, error)
	Update(ctx context.Context, id string, service *Service) (*Service, error)
	Delete(ctx context.Context, id string) error
	List(ctx context.Context, page, limit int) (*utils.PaginatedResponse[Service], error)
	ListByCategory(ctx context.Context, category string, page, limit int) (*utils.PaginatedResponse[Service], error)
}
