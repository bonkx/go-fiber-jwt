package usecase

import (
	"myapp/pkg/response"
	"myapp/pkg/utils"
	"myapp/src/models"
	"strconv"

	"github.com/gofiber/fiber/v2"
)

type ProductUsecase struct {
	pRepo models.ProductRepository
	uRepo models.UserRepository
}

// NewProductUsecase will create an object that represent the models.ProductUsecase interface
func NewProductUsecase(product models.ProductRepository, user models.UserRepository) models.ProductUsecase {
	return &ProductUsecase{
		pRepo: product,
		uRepo: user,
	}
}

// Delete implements models.ProductUsecase.
func (uc *ProductUsecase) Delete(c *fiber.Ctx, id uint) *fiber.Error {
	user, errLocal := c.Locals("user").(models.User)
	if !errLocal {
		return fiber.NewError(500, utils.ERR_CURRENT_USER_NOT_FOUND)
	}

	// get data
	obj, err := uc.pRepo.GetProduct(id)
	if err != nil {
		return err
	}

	// check the owner of data
	if obj.UserID != user.ID {
		return fiber.NewError(403, utils.ERR_FORBIDDEN_UPDATE)
	}

	// deleted obj
	err = uc.pRepo.Delete(obj)
	if err != nil {
		return err
	}

	return nil
}

// Update implements models.ProductUsecase.
func (uc *ProductUsecase) Update(c *fiber.Ctx, id uint, payload models.ProductInput) (models.Product, *fiber.Error) {
	obj := models.Product{}
	// get product data
	obj, err := uc.pRepo.GetProduct(id)
	if err != nil {
		return obj, err
	}

	user, errLocal := c.Locals("user").(models.User)
	if !errLocal {
		return obj, fiber.NewError(500, utils.ERR_CURRENT_USER_NOT_FOUND)
	}

	// check the owner of data
	if obj.UserID != user.ID {
		return obj, fiber.NewError(403, utils.ERR_FORBIDDEN_UPDATE)
	}

	// MultipartForm POST
	if form, err := c.MultipartForm(); err == nil {
		files := form.File["image"]

		// Loop through files:
		for _, file := range files {
			// // Save the files to disk:
			imageUrl, errFile := utils.ImageUpload(c, file, "products")
			if errFile != nil {
				return obj, fiber.NewError(500, errFile.Error())
			}

			// update user photo path
			obj.Image = &imageUrl
		}
	}

	// fill update data
	obj.Title = payload.Title
	obj.Description = payload.Description
	obj.Price = payload.Price

	// do update user
	obj, err = uc.pRepo.Update(obj)
	if err != nil {
		return obj, err
	}

	return obj, nil
}

// Create implements models.ProductUsecase.
func (uc *ProductUsecase) Create(c *fiber.Ctx, payload models.ProductInput) (models.Product, *fiber.Error) {
	var obj models.Product

	user, errLocal := c.Locals("user").(models.User)
	if !errLocal {
		return obj, fiber.NewError(500, utils.ERR_CURRENT_USER_NOT_FOUND)
	}

	// fill the owner
	obj.UserID = user.ID

	// MultipartForm POST
	if form, err := c.MultipartForm(); err == nil {
		files := form.File["image"]

		// Loop through files:
		for _, file := range files {
			// // Save the files to disk:
			imageUrl, errFile := utils.ImageUpload(c, file, "products")
			if errFile != nil {
				return obj, fiber.NewError(500, errFile.Error())
			}

			// update user photo path
			obj.Image = &imageUrl
		}
	}

	// fill the form data
	obj.Title = payload.Title
	obj.Description = payload.Description
	obj.Price = payload.Price

	// save the data
	obj, err := uc.pRepo.Create(obj)
	if err != nil {
		return obj, err
	}

	return obj, nil
}

// GetProduct implements models.ProductUsecase.
func (uc *ProductUsecase) GetProduct(c *fiber.Ctx) (models.Product, *fiber.Error) {
	id := utils.StringToUint(c.Params("id"))

	obj, err := uc.pRepo.GetProduct(id)
	if err != nil {
		return obj, err
	}

	return obj, nil
}

// MyProduct implements models.ProductUsecase.
func (uc *ProductUsecase) MyProduct(c *fiber.Ctx) (*response.Pagination, *fiber.Error) {
	user, errLocal := c.Locals("user").(models.User)
	if !errLocal {
		return nil, fiber.NewError(500, utils.ERR_CURRENT_USER_NOT_FOUND)
	}

	// 	Parse the query parameters
	search := c.Query("search")
	sortBy := c.Query("sort", "id|desc")
	page := c.Query("page", "1")
	limit := c.Query("per_page", "10")

	// Convert the page and limit to integers
	pageInt, _ := strconv.Atoi(page)
	limitInt, _ := strconv.Atoi(limit)

	sortQuery, errSort := utils.ValidateAndReturnSortQuery(sortBy)
	// log.Print(sortQuery)
	if errSort != nil {
		errD := fiber.NewError(fiber.StatusInternalServerError, errSort.Error())
		return nil, errD
	}

	// make param pagination struct
	pagParam := response.ParamsPagination{
		Page:      pageInt,
		Limit:     limitInt,
		SortQuery: sortQuery,
		Search:    search,
		NoPage:    c.Query("no_page"),
	}

	pagination, err := uc.pRepo.MyProduct(user, pagParam)
	if err != nil {
		return nil, err
	}
	return pagination, nil
}

// ListProduct implements models.ProductUsecase.
func (uc *ProductUsecase) ListProduct(c *fiber.Ctx) (*response.Pagination, *fiber.Error) {
	// 	Parse the query parameters
	search := c.Query("search")
	sortBy := c.Query("sort", "id|desc")
	page := c.Query("page", "1")
	limit := c.Query("per_page", "10")

	// Convert the page and limit to integers
	pageInt, _ := strconv.Atoi(page)
	limitInt, _ := strconv.Atoi(limit)

	sortQuery, errSort := utils.ValidateAndReturnSortQuery(sortBy)
	// log.Print(sortQuery)
	if errSort != nil {
		errD := fiber.NewError(fiber.StatusInternalServerError, errSort.Error())
		return nil, errD
	}

	// make param pagination struct
	pagParam := response.ParamsPagination{
		Page:      pageInt,
		Limit:     limitInt,
		SortQuery: sortQuery,
		Search:    search,
		NoPage:    c.Query("no_page"),
	}

	pagination, err := uc.pRepo.ListProduct(pagParam)
	if err != nil {
		return nil, err
	}
	return pagination, nil
}

// PopulateProducts implements models.ProductUsecase.
func (uc *ProductUsecase) PopulateProducts(userID uint, n int) *fiber.Error {
	// find user
	_, err := uc.uRepo.FindUserById(userID)
	if err != nil {
		return err
	}

	err = uc.pRepo.PopulateProducts(userID, n)
	if err != nil {
		return err
	}
	return nil
}
