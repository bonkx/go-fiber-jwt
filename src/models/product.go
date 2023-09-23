package models

import (
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

type Product struct {
	gorm.Model
	Title       string  `json:"title" validate:"required,min=4"`
	Description string  `json:"description" gorm:"default:null"`
	Image       string  `json:"image" gorm:"default:null"`
	Price       float64 `json:"price" validate:"required"`
	IsEnable    bool    `json:"is_enable" gorm:"default:true"`
	// foreignkey User
	UserID uint
}

type ProductPopulateInput struct {
	UserID uint `json:"user_id" validate:"required"`
	Amount int  `json:"amount" validate:"required"`
}

type ProductInput struct {
	Title       string  `json:"title" validate:"required,min=4"`
	Description string  `json:"description"`
	Price       float64 `json:"price" validate:"required"`
	IsEnable    bool    `json:"is_enable"`
	UserID      uint    `json:"user_id"`
	// Image       string  `json:"image"`
}

type ProductUsecase interface {
	// USECASE

	// ADMIN ROLE
	PopulateProducts(userID uint, n int) *fiber.Error
}

type ProductRepository interface {
	// FUNTIONS

	// REPOS

	// ADMIN ROLE
	PopulateProducts(userID uint, n int) *fiber.Error
}
