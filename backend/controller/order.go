package controller

import (
	"backend/domain"
	"backend/internal/utils"

	"github.com/gin-gonic/gin"
)

type OrderController struct {
	usecase domain.OrdersUsecase
}

func NewOrderController(usecase domain.OrdersUsecase) *OrderController {
	return &OrderController{usecase: usecase}
}

func (c *OrderController) Create(ctx *gin.Context) {
	var order domain.Orders
	if err := ctx.ShouldBindJSON(&order); err != nil {
		utils.RespondBadRequest(ctx, err.Error())
		return
	}
	if err := c.usecase.Create(ctx.Request.Context(), &order); err != nil {
		utils.RespondBadRequest(ctx, err.Error())
		return
	}
	utils.RespondCreated(ctx, order)
}
func (c *OrderController) Get(ctx *gin.Context) {
	orders, err := c.usecase.Get(ctx.Request.Context())
	if err != nil {
		utils.RespondNotFound(ctx, err.Error())
		return
	}
	utils.RespondOK(ctx, orders)
}
func (c *OrderController) GetByID(ctx *gin.Context) {
	id := ctx.Param("id")
	order, err := c.usecase.GetByID(ctx.Request.Context(), id)
	if err != nil {
		utils.RespondNotFound(ctx, err.Error())
		return
	}
	utils.RespondOK(ctx, order)
}

func (c *OrderController) Update(ctx *gin.Context) {
	id := ctx.Param("id")
	var order domain.Orders
	if err := ctx.ShouldBindJSON(&order); err != nil {
		utils.RespondBadRequest(ctx, err.Error())
		return
	}
	if err := c.usecase.Update(ctx.Request.Context(), id, &order); err != nil {
		utils.RespondBadRequest(ctx, err.Error())
		return
	}
	utils.RespondOK(ctx, gin.H{"message": "order updated successfully"})
}

func (c *OrderController) Delete(ctx *gin.Context) {
	id := ctx.Param("id")
	if err := c.usecase.Delete(ctx.Request.Context(), id); err != nil {
		utils.RespondBadRequest(ctx, err.Error())
		return
	}
	utils.RespondNoContent(ctx)
}

func (c *OrderController) ListByUserID(ctx *gin.Context) {
	userId := ctx.Param("user_id")
	orders, err := c.usecase.ListByUserID(ctx.Request.Context(), userId)
	if err != nil {
		utils.RespondBadRequest(ctx, err.Error())
		return
	}
	utils.RespondOK(ctx, orders)
}

func (c *OrderController) ListByBranchID(ctx *gin.Context) {
	branchId := ctx.Param("branch_id")
	orders, err := c.usecase.ListByBranchID(ctx.Request.Context(), branchId)
	if err != nil {
		utils.RespondBadRequest(ctx, err.Error())
		return
	}
	utils.RespondOK(ctx, orders)
}

type OrderItemsController struct {
	usecase domain.OrderItemsUsecase
}

func NewOrderItemsController(usecase domain.OrderItemsUsecase) *OrderItemsController {
	return &OrderItemsController{usecase: usecase}
}

func (c *OrderItemsController) Create(ctx *gin.Context) {
	var item domain.OrderItems
	if err := ctx.ShouldBindJSON(&item); err != nil {
		utils.RespondBadRequest(ctx, err.Error())
		return
	}
	if err := c.usecase.Create(ctx.Request.Context(), &item); err != nil {
		utils.RespondBadRequest(ctx, err.Error())
		return
	}
	utils.RespondCreated(ctx, item)
}

func (c *OrderItemsController) GetByID(ctx *gin.Context) {
	id := ctx.Param("id")
	item, err := c.usecase.GetByID(ctx.Request.Context(), id)
	if err != nil {
		utils.RespondNotFound(ctx, err.Error())
		return
	}
	utils.RespondOK(ctx, item)
}

func (c *OrderItemsController) ListByOrderID(ctx *gin.Context) {
	orderId := ctx.Param("order_id")
	items, err := c.usecase.ListByOrderID(ctx.Request.Context(), orderId)
	if err != nil {
		utils.RespondBadRequest(ctx, err.Error())
		return
	}
	utils.RespondOK(ctx, items)
}

func (c *OrderItemsController) Update(ctx *gin.Context) {
	id := ctx.Param("id")
	var item domain.OrderItems
	if err := ctx.ShouldBindJSON(&item); err != nil {
		utils.RespondBadRequest(ctx, err.Error())
		return
	}
	if err := c.usecase.Update(ctx.Request.Context(), id, &item); err != nil {
		utils.RespondBadRequest(ctx, err.Error())
		return
	}
	utils.RespondOK(ctx, gin.H{"message": "order item updated successfully"})
}

func (c *OrderItemsController) Delete(ctx *gin.Context) {
	id := ctx.Param("id")
	if err := c.usecase.Delete(ctx.Request.Context(), id); err != nil {
		utils.RespondBadRequest(ctx, err.Error())
		return
	}
	utils.RespondNoContent(ctx)
}
