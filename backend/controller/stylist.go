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

type StylistController struct {
	usecase   domain.StylistUsecase
	validate  *validator.Validate
	telemetry *observability.Telemetry
}

func NewStylistController(usecase domain.StylistUsecase, telemetry *observability.Telemetry) *StylistController {
	return &StylistController{
		usecase:   usecase,
		validate:  validator.New(),
		telemetry: telemetry,
	}
}

// Create handles POST /stylists
func (c *StylistController) Create(ctx *gin.Context) {
	var stylist domain.Stylist

	if err := ctx.ShouldBindJSON(&stylist); err != nil {
		utils.RespondBadRequest(ctx, err.Error())
		return
	}

	if err := c.validate.Struct(&stylist); err != nil {
		utils.RespondBadRequest(ctx, err.Error())
		return
	}

	result, err := c.usecase.Create(ctx.Request.Context(), &stylist)
	if err != nil {
		utils.RespondBadRequest(ctx, err.Error())
		return
	}

	utils.RespondCreated(ctx, result)
}

// GetByID handles GET /stylists/:id
func (c *StylistController) GetByID(ctx *gin.Context) {
	id := ctx.Param("id")

	stylist, err := c.usecase.GetByID(ctx.Request.Context(), id)
	if err != nil {
		utils.RespondNotFound(ctx, err.Error())
		return
	}

	utils.RespondOK(ctx, stylist)
}

// Update handles PUT /stylists/:id
func (c *StylistController) Update(ctx *gin.Context) {
	id := ctx.Param("id")

	var stylist domain.Stylist
	if err := ctx.ShouldBindJSON(&stylist); err != nil {
		utils.RespondBadRequest(ctx, err.Error())
		return
	}

	if err := c.validate.Struct(&stylist); err != nil {
		utils.RespondBadRequest(ctx, err.Error())
		return
	}

	result, err := c.usecase.Update(ctx.Request.Context(), id, &stylist)
	if err != nil {
		utils.RespondBadRequest(ctx, err.Error())
		return
	}

	utils.RespondOK(ctx, result)
}

// Delete handles DELETE /stylists/:id
func (c *StylistController) Delete(ctx *gin.Context) {
	id := ctx.Param("id")

	if err := c.usecase.Delete(ctx.Request.Context(), id); err != nil {
		utils.RespondBadRequest(ctx, err.Error())
		return
	}

	utils.RespondNoContent(ctx)
}

// List handles GET /stylists
func (c *StylistController) List(ctx *gin.Context) {
	page := utils.GetQueryInt(ctx, "page", 1)
	limit := utils.GetQueryInt(ctx, "limit", 10)

	reqCtx, span := c.telemetry.Tracer.Start(
		ctx.Request.Context(),
		"StylistController.List", // Tên Span nên rõ ràng
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

// ListStylistsByBranch handles GET /stylists/branch/:branch_id
func (c *StylistController) ListStylistsByBranch(ctx *gin.Context) {
	branchID := ctx.Param("branch_id")
	page := utils.GetQueryInt(ctx, "page", 1)
	limit := utils.GetQueryInt(ctx, "limit", 10)

	result, err := c.usecase.ListByBranch(ctx.Request.Context(), branchID, page, limit)
	if err != nil {
		utils.RespondBadRequest(ctx, err.Error())
		return
	}

	utils.RespondOK(ctx, result)
}
