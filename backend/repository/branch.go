package repository

import (
	"backend/domain"
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

type BranchRepository struct {
	db *pgxpool.Pool
}

func NewBranchRepository(db *pgxpool.Pool) *BranchRepository {
	return &BranchRepository{db: db}
}

// Create inserts a new branch
func (r *BranchRepository) Create(ctx context.Context, branch *domain.Branch) (*domain.Branch, error) {
	if branch.ID == "" {
		branch.ID = uuid.New().String()
	}

	query := `
		INSERT INTO branch (id, name, address, phone, opening_hours, is_active, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, NOW(), NOW())
		RETURNING id, name, address, phone, opening_hours, is_active, created_at, updated_at
	`
	row := r.db.QueryRow(ctx, query,
		branch.ID,
		branch.Name,
		branch.Address,
		branch.Phone,
		branch.OpeningHours,
		true, // is_active defaults to true
	)

	err := row.Scan(
		&branch.ID,
		&branch.Name,
		&branch.Address,
		&branch.Phone,
		&branch.OpeningHours,
		&branch.IsActive,
		&branch.CreatedAt,
		&branch.UpdatedAt,
	)

	if err != nil {
		return nil, fmt.Errorf("failed to create branch: %w", err)
	}

	return branch, nil
}

// GetByID retrieves a branch by ID
func (r *BranchRepository) GetByID(ctx context.Context, id string) (*domain.Branch, error) {
	query := `
		SELECT id, name, address, phone, opening_hours, is_active, created_at, updated_at
		FROM branch
		WHERE id = $1
	`

	branch := &domain.Branch{}
	err := r.db.QueryRow(ctx, query, id).Scan(
		&branch.ID,
		&branch.Name,
		&branch.Address,
		&branch.Phone,
		&branch.OpeningHours,
		&branch.IsActive,
		&branch.CreatedAt,
		&branch.UpdatedAt,
	)

	if err != nil {
		return nil, fmt.Errorf("failed to get branch: %w", err)
	}

	return branch, nil
}

// Update updates a branch
func (r *BranchRepository) Update(ctx context.Context, id string, branch *domain.Branch) (*domain.Branch, error) {
	query := `
		UPDATE branch
		SET name = $1, address = $2, phone = $3, opening_hours = $4, is_active = $5, updated_at = NOW()
		WHERE id = $6
		RETURNING id, name, address, phone, opening_hours, is_active, created_at, updated_at
	`

	row := r.db.QueryRow(ctx, query,
		branch.Name,
		branch.Address,
		branch.Phone,
		branch.OpeningHours,
		branch.IsActive,
		id,
	)

	err := row.Scan(
		&branch.ID,
		&branch.Name,
		&branch.Address,
		&branch.Phone,
		&branch.OpeningHours,
		&branch.IsActive,
		&branch.CreatedAt,
		&branch.UpdatedAt,
	)

	if err != nil {
		return nil, fmt.Errorf("failed to update branch: %w", err)
	}

	return branch, nil
}

// Delete deletes a branch
func (r *BranchRepository) Delete(ctx context.Context, id string) error {
	query := `DELETE FROM branch WHERE id = $1`

	result, err := r.db.Exec(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to delete branch: %w", err)
	}

	if result.RowsAffected() == 0 {
		return fmt.Errorf("branch not found")
	}

	return nil
}

// List retrieves a list of branches with pagination
func (r *BranchRepository) List(ctx context.Context, page, limit int) ([]*domain.Branch, int64, error) {
	offset := (page - 1) * limit

	// Get total count
	var total int64
	countQuery := `SELECT COUNT(*) FROM branch`
	err := r.db.QueryRow(ctx, countQuery).Scan(&total)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count branches: %w", err)
	}

	// Get branches
	query := `
		SELECT id, name, address, phone, opening_hours, is_active, created_at, updated_at
		FROM branch
		ORDER BY created_at DESC
		LIMIT $1 OFFSET $2
	`

	rows, err := r.db.Query(ctx, query, limit, offset)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to list branches: %w", err)
	}
	defer rows.Close()

	branches := make([]*domain.Branch, 0)
	for rows.Next() {
		branch := &domain.Branch{}
		err := rows.Scan(
			&branch.ID,
			&branch.Name,
			&branch.Address,
			&branch.Phone,
			&branch.OpeningHours,
			&branch.IsActive,
			&branch.CreatedAt,
			&branch.UpdatedAt,
		)
		if err != nil {
			return nil, 0, fmt.Errorf("failed to scan branch: %w", err)
		}
		branches = append(branches, branch)
	}

	if err = rows.Err(); err != nil {
		return nil, 0, fmt.Errorf("error iterating branches: %w", err)
	}

	return branches, total, nil
}
