package controller

import (
	"backend/domain"
	"backend/internal/observability"
	"backend/internal/utils"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
)

type BranchController struct {
	usecase   domain.BranchUsecase
	validate  *validator.Validate
	telemetry *observability.Telemetry
}

func NewBranchController(usecase domain.BranchUsecase, telemetry *observability.Telemetry) *BranchController {
	return &BranchController{
		usecase:   usecase,
		validate:  validator.New(),
		telemetry: telemetry,
	}
}

// Create handles POST /branches
func (c *BranchController) Create(ctx *gin.Context) {
	var branch domain.Branch

	if err := ctx.ShouldBindJSON(&branch); err != nil {
		utils.RespondBadRequest(ctx, err.Error())
		return
	}

	if err := c.validate.Struct(&branch); err != nil {
		utils.RespondBadRequest(ctx, err.Error())
		return
	}

	result, err := c.usecase.Create(ctx.Request.Context(), &branch)
	if err != nil {
		utils.RespondBadRequest(ctx, err.Error())
		return
	}

	utils.RespondCreated(ctx, result)
}

// GetByID handles GET /branches/:id
func (c *BranchController) GetByID(ctx *gin.Context) {
	id := ctx.Param("id")

	branch, err := c.usecase.GetByID(ctx.Request.Context(), id)
	if err != nil {
		utils.RespondNotFound(ctx, err.Error())
		return
	}

	utils.RespondOK(ctx, branch)
}

// Update handles PUT /branches/:id
func (c *BranchController) Update(ctx *gin.Context) {
	id := ctx.Param("id")

	var branch domain.Branch
	if err := ctx.ShouldBindJSON(&branch); err != nil {
		utils.RespondBadRequest(ctx, err.Error())
		return
	}

	if err := c.validate.Struct(&branch); err != nil {
		utils.RespondBadRequest(ctx, err.Error())
		return
	}

	result, err := c.usecase.Update(ctx.Request.Context(), id, &branch)
	if err != nil {
		utils.RespondBadRequest(ctx, err.Error())
		return
	}

	utils.RespondOK(ctx, result)
}

// Delete handles DELETE /branches/:id
func (c *BranchController) Delete(ctx *gin.Context) {
	id := ctx.Param("id")

	if err := c.usecase.Delete(ctx.Request.Context(), id); err != nil {
		utils.RespondBadRequest(ctx, err.Error())
		return
	}

	utils.RespondNoContent(ctx)
}

// List handles GET /branches
func (c *BranchController) List(ctx *gin.Context) {
	page := utils.GetQueryInt(ctx, "page", 1)
	limit := utils.GetQueryInt(ctx, "limit", 10)

	reqCtx, span := c.telemetry.Tracer.Start(
		ctx.Request.Context(),
		"BranchController.List", 
		trace.WithAttributes(
			attribute.String("http.method", ctx.Request.Method),
			attribute.String("http.path", ctx.Request.URL.Path),
			attribute.Int("pagination.page", page),
			attribute.Int("pagination.limit", limit),
		),
	)
	defer span.End()

	result, err := c.usecase.List(reqCtx, page, limit)

	if err != nil {
		utils.RespondError(ctx, err)
		// Record lỗi vào Span để debug trên Jaeger
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())

		utils.RespondError(ctx, err)
		return
	}

	utils.RespondOK(ctx, result)
}
