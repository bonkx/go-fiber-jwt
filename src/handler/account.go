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

	// public API
	acc.Post("/forgot-password", handler.ForgotPassword)
	acc.Post("/forgot-password-otp", handler.ForgotPasswordOTP)

	// private API
	acc.Get("/me", middleware.JWTAuthMiddleware(), handler.GetMe)
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

func (h *AccountHandler) ForgotPassword(c *fiber.Ctx) error {
	var payload models.EmailInput

	if err := c.BodyParser(&payload); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(
			fiber.NewError(fiber.StatusBadRequest, err.Error()),
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

	err := h.userUsecase.ForgotPassword(c.Context(), payload)
	if err != nil {
		return c.Status(err.Code).JSON(err)
	}

	message := "We sent an email with a OTP code to " + payload.Email
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"code":    fiber.StatusOK,
		"message": message,
	})
}

func (h *AccountHandler) ForgotPasswordOTP(c *fiber.Ctx) error {
	var payload models.OTPInput

	if err := c.BodyParser(&payload); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(
			fiber.NewError(fiber.StatusBadRequest, err.Error()),
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

	refNo, err := h.userUsecase.ForgotPasswordOTP(c.Context(), payload)
	if err != nil {
		return c.Status(err.Code).JSON(err)
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"code":         fiber.StatusOK,
		"message":      "OTP verification successful",
		"reference_no": refNo,
	})
}
