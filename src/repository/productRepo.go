package repository

import (
	"fmt"
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

// PopulateProducts implements models.ProductRepository.
func (r *ProductRepository) PopulateProducts(userID uint, n int) *fiber.Error {
	for i := 0; i < n; i++ {
		r.DB.Create(&models.Product{
			Title:       fmt.Sprintf("%s Product No.%s", faker.Word(), strconv.Itoa(i+1)),
			Description: faker.Paragraph(),
			Image:       fmt.Sprintf("https://loremflickr.com/320/240/product?%s", faker.UUIDDigit()),
			// Price:       rand.Float64() + 10,
			Price:  utils.GetRandFloat(5.0, 100.0),
			UserID: userID,
		})
	}

	return nil
}
