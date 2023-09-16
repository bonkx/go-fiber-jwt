package handler

import (
	middleware "myapp/pkg/middleware"
	"myapp/src/models"

	"github.com/gofiber/fiber/v2"
)

type AccountHandler struct {
	userUsecase models.UserUsecase
}

func NewAccountHandler(r fiber.Router, userUsecase models.UserUsecase) {
	handler := &AccountHandler{
		userUsecase: userUsecase,
	}

	// ROUTES
	acc := r.Group("/accounts")

	acc.Get("/me", middleware.JWTAuthMiddleware(), handler.GetMe)
}

func (h *AccountHandler) GetMe(c *fiber.Ctx) error {
	// A *model.User will eventually be added to context in middleware
	user, err := c.Locals("user").(models.User)
	if !err {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Unable to extract user from request context for unknown reason",
		})
	}

	return c.Status(fiber.StatusOK).JSON(user)
}
