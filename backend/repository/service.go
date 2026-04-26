package repository

import (
	"backend/domain"
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

type ServiceRepository struct {
	db *pgxpool.Pool
}

func NewServiceRepository(db *pgxpool.Pool) *ServiceRepository {
	return &ServiceRepository{db: db}
}

// Create inserts a new service
func (r *ServiceRepository) Create(ctx context.Context, service *domain.Service) (*domain.Service, error) {
	if service.ID == "" {
		service.ID = uuid.New().String()
	}

	query := `
		INSERT INTO service (id, name, category, description, duration_minutes, estimated_duration, is_active, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, NOW(), NOW())
		RETURNING id, name, category, description, duration_minutes, estimated_duration, is_active, created_at, updated_at
	`

	row := r.db.QueryRow(ctx, query,
		service.ID,
		service.Name,
		string(service.Category),
		service.Description,
		service.DurationMinutes,
		service.EstimatedDuration,
		true,
	)

	err := row.Scan(
		&service.ID,
		&service.Name,
		(*string)(&service.Category),
		&service.Description,
		&service.DurationMinutes,
		&service.EstimatedDuration,
		&service.IsActive,
		&service.CreatedAt,
		&service.UpdatedAt,
	)

	if err != nil {
		return nil, fmt.Errorf("failed to create service: %w", err)
	}

	return service, nil
}

// GetByID retrieves a service by ID
func (r *ServiceRepository) GetByID(ctx context.Context, id string) (*domain.Service, error) {
	query := `
		SELECT id, name, category, description, duration_minutes, estimated_duration, is_active, created_at, updated_at
		FROM service
		WHERE id = $1
	`

	service := &domain.Service{}
	err := r.db.QueryRow(ctx, query, id).Scan(
		&service.ID,
		&service.Name,
		(*string)(&service.Category),
		&service.Description,
		&service.DurationMinutes,
		&service.EstimatedDuration,
		&service.IsActive,
		&service.CreatedAt,
		&service.UpdatedAt,
	)

	if err != nil {
		return nil, fmt.Errorf("failed to get service: %w", err)
	}

	return service, nil
}

// Update updates a service
func (r *ServiceRepository) Update(ctx context.Context, id string, service *domain.Service) (*domain.Service, error) {
	query := `
		UPDATE service
		SET name = $1, category = $2, description = $3, duration_minutes = $4, estimated_duration = $5, is_active = $6, updated_at = NOW()
		WHERE id = $7
		RETURNING id, name, category, description, duration_minutes, estimated_duration, is_active, created_at, updated_at
	`

	row := r.db.QueryRow(ctx, query,
		service.Name,
		string(service.Category),
		service.Description,
		service.DurationMinutes,
		service.EstimatedDuration,
		service.IsActive,
		id,
	)

	err := row.Scan(
		&service.ID,
		&service.Name,
		(*string)(&service.Category),
		&service.Description,
		&service.DurationMinutes,
		&service.EstimatedDuration,
		&service.IsActive,
		&service.CreatedAt,
		&service.UpdatedAt,
	)

	if err != nil {
		return nil, fmt.Errorf("failed to update service: %w", err)
	}

	return service, nil
}

// Delete deletes a service
func (r *ServiceRepository) Delete(ctx context.Context, id string) error {
	query := `DELETE FROM service WHERE id = $1`

	result, err := r.db.Exec(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to delete service: %w", err)
	}

	if result.RowsAffected() == 0 {
		return fmt.Errorf("service not found")
	}

	return nil
}

// List retrieves a list of services with pagination
func (r *ServiceRepository) List(ctx context.Context, page, limit int) ([]*domain.Service, int64, error) {

	offset := (page - 1) * limit

	// Get total count
	var total int64
	countQuery := `SELECT COUNT(*) FROM service`
	err := r.db.QueryRow(ctx, countQuery).Scan(&total)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count services: %w", err)
	}

	// Get services
	query := `
		SELECT id, name, category, description, duration_minutes, estimated_duration, is_active, created_at, updated_at
		FROM service
		ORDER BY created_at DESC
		LIMIT $1 OFFSET $2
	`

	rows, err := r.db.Query(ctx, query, limit, offset)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to list services: %w", err)
	}
	defer rows.Close()

	services := make([]*domain.Service, 0)
	for rows.Next() {
		service := &domain.Service{}
		err := rows.Scan(
			&service.ID,
			&service.Name,
			(*string)(&service.Category),
			&service.Description,
			&service.DurationMinutes,
			&service.EstimatedDuration,
			&service.IsActive,
			&service.CreatedAt,
			&service.UpdatedAt,
		)
		if err != nil {
			return nil, 0, fmt.Errorf("failed to scan service: %w", err)
		}
		services = append(services, service)
	}

	if err = rows.Err(); err != nil {
		return nil, 0, fmt.Errorf("error iterating services: %w", err)
	}

	return services, total, nil
}

// ListByCategory retrieves services by category with pagination
func (r *ServiceRepository) ListByCategory(ctx context.Context, category domain.CategoryType, page, limit int) ([]*domain.Service, int64, error) {

	offset := (page - 1) * limit

	// Get total count
	var total int64
	countQuery := `SELECT COUNT(*) FROM service WHERE category = $1`
	err := r.db.QueryRow(ctx, countQuery, string(category)).Scan(&total)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count services: %w", err)
	}

	// Get services
	query := `
		SELECT id, name, category, description, duration_minutes, estimated_duration, is_active, created_at, updated_at
		FROM service
		WHERE category = $1
		ORDER BY created_at DESC
		LIMIT $2 OFFSET $3
	`

	rows, err := r.db.Query(ctx, query, string(category), limit, offset)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to list services by category: %w", err)
	}
	defer rows.Close()

	services := make([]*domain.Service, 0)
	for rows.Next() {
		service := &domain.Service{}
		err := rows.Scan(
			&service.ID,
			&service.Name,
			(*string)(&service.Category),
			&service.Description,
			&service.DurationMinutes,
			&service.EstimatedDuration,
			&service.IsActive,
			&service.CreatedAt,
			&service.UpdatedAt,
		)
		if err != nil {
			return nil, 0, fmt.Errorf("failed to scan service: %w", err)
		}
		services = append(services, service)
	}

	if err = rows.Err(); err != nil {
		return nil, 0, fmt.Errorf("error iterating services: %w", err)
	}

	return services, total, nil
}
