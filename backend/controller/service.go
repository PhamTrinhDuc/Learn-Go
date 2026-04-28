package controller

import (
	"backend/domain"
	"backend/internal/utils"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

type ServiceController struct {
	usecase  domain.ServiceUsecase
	validate *validator.Validate
}

func NewServiceController(usecase domain.ServiceUsecase) *ServiceController {
	return &ServiceController{
		usecase:  usecase,
		validate: validator.New(),
	}
}

// Create handles POST /services
func (c *ServiceController) Create(ctx *gin.Context) {
	var service domain.Service

	if err := ctx.ShouldBindJSON(&service); err != nil {
		utils.RespondBadRequest(ctx, err.Error())
		return
	}

	if err := c.validate.Struct(&service); err != nil {
		utils.RespondBadRequest(ctx, err.Error())
		return
	}

	result, err := c.usecase.Create(ctx.Request.Context(), &service)
	if err != nil {
		utils.RespondBadRequest(ctx, err.Error())
		return
	}

	utils.RespondCreated(ctx, result)
}

// GetByID handles GET /services/:id
func (c *ServiceController) GetByID(ctx *gin.Context) {
	id := ctx.Param("id")

	service, err := c.usecase.GetByID(ctx.Request.Context(), id)
	if err != nil {
		utils.RespondNotFound(ctx, err.Error())
		return
	}

	utils.RespondOK(ctx, service)
}

// Update handles PUT /services/:id
func (c *ServiceController) Update(ctx *gin.Context) {
	id := ctx.Param("id")

	var service domain.Service
	if err := ctx.ShouldBindJSON(&service); err != nil {
		utils.RespondBadRequest(ctx, err.Error())
		return
	}

	if err := c.validate.Struct(&service); err != nil {
		utils.RespondBadRequest(ctx, err.Error())
		return
	}

	result, err := c.usecase.Update(ctx.Request.Context(), id, &service)
	if err != nil {
		utils.RespondBadRequest(ctx, err.Error())
		return
	}

	utils.RespondOK(ctx, result)
}

// Delete handles DELETE /services/:id
func (c *ServiceController) Delete(ctx *gin.Context) {
	id := ctx.Param("id")

	if err := c.usecase.Delete(ctx.Request.Context(), id); err != nil {
		utils.RespondBadRequest(ctx, err.Error())
		return
	}

	utils.RespondNoContent(ctx)
}

// List handles GET /services
func (c *ServiceController) List(ctx *gin.Context) {
	page := utils.GetQueryInt(ctx, "page", 1)
	limit := utils.GetQueryInt(ctx, "limit", 10)

	result, err := c.usecase.List(ctx.Request.Context(), page, limit)
	if err != nil {
		utils.RespondError(ctx, err)
		return
	}

	utils.RespondOK(ctx, result)
}

// ListServicesByCategory handles GET /services/category/:category
func (c *ServiceController) ListServicesByCategory(ctx *gin.Context) {
	category := ctx.Param("category")
	page := utils.GetQueryInt(ctx, "page", 1)
	limit := utils.GetQueryInt(ctx, "limit", 10)

	result, err := c.usecase.ListByCategory(ctx.Request.Context(), category, page, limit)
	if err != nil {
		utils.RespondBadRequest(ctx, err.Error())
		return
	}

	utils.RespondOK(ctx, result)
}

// getQueryInt parses query parameter as int with default value
func getQueryInt(ctx *gin.Context, key string, defaultVal int) int {
	if val := ctx.Query(key); val != "" {
		if parsed, err := strconv.Atoi(val); err == nil && parsed > 0 {
			return parsed
		}
	}
	return defaultVal
}
