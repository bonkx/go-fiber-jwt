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
	r.Get("/me", middleware.AdminAuthMiddleware(), handler.GetMe)

	users := r.Group("/users")

	users.Get("", middleware.AdminAuthMiddleware(), handler.ListUser)

	users.Delete("/:id", middleware.AdminAuthMiddleware(), handler.DeleteUser)
	users.Delete("/:id/unscoped", middleware.AdminAuthMiddleware(), handler.PermanentDeleteUser)
	users.Post("/restore", middleware.AdminAuthMiddleware(), handler.RestoreUser)
}

func (h *AdminUserHandler) GetMe(c *fiber.Ctx) error {
	user, errLocal := c.Locals("user").(models.User)
	if !errLocal {
		res := models.ResponseHTTP{
			Code:    fiber.StatusInternalServerError,
			Message: "Unable to extract user from request context for unknown reason",
		}
		return c.Status(res.Code).JSON(res)
	}

	return c.Status(fiber.StatusOK).JSON(user)
}

func (h *AdminUserHandler) ListUser(c *fiber.Ctx) error {
	res := models.ResponseHTTP{
		Code:    fiber.StatusOK,
		Message: "Request has been processed successfully",
	}

	pagination, data, err := h.userUsecase.ListUser(c)
	if err != nil {
		res.Code = err.Code
		res.Message = err.Message
		return c.Status(res.Code).JSON(res)
	}

	// 	with no pagination (support searching and sorting)
	if c.Query("no_page") != "" {
		return c.Status(fiber.StatusOK).JSON(&data)
	}

	return c.Status(fiber.StatusOK).JSON(&pagination)
}

func (h *AdminUserHandler) DeleteUser(c *fiber.Ctx) error {
	res := models.ResponseHTTP{
		Code:    fiber.StatusOK,
		Message: "Request has been processed successfully",
	}

	id := utils.StringToUint(c.Params("id"))

	if err := h.userUsecase.DeleteUser(c, id); err != nil {
		res.Code = err.Code
		res.Message = err.Message
		return c.Status(res.Code).JSON(res)
	}

	return c.Status(res.Code).JSON(res)
}

func (h *AdminUserHandler) PermanentDeleteUser(c *fiber.Ctx) error {
	res := models.ResponseHTTP{
		Code:    fiber.StatusOK,
		Message: "Request has been processed successfully",
	}

	id := utils.StringToUint(c.Params("id"))

	if err := h.userUsecase.PermanentDeleteUser(c, id); err != nil {
		res.Code = err.Code
		res.Message = err.Message
		return c.Status(res.Code).JSON(res)
	}

	return c.Status(res.Code).JSON(res)
}

func (h *AdminUserHandler) RestoreUser(c *fiber.Ctx) error {
	var payload models.EmailInput
	res := models.ResponseHTTP{
		Code:    fiber.StatusOK,
		Message: "Request has been processed successfully",
	}

	if err := c.BodyParser(&payload); err != nil {
		res.Code = fiber.StatusBadRequest
		res.Message = err.Error()
		return c.Status(res.Code).JSON(res)
	}

	// form POST validations
	errD := models.ValidateStruct(payload)
	if errD.Errors != nil {
		return c.Status(errD.Code).JSON(errD)
	}

	err := h.userUsecase.RestoreUser(c, payload.Email)
	if err != nil {
		res.Code = err.Code
		res.Message = err.Message
		return c.Status(res.Code).JSON(res)
	}

	return c.Status(res.Code).JSON(res)
}
