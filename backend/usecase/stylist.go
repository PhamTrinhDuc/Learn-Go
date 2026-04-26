package usecase

import (
	"backend/domain"
	"backend/internal/utils"
	"context"
	"fmt"
)

type StylistUsecase struct {
	repo domain.StylistRepository
}

func NewStylistUsecase(repo domain.StylistRepository) *StylistUsecase {
	return &StylistUsecase{repo: repo}
}

// Create creates a new stylist
func (u *StylistUsecase) Create(ctx context.Context, stylist *domain.Stylist) (*domain.Stylist, error) {
	if stylist == nil {
		return nil, fmt.Errorf("stylist cannot be nil")
	}

	if stylist.BranchID == "" || stylist.Name == "" || stylist.Phone == "" {
		return nil, fmt.Errorf("all fields are required")
	}

	return u.repo.Create(ctx, stylist)
}

// GetByID retrieves a stylist by ID
func (u *StylistUsecase) GetByID(ctx context.Context, id string) (*domain.Stylist, error) {
	if id == "" {
		return nil, fmt.Errorf("stylist id cannot be empty")
	}

	stylist, err := u.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	return stylist, nil
}

// Update updates a stylist
func (u *StylistUsecase) Update(ctx context.Context, id string, stylist *domain.Stylist) (*domain.Stylist, error) {
	if id == "" {
		return nil, fmt.Errorf("stylist id cannot be empty")
	}

	if stylist == nil {
		return nil, fmt.Errorf("stylist cannot be nil")
	}

	if stylist.BranchID == "" || stylist.Name == "" || stylist.Phone == "" {
		return nil, fmt.Errorf("all fields are required")
	}

	// Get existing stylist to preserve ID and timestamps
	existing, err := u.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	existing.BranchID = stylist.BranchID
	existing.Name = stylist.Name
	existing.Phone = stylist.Phone
	existing.IsActive = stylist.IsActive

	return u.repo.Update(ctx, id, existing)
}

// Delete deletes a stylist
func (u *StylistUsecase) Delete(ctx context.Context, id string) error {
	if id == "" {
		return fmt.Errorf("stylist id cannot be empty")
	}

	return u.repo.Delete(ctx, id)
}

// List retrieves a list of stylists with pagination
func (u *StylistUsecase) List(ctx context.Context, page, limit int) (*utils.PaginatedResponse[domain.Stylist], error) {
	page, limit = utils.NormalizePagination(page, limit)
	stylists, total, err := u.repo.List(ctx, page, limit)
	if err != nil {
		return nil, err
	}

	return utils.CreatePaginatedResponse(stylists, total, page, limit), nil
}

// ListByBranch retrieves stylists by branch with pagination
func (u *StylistUsecase) ListByBranch(ctx context.Context, branchID string, page, limit int) (*utils.PaginatedResponse[domain.Stylist], error) {
	if branchID == "" {
		return nil, fmt.Errorf("branch id cannot be empty")
	}

	page, limit = utils.NormalizePagination(page, limit)
	stylists, total, err := u.repo.ListByBranch(ctx, branchID, page, limit)
	if err != nil {
		return nil, err
	}

	return utils.CreatePaginatedResponse(stylists, total, page, limit), nil
}
