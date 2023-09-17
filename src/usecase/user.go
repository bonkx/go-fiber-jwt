package usecase

import (
	"context"
	"myapp/src/models"

	"github.com/gofiber/fiber/v2"
)

type UserUsecase struct {
	userRepo models.UserRepository
}

// Logout implements models.UserUsecase.
func (uc *UserUsecase) Logout(authD *models.AccessDetails) *fiber.Error {
	if err := uc.userRepo.DeleteToken(authD); err != nil {
		return err
	}
	return nil
}

// ResendVerificationCode implements models.UserUsecase.
func (uc *UserUsecase) ResendVerificationCode(ctx context.Context, email string) *fiber.Error {
	// get user based on email param
	user, err := uc.userRepo.FindUserByEmail(ctx, email)
	if err != nil {
		return err
	}

	err = uc.userRepo.ResendVerificationCode(user)
	if err != nil {
		return err
	}

	return nil
}

// VerificationEmail implements models.UserUsecase.
func (uc *UserUsecase) VerificationEmail(ctx context.Context, code string) *fiber.Error {
	err := uc.userRepo.VerificationEmail(ctx, code)
	if err != nil {
		return err
	}
	return nil
}

// RefreshToken implements models.UserUsecase.
func (uc *UserUsecase) RefreshToken(ctx context.Context, payload models.RefreshTokenInput) (models.Token, *fiber.Error) {
	data, err := uc.userRepo.RefreshToken(ctx, payload)
	if err != nil {
		return data, err
	}
	return data, nil
}

// Login implements models.UserUsecase.
func (uc *UserUsecase) Login(ctx context.Context, payload models.LoginInput) (models.Token, *fiber.Error) {
	// check email or username exists
	user, err := uc.userRepo.FindUserByIdentity(ctx, payload.Email)
	if err != nil {
		return models.Token{}, err
	}

	if !user.Verified {
		return models.Token{}, fiber.NewError(400, "Your account is not active yet, please verify your email.")
	}

	if err := user.ValidatePassword(payload.Password); err != nil {
		return models.Token{}, fiber.NewError(400, "Invalid Email or Password.")
	}

	data, err := uc.userRepo.Login(ctx, user)
	if err != nil {
		return models.Token{}, err
	}
	return data, nil
}

// Register implements models.UserUsecase.
func (uc *UserUsecase) Register(ctx context.Context, payload models.RegisterInput) (models.User, *fiber.Error) {

	// cek email of user
	if err := uc.userRepo.EmailExists(payload.Email); err != nil {
		return models.User{}, err
	}

	// cek username of user
	if err := uc.userRepo.UsernameExists(payload.Username); err != nil {
		return models.User{}, err
	}

	user := models.User{
		Username:  payload.Username,
		Password:  payload.Password,
		Email:     payload.Email,
		FirstName: payload.FirstName,
		LastName:  payload.LastName,
		UserProfile: models.UserProfile{
			Phone:    payload.Phone,
			StatusID: 3, // pending
		},
	}

	user, err := uc.userRepo.Register(ctx, user)
	if err != nil {
		return user, err
	}
	return user, nil
}

// NewAuthEntity will create new an articleUsecase object representation of domain.ArticleUsecase interface
func NewUserUsecase(userRepo models.UserRepository) models.UserUsecase {
	return &UserUsecase{
		userRepo: userRepo,
	}
}
