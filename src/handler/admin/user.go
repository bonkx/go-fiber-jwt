package admin

import (
	middleware "myapp/pkg/middleware"
	"myapp/pkg/utils"
	"myapp/src/models"

	"github.com/gofiber/fiber/v2"
)

type AdminUserHandler struct {
	userUsecase models.UserUsecase
}

func NewAdminUserHandler(r fiber.Router, uc models.UserUsecase) {
	handler := &AdminUserHandler{
		userUsecase: uc,
	}

	// ROUTES
	acc := r.Group("/user")

	// private API
	acc.Get("/me", middleware.AdminAuthMiddleware(), handler.GetMe)

	acc.Delete("/:id", middleware.AdminAuthMiddleware(), handler.DeleteUser)
	acc.Post("/restore", middleware.AdminAuthMiddleware(), handler.RestoreUser)
}

func (h *AdminUserHandler) GetMe(c *fiber.Ctx) error {
	user, errLocal := c.Locals("user").(models.User)
	if !errLocal {
		errD := fiber.NewError(500, "Unable to extract user from request context for unknown reason")
		return c.Status(errD.Code).JSON(errD)
	}

	return c.Status(fiber.StatusOK).JSON(user)
}

func (h *AdminUserHandler) DeleteUser(c *fiber.Ctx) error {
	uID := utils.StringToUint(c.Params("id"))

	if err := h.userUsecase.DeleteUser(c, uID); err != nil {
		errD := fiber.NewError(err.Code, err.Message)
		return c.Status(errD.Code).JSON(errD)
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"code":    fiber.StatusOK,
		"message": "Request has been processed successfully",
	})
}

func (h *AdminUserHandler) RestoreUser(c *fiber.Ctx) error {
	var payload models.EmailInput

	if err := c.BodyParser(&payload); err != nil {
		errD := fiber.NewError(fiber.StatusBadRequest, err.Error())
		return c.Status(errD.Code).JSON(errD)
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

	err := h.userUsecase.RestoreUser(c, payload.Email)
	if err != nil {
		return c.Status(err.Code).JSON(err)
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"code":    fiber.StatusOK,
		"message": "Request has been processed successfully",
	})
}
