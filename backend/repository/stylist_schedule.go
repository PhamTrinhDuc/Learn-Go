package repository

import (
	"backend/domain"
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

type StylistScheduleRepository struct {
	db *pgxpool.Pool
}

func NewStylistScheduleRepository(db *pgxpool.Pool) *StylistScheduleRepository {
	return &StylistScheduleRepository{db: db}
}

// Create inserts a new stylist schedule
func (r *StylistScheduleRepository) Create(ctx context.Context, schedule *domain.StylistSchedule) (*domain.StylistSchedule, error) {
	if schedule.ID == "" {
		schedule.ID = uuid.New().String()
	}

	query := `
		INSERT INTO stylist_schedule (id, stylist_id, day_of_week, start_time, end_time, is_active, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, NOW(), NOW())
		RETURNING id, stylist_id, day_of_week, start_time, end_time, is_active, created_at, updated_at
	`

	err := r.db.QueryRow(ctx, query,
		schedule.ID,
		schedule.StylistID,
		schedule.DayOfWeek,
		schedule.StartTime,
		schedule.EndTime,
		schedule.IsActive,
	).Scan(
		&schedule.ID,
		&schedule.StylistID,
		&schedule.DayOfWeek,
		&schedule.StartTime,
		&schedule.EndTime,
		&schedule.IsActive,
		&schedule.CreatedAt,
		&schedule.UpdatedAt,
	)

	if err != nil {
		return nil, fmt.Errorf("failed to create stylist schedule: %w", err)
	}

	return schedule, nil
}

// GetByID retrieves a stylist schedule by ID
func (r *StylistScheduleRepository) GetByID(ctx context.Context, id string) (*domain.StylistSchedule, error) {
	query := `
		SELECT id, stylist_id, day_of_week, start_time, end_time, is_active, created_at, updated_at
		FROM stylist_schedule
		WHERE id = $1
	`

	schedule := &domain.StylistSchedule{}
	err := r.db.QueryRow(ctx, query, id).Scan(
		&schedule.ID,
		&schedule.StylistID,
		&schedule.DayOfWeek,
		&schedule.StartTime,
		&schedule.EndTime,
		&schedule.IsActive,
		&schedule.CreatedAt,
		&schedule.UpdatedAt,
	)

	if err != nil {
		return nil, fmt.Errorf("failed to get stylist schedule: %w", err)
	}

	return schedule, nil
}

// Update updates an existing stylist schedule
func (r *StylistScheduleRepository) Update(ctx context.Context, schedule *domain.StylistSchedule) (*domain.StylistSchedule, error) {
	query := `
		UPDATE stylist_schedule
		SET day_of_week = $1, start_time = $2, end_time = $3, is_active = $4, updated_at = NOW()
		WHERE id = $5
		RETURNING id, stylist_id, day_of_week, start_time, end_time, is_active, created_at, updated_at
	`

	err := r.db.QueryRow(ctx, query,
		schedule.DayOfWeek,
		schedule.StartTime,
		schedule.EndTime,
		schedule.IsActive,
		schedule.ID,
	).Scan(
		&schedule.ID,
		&schedule.StylistID,
		&schedule.DayOfWeek,
		&schedule.StartTime,
		&schedule.EndTime,
		&schedule.IsActive,
		&schedule.CreatedAt,
		&schedule.UpdatedAt,
	)

	if err != nil {
		return nil, fmt.Errorf("failed to update stylist schedule: %w", err)
	}

	return schedule, nil
}

// Delete removes a stylist schedule by ID
func (r *StylistScheduleRepository) Delete(ctx context.Context, id string) error {
	query := `DELETE FROM stylist_schedule WHERE id = $1`

	result, err := r.db.Exec(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to delete stylist schedule: %w", err)
	}

	if result.RowsAffected() == 0 {
		return fmt.Errorf("stylist schedule not found")
	}

	return nil
}

// ListByStylistID retrieves all schedules for a specific stylist
func (r *StylistScheduleRepository) ListByStylistID(ctx context.Context, stylistID string) ([]*domain.StylistSchedule, error) {
	query := `
		SELECT id, stylist_id, day_of_week, start_time, end_time, is_active, created_at, updated_at
		FROM stylist_schedule
		WHERE stylist_id = $1
		ORDER BY day_of_week ASC, start_time ASC
	`

	rows, err := r.db.Query(ctx, query, stylistID)
	if err != nil {
		return nil, fmt.Errorf("failed to list stylist schedules: %w", err)
	}
	defer rows.Close()

	schedules := make([]*domain.StylistSchedule, 0)
	for rows.Next() {
		s := &domain.StylistSchedule{}
		err := rows.Scan(
			&s.ID,
			&s.StylistID,
			&s.DayOfWeek,
			&s.StartTime,
			&s.EndTime,
			&s.IsActive,
			&s.CreatedAt,
			&s.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan stylist schedule: %w", err)
		}
		schedules = append(schedules, s)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating stylist schedules: %w", err)
	}

	return schedules, nil
}
