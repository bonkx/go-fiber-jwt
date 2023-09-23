package admin

import (
	"myapp/pkg/middleware"
	"myapp/src/models"

	"github.com/gofiber/fiber/v2"
)

type AdminProductHandler struct {
	pUsecase models.ProductUsecase
}

func NewAdminProductHandler(r fiber.Router, uc models.ProductUsecase) {
	handler := &AdminProductHandler{
		pUsecase: uc,
	}

	// ROUTES
	acc := r.Group("/products")

	// private API
	acc.Post("/populate", middleware.AdminAuthMiddleware(), handler.PopulateProducts)

}

func (h *AdminProductHandler) PopulateProducts(c *fiber.Ctx) error {
	var payload models.ProductPopulateInput
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
	errors := models.ValidateStruct(payload)
	if errors != nil {
		res.Code = fiber.ErrUnprocessableEntity.Code
		res.Message = fiber.ErrUnprocessableEntity.Message
		res.Errors = errors
		return c.Status(res.Code).JSON(res)
	}

	err := h.pUsecase.PopulateProducts(payload.UserID, payload.Amount)
	if err != nil {
		res.Code = err.Code
		res.Message = err.Message
		return c.Status(res.Code).JSON(res)
	}

	return c.Status(res.Code).JSON(res)
}
