package usecase

import (
	"context"
	"myapp/src/models"
)

type UserUsecase struct {
	Repo models.UserRepository
}

// UsernameExists implements models.UserUsecase.
func (*UserUsecase) UsernameExists(username string) error {
	panic("unimplemented")
}

// EmailExists implements models.UserUsecase.
func (uc *UserUsecase) EmailExists(ucmail string) error {
	panic("unimplemented")
}

// RefreshToken implements models.UserUsecase.
func (uc *UserUsecase) RefreshToken(ctx context.Context, payload models.RefreshTokenInput) (models.Token, error) {
	data, err := uc.Repo.RefreshToken(ctx, payload)
	if err != nil {
		return data, err
	}
	return data, nil
}

// Create implements models.UserUsecase.
func (uc *UserUsecase) Create(ctx context.Context, md models.User) error {
	panic("unimplemented")
}

// FindUserById implements models.UserUsecase.
func (uc *UserUsecase) FindUserById(ctx context.Context, id uint) (models.User, error) {
	panic("unimplemented")
}

// FindUserByUsername implements models.UserUsecase.
func (uc *UserUsecase) FindUserByUsername(ctx context.Context, username string) (models.User, error) {
	panic("unimplemented")
}

// Login implements models.UserUsecase.
func (uc *UserUsecase) Login(ctx context.Context, payload models.LoginInput) (models.Token, error) {
	data, err := uc.Repo.Login(ctx, payload)
	if err != nil {
		return data, err
	}
	return data, nil
}

// Register implements models.UserUsecase.
func (uc *UserUsecase) Register(ctx context.Context, payload models.RegisterInput) (models.User, error) {

	// cek email of user
	if err := uc.EmailExists(payload.Email); err != nil {
		return models.User{}, err
	}

	// cek username of user
	if err := uc.UsernameExists(payload.Username); err != nil {
		return models.User{}, err
	}

	data, err := uc.Repo.Register(ctx, payload)
	if err != nil {
		return data, err
	}
	return data, nil
}

// NewAuthEntity will create new an articleUsecase object representation of domain.ArticleUsecase interface
func NewUserUsecase(repo models.UserRepository) models.UserUsecase {
	return &UserUsecase{
		Repo: repo,
	}
}
