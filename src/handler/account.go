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

	// private API
	acc.Get("/me", middleware.JWTAuthMiddleware(), handler.GetMe)
	acc.Post("/change-password", middleware.JWTAuthMiddleware(), handler.ChangePassword)
}

func (h *AccountHandler) GetMe(c *fiber.Ctx) error {
	// A *model.User will eventually be added to context in middleware
	user, err := c.Locals("user").(models.User)
	if !err {
		errD := fiber.NewError(500, "Unable to extract user from request context for unknown reason")
		return c.Status(errD.Code).JSON(errD)
	}

	return c.Status(fiber.StatusOK).JSON(user)
}

func (h *AccountHandler) ChangePassword(c *fiber.Ctx) error {
	var payload models.ChangePasswordInput

	if err := c.BodyParser(&payload); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(
			fiber.NewError(fiber.StatusBadRequest, err.Error()),
		)
	}

	if payload.Password != payload.PasswordConfirm {
		return c.Status(fiber.StatusBadRequest).JSON(
			fiber.NewError(400, "Passwords do not match!"),
		)
	}

	// form POST validations
	errors := models.ValidateStruct(payload)
	if errors != nil {
		errD := models.ErrorDetailsResponse{
			Code:    fiber.ErrUnprocessableEntity.Code,
			Message: fiber.ErrUnprocessableEntity.Message,
			Errors:  errors,
		}
		return c.Status(errD.Code).JSON(errD)
	}

	user, errUser := c.Locals("user").(models.User)
	if !errUser {
		errD := fiber.NewError(500, "Unable to extract user from request context for unknown reason")
		return c.Status(errD.Code).JSON(errD)
	}

	err := h.userUsecase.ChangePassword(c.Context(), user, payload)
	if err != nil {
		return c.Status(err.Code).JSON(err)
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"code":    fiber.StatusOK,
		"message": "Your password has been changed successfully",
	})
}
