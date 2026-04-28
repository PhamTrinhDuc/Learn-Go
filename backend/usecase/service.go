package usecase

import (
	"backend/domain"
	"backend/internal/utils"
	"context"
)

type ServiceUsecase struct {
	repo domain.ServiceRepository
}

func NewServiceUsecase(repo domain.ServiceRepository) *ServiceUsecase {
	return &ServiceUsecase{repo: repo}
}

// Create validates and creates a new service
func (u *ServiceUsecase) Create(ctx context.Context, service *domain.Service) (*domain.Service, error) {
	if service == nil {
		return nil, utils.ErrNilInput
	}

	if err := utils.ValidateRequired("name", service.Name); err != nil {
		return nil, err
	}
	if service.DurationMinutes <= 0 {
		return nil, utils.ValidatePositiveInt("duration_minutes", service.DurationMinutes)
	}
	if service.EstimatedDuration <= 0 {
		return nil, utils.ValidatePositiveInt("estimated_duration", service.EstimatedDuration)
	}

	return u.repo.Create(ctx, service)
}

// GetByID retrieves a service by ID
func (u *ServiceUsecase) GetByID(ctx context.Context, id string) (*domain.Service, error) {
	if err := utils.ValidateID(id); err != nil {
		return nil, err
	}
	return u.repo.GetByID(ctx, id)
}

// Update validates and updates a service
func (u *ServiceUsecase) Update(ctx context.Context, id string, service *domain.Service) (*domain.Service, error) {
	if err := utils.ValidateID(id); err != nil {
		return nil, err
	}
	if service == nil {
		return nil, utils.ErrNilInput
	}

	if err := utils.ValidateRequired("name", service.Name); err != nil {
		return nil, err
	}
	if service.DurationMinutes <= 0 {
		return nil, utils.ValidatePositiveInt("duration_minutes", service.DurationMinutes)
	}
	if service.EstimatedDuration <= 0 {
		return nil, utils.ValidatePositiveInt("estimated_duration", service.EstimatedDuration)
	}

	// Get existing service to preserve ID and timestamps
	existing, err := u.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	existing.Name = service.Name
	existing.Category = service.Category
	existing.Description = service.Description
	existing.DurationMinutes = service.DurationMinutes
	existing.EstimatedDuration = service.EstimatedDuration
	existing.IsActive = service.IsActive

	return u.repo.Update(ctx, id, existing)
}

// Delete deletes a service
func (u *ServiceUsecase) Delete(ctx context.Context, id string) error {
	if err := utils.ValidateID(id); err != nil {
		return err
	}
	return u.repo.Delete(ctx, id)
}

// List retrieves a list of services with pagination
func (u *ServiceUsecase) List(ctx context.Context, page, limit int) (*utils.PaginatedResponse[domain.Service], error) {
	page, limit = utils.NormalizePagination(page, limit)
	services, total, err := u.repo.List(ctx, page, limit)
	if err != nil {
		return nil, err
	}
	return utils.CreatePaginatedResponse(services, total, page, limit), nil
}

// ListByCategory retrieves services by category with pagination
func (u *ServiceUsecase) ListByCategory(ctx context.Context, category string, page, limit int) (*utils.PaginatedResponse[domain.Service], error) {
	if err := utils.ValidateRequired("category", category); err != nil {
		return nil, err
	}
	page, limit = utils.NormalizePagination(page, limit)

	services, total, err := u.repo.ListByCategory(ctx, domain.CategoryType(category), page, limit)
	if err != nil {
		return nil, err
	}
	return utils.CreatePaginatedResponse(services, total, page, limit), nil
}
