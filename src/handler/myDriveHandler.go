package handler

import (
	"myapp/pkg/middleware"
	"myapp/src/models"

	"github.com/gofiber/fiber/v2"
)

type MyDriveHandler struct {
	uCase models.MyDriveUsecase
}

func NewMyDriveHandler(r fiber.Router, uc models.MyDriveUsecase) {
	handler := &MyDriveHandler{
		uCase: uc,
	}

	products := r.Group("/drives")

	products.Get("", middleware.JWTAuthMiddleware(), handler.MyDrives)
	products.Post("", middleware.JWTAuthMiddleware(), handler.Create)
}

// MyDrives
// @Summary      List of My Drive
// @Description  List of all files
// @Tags         My Drive
// @Accept       json
// @Produce      json
// @Success      200  {object}  response.Pagination
// @Failure      500  {object}  models.ResponseError
// @Security 	 BearerAuth
// @Router       /v1/drives [get]
func (h *MyDriveHandler) MyDrives(c *fiber.Ctx) error {
	res := models.ResponseHTTP{
		Code:    fiber.StatusOK,
		Message: "Request has been processed successfully",
	}

	pagination, err := h.uCase.MyDrive(c)
	if err != nil {
		res.Code = err.Code
		res.Message = err.Message
		return c.Status(res.Code).JSON(res)
	}

	return c.Status(res.Code).JSON(&pagination)
}

// Create
// @Summary      Upload File
// @Description  Upload new file
// @Tags         My Drive
// @Accept       multipart/form-data
// @Produce      json
// @Param 		 file formData file false "File to upload" format(multipart/form-data)
// @Success      200  {object}  models.MyDrive
// @Failure      400  {object}  models.ResponseError
// @Failure      500  {object}  models.ResponseError
// @Security 	 BearerAuth
// @Router       /v1/products [post]
func (h *MyDriveHandler) Create(c *fiber.Ctx) error {
	res := models.ResponseHTTP{
		Code:    fiber.StatusOK,
		Message: "Request has been processed successfully",
	}

	obj, err := h.uCase.Create(c)
	if err != nil {
		res.Code = err.Code
		res.Message = err.Message
		return c.Status(res.Code).JSON(res)
	}

	return c.Status(res.Code).JSON(obj)
}
