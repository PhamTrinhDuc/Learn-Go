package usecase

import (
	"backend/domain"
	"backend/internal/utils"
	"context"
)

type OrderUsecase struct {
	repo domain.OrdersRepo
}

func NewOrderUsecase(repo domain.OrdersRepo) *OrderUsecase {
	return &OrderUsecase{repo: repo}
}

func (u *OrderUsecase) GetByID(ctx context.Context, id string) (*domain.Orders, error) {
	if err := utils.ValidateID(id); err != nil {
		return nil, err
	}
	return u.repo.GetByID(ctx, id)
}

func (u *OrderUsecase) Create(ctx context.Context, order *domain.Orders) error {
	if order == nil {
		return utils.ErrNilInput
	}
	if err := utils.ValidateRequired("user_id", order.UserID); err != nil {
		return err
	}
	if err := utils.ValidateRequired("branch_id", order.BranchID); err != nil {
		return err
	}
	if order.TotalAmount < 0 {
		return utils.ValidatePositiveInt("total_amount", int(order.TotalAmount)) // simple validation
	}
	return u.repo.Create(ctx, order)
}

func (u *OrderUsecase) Update(ctx context.Context, id string, order *domain.Orders) error {
	if err := utils.ValidateID(id); err != nil {
		return err
	}
	if order == nil {
		return utils.ErrNilInput
	}

	existing, err := u.repo.GetByID(ctx, id)
	if err != nil {
		return err
	}

	existing.UserID = order.UserID
	existing.BranchID = order.BranchID
	existing.TotalAmount = order.TotalAmount
	existing.PaymentMethod = order.PaymentMethod
	existing.Status = order.Status

	return u.repo.Update(ctx, existing)
}

func (u *OrderUsecase) Delete(ctx context.Context, id string) error {
	if err := utils.ValidateID(id); err != nil {
		return err
	}
	return u.repo.Delete(ctx, id)
}

func (u *OrderUsecase) ListByUserID(ctx context.Context, userId string) ([]*domain.Orders, error) {
	if err := utils.ValidateID(userId); err != nil {
		return nil, err
	}
	return u.repo.ListByUserID(ctx, userId)
}

func (u *OrderUsecase) ListByBranchID(ctx context.Context, branchId string) ([]*domain.Orders, error) {
	if err := utils.ValidateID(branchId); err != nil {
		return nil, err
	}
	return u.repo.ListByBranchID(ctx, branchId)
}

// OrderItemsUsecase implementation

type OrderItemsUsecase struct {
	repo domain.OrderItemsRepo
}

func NewOrderItemsUsecase(repo domain.OrderItemsRepo) *OrderItemsUsecase {
	return &OrderItemsUsecase{repo: repo}
}

func (u *OrderItemsUsecase) GetByID(ctx context.Context, id string) (*domain.OrderItems, error) {
	if err := utils.ValidateID(id); err != nil {
		return nil, err
	}
	return u.repo.GetByID(ctx, id)
}

func (u *OrderUsecase) Get(ctx context.Context) ([]*domain.Orders, error) {
	return u.repo.Get(ctx)
}

func (u *OrderItemsUsecase) ListByOrderID(ctx context.Context, orderId string) ([]*domain.OrderItems, error) {
	if err := utils.ValidateID(orderId); err != nil {
		return nil, err
	}
	return u.repo.ListByOrderID(ctx, orderId)
}

func (u *OrderItemsUsecase) Create(ctx context.Context, orderItem *domain.OrderItems) error {
	if orderItem == nil {
		return utils.ErrNilInput
	}
	if err := utils.ValidateRequired("order_id", orderItem.OrderID); err != nil {
		return err
	}
	if err := utils.ValidateRequired("product_id", orderItem.ProductID); err != nil {
		return err
	}
	if orderItem.Quantity <= 0 {
		return utils.ValidatePositiveInt("quantity", orderItem.Quantity)
	}
	return u.repo.Create(ctx, orderItem)
}

func (u *OrderItemsUsecase) Update(ctx context.Context, id string, orderItem *domain.OrderItems) error {
	if err := utils.ValidateID(id); err != nil {
		return err
	}
	if orderItem == nil {
		return utils.ErrNilInput
	}

	existing, err := u.repo.GetByID(ctx, id)
	if err != nil {
		return err
	}

	existing.OrderID = orderItem.OrderID
	existing.ProductID = orderItem.ProductID
	existing.Quantity = orderItem.Quantity
	existing.UnitPrice = orderItem.UnitPrice

	return u.repo.Update(ctx, existing)
}

func (u *OrderItemsUsecase) Delete(ctx context.Context, id string) error {
	if err := utils.ValidateID(id); err != nil {
		return err
	}
	return u.repo.Delete(ctx, id)
}
