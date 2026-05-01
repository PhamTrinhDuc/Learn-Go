package domain

import (
	"context"
	"time"
)

type StylistSchedule struct {
	ID        string    `json:"id" db:"id"`
	StylistID string    `json:"stylist_id" db:"stylist_id"`
	DayOfWeek int       `json:"day_of_week" db:"day_of_week"` // 0 = Sunday, 6 = Saturday
	StartTime string    `json:"start_time" db:"start_time"`   // "09:00"
	EndTime   string    `json:"end_time" db:"end_time"`       // "18:00"
	IsActive  bool      `json:"is_active" db:"is_active"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
}

const (
	Sunday = iota
	Monday
	Tuesday
	Wednesday
	Thursday
	Friday
	Saturday
)

// Repository for database operations
type StylistScheduleRepo interface {
	GetByID(ctx context.Context, id string) (*StylistSchedule, error)
	Create(ctx context.Context, schedule *StylistSchedule) (*StylistSchedule, error)
	Update(ctx context.Context, schedule *StylistSchedule) (*StylistSchedule, error)
	Delete(ctx context.Context, id string) error
	ListByStylistID(ctx context.Context, stylistID string) ([]*StylistSchedule, error)
}

// Interface for service layer
type StylistScheduleUsecase interface {
	GetByID(ctx context.Context, id string) (*StylistSchedule, error)
	Create(ctx context.Context, schedule *StylistSchedule) (*StylistSchedule, error)
	Update(ctx context.Context, schedule *StylistSchedule) (*StylistSchedule, error)
	Delete(ctx context.Context, id string) error
	ListByStylistID(ctx context.Context, stylistID string) ([]*StylistSchedule, error)
}
