package routes

import (
	_handler "myapp/src/handler"
	_repo "myapp/src/repository"
	_useCase "myapp/src/usecase"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

// APIRoutes func for describe group of private routes.
func APIRoutes(a *fiber.App, db *gorm.DB) {
	// Create routes group.
	v1 := a.Group("/api/v1")

	repoUser := _repo.NewUserRepository(db)
	ucUser := _useCase.NewUserUsecase(repoUser)

	// ROUTES
	// AuthRoute
	_handler.NewAuthHandler(v1, ucUser)

}
