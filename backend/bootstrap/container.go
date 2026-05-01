// bootstrap/container.go
package bootstrap

import (
	"backend/controller"
	"backend/internal/observability"
	"backend/repository"
	"backend/usecase"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Container struct {
	// Repositories
	BranchRepo          *repository.BranchRepository
	UserRepo            *repository.UserRepository
	StylistRepo         *repository.StylistRepository
	StylistScheduleRepo *repository.StylistScheduleRepository

	// Usecases
	BranchUC          *usecase.BranchUsecase
	UserUC            *usecase.UserUsecase
	StylistUC         *usecase.StylistUsecase
	StylistScheduleUC *usecase.StylistScheduleUsecase

	// Controllers
	BranchCtl          *controller.BranchController
	UserCtl            *controller.UserController
	StylistCtl         *controller.StylistController
	StylistScheduleCtl *controller.StylistScheduleController
}

func NewContainer(pool *pgxpool.Pool, telemetry *observability.Telemetry) *Container {
	// Repositories
	branchRepo := repository.NewBranchRepository(pool)
	userRepo := repository.NewUserRepository(pool)
	stylistRepo := repository.NewStylistRepository(pool)
	stylistScheduleRepo := repository.NewStylistScheduleRepository(pool)

	// Usecases
	branchUC := usecase.NewBranchUsecase(branchRepo)
	userUC := usecase.NewUserUsecase(userRepo)
	stylistUC := usecase.NewStylistUsecase(stylistRepo)
	stylistScheduleUC := usecase.NewStylistScheduleUsecase(stylistScheduleRepo)

	// Controllers
	branchCtl := controller.NewBranchController(branchUC)
	userCtl := controller.NewUserController(userUC)
	stylistCtl := controller.NewStylistController(stylistUC, telemetry)
	stylistScheduleCtl := controller.NewStylistScheduleController(stylistScheduleUC)

	return &Container{
		BranchRepo:          branchRepo,
		UserRepo:            userRepo,
		StylistRepo:         stylistRepo,
		StylistScheduleRepo: stylistScheduleRepo,
		BranchUC:            branchUC,
		UserUC:              userUC,
		StylistUC:           stylistUC,
		StylistScheduleUC:   stylistScheduleUC,
		BranchCtl:           branchCtl,
		UserCtl:             userCtl,
		StylistCtl:          stylistCtl,
		StylistScheduleCtl:  stylistScheduleCtl,
	}
}
