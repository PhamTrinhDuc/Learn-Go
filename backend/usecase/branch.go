package usecase

import (
	"backend/domain"
	"backend/internal/utils"
	"context"
)

type BranchUsecase struct {
	repo domain.BranchRepository
}

func NewBranchUsecase(repo domain.BranchRepository) *BranchUsecase {
	return &BranchUsecase{repo: repo}
}

// Create validates and creates a new branch
func (u *BranchUsecase) Create(ctx context.Context, branch *domain.Branch) (*domain.Branch, error) {
	if branch == nil {
		return nil, utils.ErrNilInput
	}

	if err := utils.CombineErrors(
		utils.ValidateRequired("name", branch.Name),
		utils.ValidateRequired("address", branch.Address),
		utils.ValidateRequired("phone", branch.Phone),
		utils.ValidateRequired("opening_hours", branch.OpeningHours),
	); err != nil {
		return nil, err
	}

	return u.repo.Create(ctx, branch)
}

// GetByID retrieves a branch by ID
func (u *BranchUsecase) GetByID(ctx context.Context, id string) (*domain.Branch, error) {
	if err := utils.ValidateID(id); err != nil {
		return nil, err
	}
	return u.repo.GetByID(ctx, id)
}

// Update validates and updates a branch
func (u *BranchUsecase) Update(ctx context.Context, id string, branch *domain.Branch) (*domain.Branch, error) {
	if err := utils.ValidateID(id); err != nil {
		return nil, err
	}
	if branch == nil {
		return nil, utils.ErrNilInput
	}

	if err := utils.CombineErrors(
		utils.ValidateRequired("name", branch.Name),
		utils.ValidateRequired("address", branch.Address),
		utils.ValidateRequired("phone", branch.Phone),
		utils.ValidateRequired("opening_hours", branch.OpeningHours),
	); err != nil {
		return nil, err
	}

	existing, err := u.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	existing.Name = branch.Name
	existing.Address = branch.Address
	existing.Phone = branch.Phone
	existing.OpeningHours = branch.OpeningHours
	existing.IsActive = branch.IsActive

	return u.repo.Update(ctx, id, existing)
}

// Delete deletes a branch
func (u *BranchUsecase) Delete(ctx context.Context, id string) error {
	if err := utils.ValidateID(id); err != nil {
		return err
	}
	return u.repo.Delete(ctx, id)
}

// List retrieves a list of branches with pagination
func (u *BranchUsecase) List(ctx context.Context, page, limit int) (*utils.PaginatedResponse[domain.Branch], error) {
	page, limit = utils.NormalizePagination(page, limit)
	branches, total, err := u.repo.List(ctx, page, limit)
	if err != nil {
		return nil, err
	}
	return utils.CreatePaginatedResponse(branches, total, page, limit), nil
}
