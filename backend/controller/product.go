package controller

import (
	"backend/domain"
	"backend/internal/utils"
	"backend/usecase"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
)

type ProductController struct {
	useCase  *usecase.ProductUseCase
	validate *validator.Validate
}

func NewProductController(useCase *usecase.ProductUseCase) *ProductController {
	return &ProductController{
		useCase:  useCase,
		validate: validator.New(),
	}
}

// List handles GET /products
func (c *ProductController) List(ctx *gin.Context) {
	search := ctx.Query("search")
	category := ctx.Query("category")
	page := utils.GetQueryInt(ctx, "page", 1)
	limit := utils.GetQueryInt(ctx, "limit", 10)

	var minPrice, maxPrice *float64
	if mp := ctx.Query("min_price"); mp != "" {
		if parsed, err := strconv.ParseFloat(mp, 64); err == nil {
			minPrice = &parsed
		}
	}
	if mp := ctx.Query("max_price"); mp != "" {
		if parsed, err := strconv.ParseFloat(mp, 64); err == nil {
			maxPrice = &parsed
		}
	}

	var usageType *domain.UsageType
	if ut := ctx.Query("usage_type"); ut != "" {
		t := domain.UsageType(ut)
		usageType = &t
	}

	products, err := c.useCase.List(ctx.Request.Context(), search, category, minPrice, maxPrice, usageType, page, limit)
	if err != nil {
		utils.RespondError(ctx, err)
		return
	}
	utils.RespondOK(ctx, products)
}

// Create handles POST /products
func (c *ProductController) Create(ctx *gin.Context) {
	var product domain.Product
	if err := ctx.ShouldBindJSON(&product); err != nil {
		utils.RespondBadRequest(ctx, err.Error())
		return
	}

	if err := c.validate.Struct(&product); err != nil {
		utils.RespondBadRequest(ctx, err.Error())
		return
	}

	result, err := c.useCase.Create(ctx.Request.Context(), &product)
	if err != nil {
		utils.RespondBadRequest(ctx, err.Error())
		return
	}
	utils.RespondCreated(ctx, result)
}

// GetByID handles GET /products/:id
func (c *ProductController) GetByID(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		utils.RespondBadRequest(ctx, "invalid id format")
		return
	}

	product, err := c.useCase.GetByID(ctx.Request.Context(), id)
	if err != nil {
		utils.RespondNotFound(ctx, err.Error())
		return
	}
	utils.RespondOK(ctx, product)
}

// Update handles PUT /products/:id
func (c *ProductController) Update(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		utils.RespondBadRequest(ctx, "invalid id format")
		return
	}

	var product domain.Product
	if err := ctx.ShouldBindJSON(&product); err != nil {
		utils.RespondBadRequest(ctx, err.Error())
		return
	}

	if err := c.validate.Struct(&product); err != nil {
		utils.RespondBadRequest(ctx, err.Error())
		return
	}

	updatedProduct, err := c.useCase.Update(ctx.Request.Context(), id, &product)
	if err != nil {
		utils.RespondBadRequest(ctx, err.Error())
		return
	}
	utils.RespondOK(ctx, updatedProduct)
}

// Delete handles DELETE /products/:id
func (c *ProductController) Delete(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		utils.RespondBadRequest(ctx, "invalid id format")
		return
	}

	err = c.useCase.Delete(ctx.Request.Context(), id)
	if err != nil {
		utils.RespondError(ctx, err)
		return
	}
	utils.RespondNoContent(ctx)
}
