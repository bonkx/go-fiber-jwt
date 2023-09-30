package models

import (
	"encoding/json"
	"fmt"
	"myapp/pkg/response"
	"os"
	"strings"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

type Product struct {
	gorm.Model
	Title       string  `json:"title" validate:"required,min=4"`
	Description string  `json:"description" gorm:"default:null"`
	Image       *string `json:"image" gorm:"default:null"`
	Price       float64 `json:"price" validate:"required"`
	IsEnable    bool    `json:"is_enable" gorm:"default:true"`
	// foreignkey User
	UserID uint
	User   User `gorm:"foreignkey:UserID;constraint:OnDelete:CASCADE;" json:"user"`
}

func (md Product) MarshalJSON() ([]byte, error) {
	type Alias Product
	var image *string = md.Image
	// check image is empty
	if image != nil {
		// check image startwith http/https
		if !strings.HasPrefix(*image, "http") {
			*image = fmt.Sprintf("%s/%s", os.Getenv("CLIENT_ORIGIN"), *md.Image)
		}
	}

	aux := struct {
		Alias
		Image *string `json:"image"`
	}{
		Alias: (Alias)(md),
		Image: image,
	}
	return json.Marshal(aux)
}

type ProductPopulateInput struct {
	UserID uint `json:"user_id" validate:"required"`
	Amount int  `json:"amount" validate:"required,number"`
}

type ProductInput struct {
	Title       string  `json:"title" form:"title" validate:"required,min=4"`
	Description string  `json:"description" form:"description"`
	Price       float64 `json:"price" form:"price" validate:"required,number"`
}

type ProductUsecase interface {
	// USECASE
	MyProduct(c *fiber.Ctx) (*response.Pagination, *fiber.Error)
	ListProduct(c *fiber.Ctx) (*response.Pagination, *fiber.Error)
	GetProduct(c *fiber.Ctx) (Product, *fiber.Error)
	Create(c *fiber.Ctx, payload ProductInput) (Product, *fiber.Error)
	Update(c *fiber.Ctx, id uint, payload ProductInput) (Product, *fiber.Error)
	Delete(c *fiber.Ctx, id uint) *fiber.Error

	// ADMIN ROLE
	PopulateProducts(userID uint, n int) *fiber.Error
}

type ProductRepository interface {
	// FUNTIONS

	// REPOS
	MyProduct(user User, param response.ParamsPagination) (*response.Pagination, *fiber.Error)
	ListProduct(param response.ParamsPagination) (*response.Pagination, *fiber.Error)
	GetProduct(id uint) (Product, *fiber.Error)
	Create(obj Product) (Product, *fiber.Error)
	Update(obj Product) (Product, *fiber.Error)
	Delete(obj Product) *fiber.Error

	// ADMIN ROLE
	PopulateProducts(userID uint, n int) *fiber.Error
}
