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
	products.Get("/:id", middleware.JWTAuthMiddleware(), handler.Get)
	products.Post("", middleware.JWTAuthMiddleware(), handler.Create)
	products.Put("/:id", middleware.JWTAuthMiddleware(), handler.Update)
	products.Delete("/:id", middleware.JWTAuthMiddleware(), handler.Delete)
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
// @Param 		 files formData file false "File to upload" format(multipart/form-data)
// @Success      200  {object}  []models.MyDrive
// @Failure      400  {object}  models.ResponseError
// @Failure      500  {object}  models.ResponseError
// @Security 	 BearerAuth
// @Router       /v1/drives [post]
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

// Update
// @Summary      Rename file
// @Description  Rename file
// @Tags         My Drive
// @Accept       json
// @Produce      json
// @Param        id   path      string  true  "ID"
// @Param 		 body formData models.MyDriveRenameInput true "Body"
// @Success      200  {object}  models.MyDrive
// @Failure      404  {object}  models.ResponseError
// @Failure      422  {object}  models.ResponseHTTP
// @Failure      500  {object}  models.ResponseError
// @Security 	 BearerAuth
// @Router       /v1/drives/{id} [put]
func (h *MyDriveHandler) Update(c *fiber.Ctx) error {
	var payload models.MyDriveRenameInput
	res := models.ResponseHTTP{
		Code:    fiber.StatusOK,
		Message: "Request has been processed successfully",
	}

	if err := c.BodyParser(&payload); err != nil {
		res.Code = fiber.StatusBadRequest
		res.Message = err.Error()
		return c.Status(res.Code).JSON(res)
	}

	// form POST validation
	errD := models.ValidateStruct(payload)
	if errD.Errors != nil {
		return c.Status(errD.Code).JSON(errD)
	}

	obj, err := h.uCase.Update(c, payload)
	if err != nil {
		res.Code = err.Code
		res.Message = err.Message
		return c.Status(res.Code).JSON(res)
	}

	return c.Status(res.Code).JSON(obj)
}

// Get
// @Summary      Get file
// @Description  Get file
// @Tags         My Drive
// @Accept       json
// @Produce      json
// @Param        id   path      string  true  "ID"
// @Success      200  {object}  models.MyDrive
// @Failure      400  {object}  models.ResponseError
// @Failure      500  {object}  models.ResponseError
// @Security 	 BearerAuth
// @Router       /v1/drives/{id} [get]
func (h *MyDriveHandler) Get(c *fiber.Ctx) error {
	res := models.ResponseHTTP{
		Code:    fiber.StatusOK,
		Message: "Request has been processed successfully",
	}

	obj, err := h.uCase.Get(c)
	if err != nil {
		res.Code = err.Code
		res.Message = err.Message
		return c.Status(res.Code).JSON(res)
	}

	return c.Status(res.Code).JSON(obj)
}

// Delete
// @Summary      Delete file
// @Description  Delete file
// @Tags         My Drive
// @Accept       json
// @Produce      json
// @Param        id   path      string  true  "ID"
// @Success      200  {object}  models.ResponseSuccess
// @Failure      400  {object}  models.ResponseError
// @Failure      500  {object}  models.ResponseError
// @Security 	 BearerAuth
// @Router       /v1/drives/{id} [delete]
func (h *MyDriveHandler) Delete(c *fiber.Ctx) error {
	res := models.ResponseHTTP{
		Code:    fiber.StatusOK,
		Message: "Request has been processed successfully",
	}

	if err := h.uCase.Delete(c); err != nil {
		res.Code = err.Code
		res.Message = err.Message
		return c.Status(res.Code).JSON(res)
	}

	return c.Status(res.Code).JSON(res)
}
