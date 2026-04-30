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
	BranchRepo  *repository.BranchRepository
	UserRepo    *repository.UserRepository
	StylistRepo *repository.StylistRepository

	// Usecases
	BranchUC  *usecase.BranchUsecase
	UserUC    *usecase.UserUsecase
	StylistUC *usecase.StylistUsecase

	// Controllers
	BranchCtl  *controller.BranchController
	UserCtl    *controller.UserController
	StylistCtl *controller.StylistController
}

func NewContainer(pool *pgxpool.Pool, telemetry *observability.Telemetry) *Container {
	// Repositories
	branchRepo := repository.NewBranchRepository(pool)
	userRepo := repository.NewUserRepository(pool)
	stylistRepo := repository.NewStylistRepository(pool)

	// Usecases
	branchUC := usecase.NewBranchUsecase(branchRepo)
	userUC := usecase.NewUserUsecase(userRepo)
	stylistUC := usecase.NewStylistUsecase(stylistRepo)

	// Controllers
	branchCtl := controller.NewBranchController(branchUC)
	userCtl := controller.NewUserController(userUC)
	stylistCtl := controller.NewStylistController(stylistUC, telemetry)

	return &Container{
		BranchRepo:  branchRepo,
		UserRepo:    userRepo,
		StylistRepo: stylistRepo,
		BranchUC:    branchUC,
		UserUC:      userUC,
		StylistUC:   stylistUC,
		BranchCtl:   branchCtl,
		UserCtl:     userCtl,
		StylistCtl:  stylistCtl,
	}
}
