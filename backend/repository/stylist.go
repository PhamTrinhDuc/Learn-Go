package repository

import (
	"backend/domain"
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

type StylistRepository struct {
	db *pgxpool.Pool
}

func NewStylistRepository(db *pgxpool.Pool) *StylistRepository {
	return &StylistRepository{db: db}
}

// Create inserts a new stylist
func (r *StylistRepository) Create(ctx context.Context, stylist *domain.Stylist) (*domain.Stylist, error) {
	if stylist.ID == "" {
		stylist.ID = uuid.New().String()
	}

	query := `
		INSERT INTO stylist (id, branch_id, name, phone, is_active, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, NOW(), NOW())
		RETURNING id, branch_id, name, phone, is_active, created_at, updated_at
	`

	row := r.db.QueryRow(ctx, query,
		stylist.ID,
		stylist.BranchID,
		stylist.Name,
		stylist.Phone,
		true,
	)

	err := row.Scan(
		&stylist.ID,
		&stylist.BranchID,
		&stylist.Name,
		&stylist.Phone,
		&stylist.IsActive,
		&stylist.CreatedAt,
		&stylist.UpdatedAt,
	)

	if err != nil {
		return nil, fmt.Errorf("failed to create stylist: %w", err)
	}

	return stylist, nil
}

// GetByID retrieves a stylist by ID
func (r *StylistRepository) GetByID(ctx context.Context, id string) (*domain.Stylist, error) {
	query := `
		SELECT id, branch_id, name, phone, is_active, created_at, updated_at
		FROM stylist
		WHERE id = $1
	`

	stylist := &domain.Stylist{}
	err := r.db.QueryRow(ctx, query, id).Scan(
		&stylist.ID,
		&stylist.BranchID,
		&stylist.Name,
		&stylist.Phone,
		&stylist.IsActive,
		&stylist.CreatedAt,
		&stylist.UpdatedAt,
	)

	if err != nil {
		return nil, fmt.Errorf("failed to get stylist: %w", err)
	}

	return stylist, nil
}

// Update updates a stylist
func (r *StylistRepository) Update(ctx context.Context, id string, stylist *domain.Stylist) (*domain.Stylist, error) {
	query := `
		UPDATE stylist
		SET branch_id = $1, name = $2, phone = $3, is_active = $4, updated_at = NOW()
		WHERE id = $5
		RETURNING id, branch_id, name, phone, is_active, created_at, updated_at
	`

	row := r.db.QueryRow(ctx, query,
		stylist.BranchID,
		stylist.Name,
		stylist.Phone,
		stylist.IsActive,
		id,
	)

	err := row.Scan(
		&stylist.ID,
		&stylist.BranchID,
		&stylist.Name,
		&stylist.Phone,
		&stylist.IsActive,
		&stylist.CreatedAt,
		&stylist.UpdatedAt,
	)

	if err != nil {
		return nil, fmt.Errorf("failed to update stylist: %w", err)
	}

	return stylist, nil
}

// Delete deletes a stylist
func (r *StylistRepository) Delete(ctx context.Context, id string) error {
	query := `DELETE FROM stylist WHERE id = $1`

	result, err := r.db.Exec(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to delete stylist: %w", err)
	}

	if result.RowsAffected() == 0 {
		return fmt.Errorf("stylist not found")
	}

	return nil
}

// List retrieves a list of stylists with pagination
func (r *StylistRepository) List(ctx context.Context, page, limit int) ([]*domain.Stylist, int64, error) {
	offset := (page - 1) * limit

	// Get total count
	var total int64
	countQuery := `SELECT COUNT(*) FROM stylist`
	err := r.db.QueryRow(ctx, countQuery).Scan(&total)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count stylists: %w", err)
	}

	// Get stylists
	query := `
		SELECT id, branch_id, name, phone, is_active, created_at, updated_at
		FROM stylist
		ORDER BY created_at DESC
		LIMIT $1 OFFSET $2
	`

	rows, err := r.db.Query(ctx, query, limit, offset)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to list stylists: %w", err)
	}
	defer rows.Close()

	stylists := make([]*domain.Stylist, 0)
	for rows.Next() {
		stylist := &domain.Stylist{}
		err := rows.Scan(
			&stylist.ID,
			&stylist.BranchID,
			&stylist.Name,
			&stylist.Phone,
			&stylist.IsActive,
			&stylist.CreatedAt,
			&stylist.UpdatedAt,
		)
		if err != nil {
			return nil, 0, fmt.Errorf("failed to scan stylist: %w", err)
		}
		stylists = append(stylists, stylist)
	}

	if err = rows.Err(); err != nil {
		return nil, 0, fmt.Errorf("error iterating stylists: %w", err)
	}

	return stylists, total, nil
}

// ListByBranch retrieves stylists by branch with pagination
func (r *StylistRepository) ListByBranch(ctx context.Context, branchID string, page, limit int) ([]*domain.Stylist, int64, error) {
	offset := (page - 1) * limit

	// Get total count
	var total int64
	countQuery := `SELECT COUNT(*) FROM stylist WHERE branch_id = $1`
	err := r.db.QueryRow(ctx, countQuery, branchID).Scan(&total)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count stylists: %w", err)
	}

	// Get stylists
	query := `
		SELECT id, branch_id, name, phone, is_active, created_at, updated_at
		FROM stylist
		WHERE branch_id = $1
		ORDER BY created_at DESC
		LIMIT $2 OFFSET $3
	`

	rows, err := r.db.Query(ctx, query, branchID, limit, offset)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to list stylists by branch: %w", err)
	}
	defer rows.Close()

	stylists := make([]*domain.Stylist, 0)
	for rows.Next() {
		stylist := &domain.Stylist{}
		err := rows.Scan(
			&stylist.ID,
			&stylist.BranchID,
			&stylist.Name,
			&stylist.Phone,
			&stylist.IsActive,
			&stylist.CreatedAt,
			&stylist.UpdatedAt,
		)
		if err != nil {
			return nil, 0, fmt.Errorf("failed to scan stylist: %w", err)
		}
		stylists = append(stylists, stylist)
	}

	if err = rows.Err(); err != nil {
		return nil, 0, fmt.Errorf("error iterating stylists: %w", err)
	}

	return stylists, total, nil
}
