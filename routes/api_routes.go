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
	repoProduct := _repo.NewProductRepository(db)
	repoMyDrive := _repo.NewMyDriveRepository(db)

	// register All USECASE
	ucUser := _useCase.NewUserUsecase(repoUser)
	ucProduct := _useCase.NewProductUsecase(repoProduct, repoUser)
	ucMyDrive := _useCase.NewMyDriveUsecase(repoMyDrive, repoUser)

	// ROUTES
	_handler.NewAuthHandler(v1, ucUser)
	_handler.NewAccountHandler(v1, ucUser)
	_handler.NewProductHandler(v1, ucProduct)
	_handler.NewMyDriveHandler(v1, ucMyDrive)

	// ADMIN Routes
	admin := v1.Group("/admin")
	_admin.NewAdminUserHandler(admin, ucUser)
	_admin.NewAdminProductHandler(admin, ucProduct)
	// test routes
	_handler.NewEmailHandler(a, ucUser)
}
