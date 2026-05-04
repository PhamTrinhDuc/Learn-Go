package controller

import (
	"backend/domain"
	"backend/internal/utils"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

type StylistScheduleController struct {
	usecase  domain.StylistScheduleUsecase
	validate *validator.Validate
}

func NewStylistScheduleController(usecase domain.StylistScheduleUsecase) *StylistScheduleController {
	return &StylistScheduleController{
		usecase:  usecase,
		validate: validator.New(),
	}
}

// Create handles POST /stylist-schedules
func (c *StylistScheduleController) Create(ctx *gin.Context) {
	var schedule domain.StylistSchedule

	if err := ctx.ShouldBindJSON(&schedule); err != nil {
		utils.RespondBadRequest(ctx, err.Error())
		return
	}

	result, err := c.usecase.Create(ctx.Request.Context(), &schedule)
	if err != nil {
		utils.RespondBadRequest(ctx, err.Error())
		return
	}

	utils.RespondCreated(ctx, result)
}

// GetByID handles GET /stylist-schedules/:id
func (c *StylistScheduleController) GetByID(ctx *gin.Context) {
	id := ctx.Param("id")

	schedule, err := c.usecase.GetByID(ctx.Request.Context(), id)
	if err != nil {
		utils.RespondNotFound(ctx, err.Error())
		return
	}

	utils.RespondOK(ctx, schedule)
}

// Update handles PUT /stylist-schedules/:id
func (c *StylistScheduleController) Update(ctx *gin.Context) {
	id := ctx.Param("id")

	var schedule domain.StylistSchedule
	if err := ctx.ShouldBindJSON(&schedule); err != nil {
		utils.RespondBadRequest(ctx, err.Error())
		return
	}

	schedule.ID = id
	result, err := c.usecase.Update(ctx.Request.Context(), &schedule)
	if err != nil {
		utils.RespondBadRequest(ctx, err.Error())
		return
	}

	utils.RespondOK(ctx, result)
}

// Delete handles DELETE /stylist-schedules/:id
func (c *StylistScheduleController) Delete(ctx *gin.Context) {
	id := ctx.Param("id")

	if err := c.usecase.Delete(ctx.Request.Context(), id); err != nil {
		utils.RespondBadRequest(ctx, err.Error())
		return
	}

	utils.RespondNoContent(ctx)
}

// ListByStylistID handles GET /stylists/:stylist_id/schedules
func (c *StylistScheduleController) ListByStylistID(ctx *gin.Context) {
	stylistID := ctx.Param("id")

	result, err := c.usecase.ListByStylistID(ctx.Request.Context(), stylistID)
	if err != nil {
		utils.RespondBadRequest(ctx, err.Error())
		return
	}

	utils.RespondOK(ctx, result)
}
