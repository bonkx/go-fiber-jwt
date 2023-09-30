package handler

import (
	"myapp/pkg/middleware"
	"myapp/src/models"

	"github.com/gofiber/fiber/v2"
)

type ProductHandler struct {
	uCase models.ProductUsecase
}

func NewProductHandler(r fiber.Router, uc models.ProductUsecase) {
	handler := &ProductHandler{
		uCase: uc,
	}

	products := r.Group("/products")

	products.Get("/my-products", middleware.JWTAuthMiddleware(), handler.MyProduct)
	products.Get("", middleware.JWTAuthMiddleware(), handler.ListProduct)
	products.Get("/:id", middleware.JWTAuthMiddleware(), handler.GetProduct)
	products.Post("", middleware.JWTAuthMiddleware(), handler.CreateProduct)
	products.Put("/:id", middleware.JWTAuthMiddleware(), handler.UpdateProduct)
	products.Delete("/:id", middleware.JWTAuthMiddleware(), handler.DeleteProduct)
}

// ListProduct
// @Summary      List of Product
// @Description  List of all Products
// @Tags         Products
// @Accept       json
// @Produce      json
// @Success      200  {object}  response.Pagination
// @Failure      500  {object}  models.ResponseError
// @Security 	 BearerAuth
// @Router       /v1/products [get]
func (h *ProductHandler) ListProduct(c *fiber.Ctx) error {
	res := models.ResponseHTTP{
		Code:    fiber.StatusOK,
		Message: "Request has been processed successfully",
	}

	pagination, err := h.uCase.ListProduct(c)
	if err != nil {
		res.Code = err.Code
		res.Message = err.Message
		return c.Status(res.Code).JSON(res)
	}

	return c.Status(res.Code).JSON(&pagination)
}

// MyProduct
// @Summary      My Product
// @Description  List of all my Products
// @Tags         Products
// @Accept       json
// @Produce      json
// @Success      200  {object}  response.Pagination
// @Failure      500  {object}  models.ResponseError
// @Security 	 BearerAuth
// @Router       /v1/products/my-product [get]
func (h *ProductHandler) MyProduct(c *fiber.Ctx) error {
	res := models.ResponseHTTP{
		Code:    fiber.StatusOK,
		Message: "Request has been processed successfully",
	}

	pagination, err := h.uCase.MyProduct(c)
	if err != nil {
		res.Code = err.Code
		res.Message = err.Message
		return c.Status(res.Code).JSON(res)
	}

	return c.Status(res.Code).JSON(&pagination)
}

// GetProduct
// @Summary      GetProduct
// @Description  Get details of product
// @Tags         Products
// @Accept       json
// @Produce      json
// @Param        id   path      int  true  "Product ID"
// @Success      200  {object}  models.Product
// @Success      404  {object}  models.ResponseError
// @Failure      500  {object}  models.ResponseError
// @Security 	 BearerAuth
// @Router       /v1/products/{id} [get]
func (h *ProductHandler) GetProduct(c *fiber.Ctx) error {
	res := models.ResponseHTTP{
		Code:    fiber.StatusOK,
		Message: "Request has been processed successfully",
	}

	obj, err := h.uCase.GetProduct(c)
	if err != nil {
		res.Code = err.Code
		res.Message = err.Message
		return c.Status(res.Code).JSON(res)
	}

	return c.Status(res.Code).JSON(obj)
}

// CreateProduct
// @Summary      Create Product
// @Description  Create new product
// @Tags         Products
// @Accept       json
// @Accept       multipart/form-data
// @Produce      json
// @Param 		 body formData models.ProductInput true "Body"
// @Param 		 file formData file false "File to upload" format(multipart/form-data)
// @Success      200  {object}  models.Product
// @Failure      400  {object}  models.ResponseError
// @Failure      422  {object}  models.ResponseHTTP
// @Failure      500  {object}  models.ResponseError
// @Security 	 BearerAuth
// @Router       /v1/products [post]
func (h *ProductHandler) CreateProduct(c *fiber.Ctx) error {
	var payload models.ProductInput
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

	obj, err := h.uCase.Create(c, payload)
	if err != nil {
		res.Code = err.Code
		res.Message = err.Message
		return c.Status(res.Code).JSON(res)
	}

	return c.Status(res.Code).JSON(obj)
}

// UpdateProduct
// @Summary      Update Product
// @Description  Update product's data
// @Tags         Products
// @Accept       json
// @Accept       multipart/form-data
// @Produce      json
// @Param        id   path      int  true  "Product ID"
// @Param 		 body formData models.ProductInput true "Body"
// @Param 		 file formData file false "File to upload" format(multipart/form-data)
// @Success      200  {object}  models.Product
// @Failure      404  {object}  models.ResponseError
// @Failure      422  {object}  models.ResponseHTTP
// @Failure      500  {object}  models.ResponseError
// @Security 	 BearerAuth
// @Router       /v1/products/{id} [put]
func (h *ProductHandler) UpdateProduct(c *fiber.Ctx) error {
	var payload models.ProductInput
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

// DeleteProduct
// @Summary      Delete Product
// @Description  Delete product's data
// @Tags         Products
// @Accept       json
// @Produce      json
// @Param        id   path      int  true  "Product ID"
// @Success      200  {object}  models.ResponseSuccess
// @Failure      404  {object}  models.ResponseError
// @Failure      500  {object}  models.ResponseError
// @Security 	 BearerAuth
// @Router       /v1/products/{id} [delete]
func (h *ProductHandler) DeleteProduct(c *fiber.Ctx) error {
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
