package usecase

import (
	"backend/domain"
	"backend/internal/utils"
	"context"

	"github.com/google/uuid"
)

type ProductUseCase struct {
	repo domain.ProductRepo
}

func NewProductUseCase(repo domain.ProductRepo) *ProductUseCase {
	return &ProductUseCase{repo: repo}
}

// Create validates and creates a product
func (u *ProductUseCase) Create(ctx context.Context, product *domain.Product) (*domain.Product, error) {
	if product == nil {
		return nil, utils.ErrNilInput
	}

	if err := utils.CombineErrors(
		utils.ValidateRequired("name", product.Name),
		utils.ValidateRequired("category", product.Category),
		utils.ValidatePositiveInt("price_out", int(product.PriceOut)),
	); err != nil {
		return nil, err
	}

	if product.ID == uuid.Nil {
		product.ID = uuid.New()
	}

	err := u.repo.Create(ctx, product)
	if err != nil {
		return nil, err
	}
	return product, nil
}

// GetByID retrieves a product by ID
func (u *ProductUseCase) GetByID(ctx context.Context, id uuid.UUID) (*domain.Product, error) {
	return u.repo.GetByID(ctx, id)
}

// GetProducts retrieves products with pagination and filters
func (u *ProductUseCase) List(ctx context.Context, search string, category string, minPrice *float64, maxPrice *float64, usageType *domain.UsageType, page, limit int) (*utils.PaginatedResponse[domain.Product], error) {
	page, limit = utils.NormalizePagination(page, limit)
	products, total, err := u.repo.List(ctx, search, category, minPrice, maxPrice, usageType, page, limit)
	if err != nil {
		return nil, err
	}
	return utils.CreatePaginatedResponse(products, total, page, limit), nil
}

// Update validates and updates a product
func (u *ProductUseCase) Update(ctx context.Context, id uuid.UUID, product *domain.Product) (*domain.Product, error) {
	if product == nil {
		return nil, utils.ErrNilInput
	}

	if err := utils.CombineErrors(
		utils.ValidateRequired("name", product.Name),
		utils.ValidateRequired("category", product.Category),
		utils.ValidatePositiveInt("price_out", int(product.PriceOut)),
	); err != nil {
		return nil, err
	}

	existing, err := u.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	existing.Name = product.Name
	existing.Description = product.Description
	existing.PriceIn = product.PriceIn
	existing.PriceOut = product.PriceOut
	existing.Category = product.Category
	existing.UsageType = product.UsageType
	existing.LowStockThresholdRetail = product.LowStockThresholdRetail
	existing.LowStockThresholdInternal = product.LowStockThresholdInternal

	err = u.repo.Update(ctx, existing)
	if err != nil {
		return nil, err
	}
	return existing, nil
}

// Delete deletes a product
func (u *ProductUseCase) Delete(ctx context.Context, id uuid.UUID) error {
	return u.repo.Delete(ctx, id)
}
