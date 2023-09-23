package usecase

import (
	"myapp/src/models"

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
