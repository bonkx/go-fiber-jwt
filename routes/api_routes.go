package routes

import (
	_handler "myapp/src/handler"
	_admin "myapp/src/handler/admin"
	_repo "myapp/src/repository"
	_useCase "myapp/src/usecase"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

// APIRoutes func for describe group of private routes.
func APIRoutes(a *fiber.App, db *gorm.DB) {
	// Create routes group.
	v1 := a.Group("/api/v1")

	// register All REPOSITORY
	repoUser := _repo.NewUserRepository(db)

	// register All USECASE
	ucUser := _useCase.NewUserUsecase(repoUser)

	// ROUTES
	// Auth route
	_handler.NewAuthHandler(v1, ucUser)
	// Account route
	_handler.NewAccountHandler(v1, ucUser)

	// ADMIN Routes
	admin := v1.Group("/admin")
	_admin.NewAdminUserHandler(admin, ucUser)
	// test routes
	_handler.NewEmailHandler(a, ucUser)
}
