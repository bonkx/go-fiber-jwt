package repository

import (
	"fmt"
	"myapp/pkg/response"
	"myapp/pkg/utils"
	"myapp/src/models"
	"strconv"

	"github.com/go-faker/faker/v4"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

type ProductRepository struct {
	DB *gorm.DB
}

// NewProductRepository will create an object that represent the models.ProductRepository interface
func NewProductRepository(Conn *gorm.DB) models.ProductRepository {
	return &ProductRepository{Conn}
}

// Create implements models.ProductRepository.
func (r *ProductRepository) Create(product *models.Product) (*models.Product, *fiber.Error) {
	result := r.DB.Create(&product)
	if result.Error != nil {
		return nil, fiber.NewError(500, result.Error.Error())
	}

	return product, nil
}

// GetProduct implements models.ProductRepository.
func (r *ProductRepository) GetProduct(id uint) (*models.Product, *fiber.Error) {
	var obj *models.Product
	result := r.DB.Preload("User.UserProfile.Status").Where("id = ?", id).Find(&obj)
	if result.Error != nil {
		return nil, fiber.NewError(500, result.Error.Error())
	}
	if result.RowsAffected == 0 {
		return nil, fiber.NewError(404, utils.ERR_DATA_NOT_FOUND)
	}
	return obj, nil
}

// MyProduct implements models.ProductRepository.
func (r *ProductRepository) MyProduct(user models.User, param response.ParamsPagination) (*response.Pagination, *fiber.Error) {
	var data []*models.Product
	// var count int64
	var pagination response.Pagination

	db := r.DB.Preload("User.UserProfile.Status").Where("user_id = ?", user.ID)
	// db := r.DB.Preload("User").Where("user_id = ?", user.ID)
	// db := r.DB.Preload(clause.Associations).Where("user_id = ?", user.ID)

	if param.Search != "" {
		// search data based on title, description
		db = db.Where("title ILIKE ?", "%"+param.Search+"%").
			Or("description ILIKE ?", "%"+param.Search+"%")
	}

	// 	fill all params pagination
	pagination.Sort = param.SortQuery
	pagination.Page = param.Page
	pagination.Limit = param.Limit

	err := db.Scopes(response.Paginate(data, &pagination, db)).Find(&data).Error
	if err != nil {
		return nil, fiber.NewError(500, db.Error.Error())
	}

	pagination.Data = data

	return &pagination, nil
}

// ListProduct implements models.ProductRepository.
func (r *ProductRepository) ListProduct(param response.ParamsPagination) (*response.Pagination, *fiber.Error) {
	var data []*models.Product
	// var count int64
	var pagination response.Pagination

	db := r.DB.Preload("User.UserProfile.Status")

	if param.Search != "" {
		// search data based on title, description
		db = db.Where("title ILIKE ?", "%"+param.Search+"%").
			Or("description ILIKE ?", "%"+param.Search+"%")
	}

	// 	fill all params pagination
	pagination.Sort = param.SortQuery
	pagination.Page = param.Page
	pagination.Limit = param.Limit

	err := db.Scopes(response.Paginate(data, &pagination, db)).Find(&data).Error
	if err != nil {
		return nil, fiber.NewError(500, db.Error.Error())
	}

	pagination.Data = data

	return &pagination, nil
}

// PopulateProducts implements models.ProductRepository.
func (r *ProductRepository) PopulateProducts(userID uint, n int) *fiber.Error {
	for i := 0; i < n; i++ {
		image := fmt.Sprintf("https://loremflickr.com/320/240/product?%s", faker.UUIDDigit())
		r.DB.Create(&models.Product{
			Title:       fmt.Sprintf("%s Product No.%s", faker.Word(), strconv.Itoa(i+1)),
			Description: faker.Paragraph(),
			Image:       &image,
			Price:       utils.GetRandFloat(5.0, 100.0),
			IsEnable:    false,
			UserID:      userID,
			User:        models.User{},
		})
	}

	return nil
}
