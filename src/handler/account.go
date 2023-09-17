package handler

import (
	middleware "myapp/pkg/middleware"
	"myapp/src/models"

	"github.com/gofiber/fiber/v2"
)

type AccountHandler struct {
	userUsecase models.UserUsecase
}

func NewAccountHandler(r fiber.Router, uc models.UserUsecase) {
	handler := &AccountHandler{
		userUsecase: uc,
	}

	// ROUTES
	acc := r.Group("/accounts")

	// private API
	acc.Get("/me", middleware.JWTAuthMiddleware(), handler.GetMe)
	acc.Post("/change-password", middleware.JWTAuthMiddleware(), handler.ChangePassword)
	acc.Put("/update", middleware.JWTAuthMiddleware(), handler.UpdateProfile)
	// acc.Post("/photo", middleware.JWTAuthMiddleware(), handler.UploadPhotoProfile)
}

func (h *AccountHandler) GetMe(c *fiber.Ctx) error {
	// A *model.User will eventually be added to context in middleware
	user, errLocal := c.Locals("user").(models.User)
	if !errLocal {
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

	user, errLocal := c.Locals("user").(models.User)
	if !errLocal {
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

func (h *AccountHandler) UpdateProfile(c *fiber.Ctx) error {
	var payload models.UpdateProfileInput

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

	user, err := h.userUsecase.Update(c, payload)
	if err != nil {
		return c.Status(err.Code).JSON(err)
	}

	return c.Status(fiber.StatusOK).JSON(user)
}

// func (h *AccountHandler) UploadPhotoProfile(c *fiber.Ctx) error {

// 	user, errLocal := c.Locals("user").(models.User)
// 	if !errLocal {
// 		errD := fiber.NewError(500, "Unable to extract user from request context for unknown reason")
// 		return c.Status(errD.Code).JSON(errD)
// 	}

// 	fileheader, err := c.FormFile("file")
// 	if err != nil {
// 		errD := fiber.NewError(500, err.Error())
// 		return c.Status(errD.Code).JSON(errD)
// 	}

// 	file, err := fileheader.Open()
// 	if err != nil {
// 		errD := fiber.NewError(500, err.Error())
// 		return c.Status(errD.Code).JSON(errD)
// 	}
// 	defer file.Close()

// 	buffer, err := io.ReadAll(file)
// 	if err != nil {
// 		errD := fiber.NewError(500, err.Error())
// 		return c.Status(errD.Code).JSON(errD)
// 	}

// 	err := h.userUsecase.ChangePassword(c.Context(), user, payload)
// 	if err != nil {
// 		return c.Status(err.Code).JSON(err)
// 	}

// 	return c.Status(fiber.StatusOK).JSON(fiber.Map{
// 		"code":    fiber.StatusOK,
// 		"message": "Your password has been changed successfully",
// 	})
// }
