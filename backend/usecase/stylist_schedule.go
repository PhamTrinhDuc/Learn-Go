package usecase

import (
	"backend/domain"
	"context"
	"fmt"
)

type StylistScheduleUsecase struct {
	repo domain.StylistScheduleRepo
}

func NewStylistScheduleUsecase(repo domain.StylistScheduleRepo) *StylistScheduleUsecase {
	return &StylistScheduleUsecase{repo: repo}
}

func (u *StylistScheduleUsecase) GetByID(ctx context.Context, id string) (*domain.StylistSchedule, error) {
	if id == "" {
		return nil, fmt.Errorf("id cannot be empty")
	}
	return u.repo.GetByID(ctx, id)
}

func (u *StylistScheduleUsecase) Create(ctx context.Context, schedule *domain.StylistSchedule) (*domain.StylistSchedule, error) {
	if schedule == nil {
		return nil, fmt.Errorf("schedule cannot be nil")
	}
	if schedule.StylistID == "" {
		return nil, fmt.Errorf("stylist id is required")
	}
	// Basic validation for time format "HH:MM"
	if len(schedule.StartTime) != 5 || len(schedule.EndTime) != 5 {
		return nil, fmt.Errorf("invalid time format, expected HH:MM")
	}
	return u.repo.Create(ctx, schedule)
}

func (u *StylistScheduleUsecase) Update(ctx context.Context, schedule *domain.StylistSchedule) (*domain.StylistSchedule, error) {
	if schedule == nil {
		return nil, fmt.Errorf("schedule cannot be nil")
	}
	if schedule.ID == "" {
		return nil, fmt.Errorf("id is required for update")
	}
	
	existing, err := u.repo.GetByID(ctx, schedule.ID)
	if err != nil {
		return nil, err
	}

	existing.DayOfWeek = schedule.DayOfWeek
	existing.StartTime = schedule.StartTime
	existing.EndTime = schedule.EndTime
	existing.IsActive = schedule.IsActive

	return u.repo.Update(ctx, existing)
}

func (u *StylistScheduleUsecase) Delete(ctx context.Context, id string) error {
	if id == "" {
		return fmt.Errorf("id cannot be empty")
	}
	return u.repo.Delete(ctx, id)
}

func (u *StylistScheduleUsecase) ListByStylistID(ctx context.Context, stylistID string) ([]*domain.StylistSchedule, error) {
	if stylistID == "" {
		return nil, fmt.Errorf("stylist id cannot be empty")
	}
	return u.repo.ListByStylistID(ctx, stylistID)
}
