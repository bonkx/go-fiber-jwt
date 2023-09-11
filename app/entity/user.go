package entity

import (
	"context"
	"myapp/models"
)

type UserEntity struct {
	Repo models.UserRepository
}

// RefreshToken implements models.UserEntity.
func (e *UserEntity) RefreshToken(ctx context.Context, payload models.RefreshTokenInput) (models.Token, error) {
	data, err := e.Repo.RefreshToken(ctx, payload)
	if err != nil {
		return data, err
	}
	return data, nil
}

// Create implements models.UserEntity.
func (e *UserEntity) Create(ctx context.Context, md models.User) error {
	panic("unimplemented")
}

// FindUserById implements models.UserEntity.
func (e *UserEntity) FindUserById(ctx context.Context, id uint) (models.User, error) {
	panic("unimplemented")
}

// FindUserByUsername implements models.UserEntity.
func (e *UserEntity) FindUserByUsername(ctx context.Context, username string) (models.User, error) {
	panic("unimplemented")
}

// Login implements models.UserEntity.
func (e *UserEntity) Login(ctx context.Context, payload models.LoginInput) (models.Token, error) {
	data, err := e.Repo.Login(ctx, payload)
	if err != nil {
		return data, err
	}
	return data, nil
}

// Register implements models.UserEntity.
func (e *UserEntity) Register(ctx context.Context, payload models.RegisterInput) (models.User, error) {
	data, err := e.Repo.Register(ctx, payload)
	if err != nil {
		return data, err
	}
	return data, nil
}

// NewAuthEntity will create new an articleUsecase object representation of domain.ArticleUsecase interface
func NewUserEntity(repo models.UserRepository) models.UserEntity {
	return &UserEntity{
		Repo: repo,
	}
}
