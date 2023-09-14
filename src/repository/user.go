package repository

import (
	"context"
	"fmt"
	initializers "myapp/pkg/configs"
	"myapp/pkg/helpers"
	"myapp/src/models"
	"strings"

	"gorm.io/gorm"
)

type UserRepository struct {
	Conn *gorm.DB
}

// EmailExists implements models.UserRepository.
func (*UserRepository) EmailExists(email string) error {
	panic("unimplemented")
}

// UsernameExists implements models.UserRepository.
func (*UserRepository) UsernameExists(username string) error {
	panic("unimplemented")
}

// RefreshToken implements models.UserRepository.
func (r *UserRepository) RefreshToken(ctx context.Context, payload models.RefreshTokenInput) (models.Token, error) {
	var token models.Token

	refresh_token := payload.RefreshToken

	config, _ := initializers.LoadConfig(".")

	// validate refrefresh_token
	tokenClaims, err := helpers.ValidateToken(refresh_token, config.RefreshTokenPublicKey)
	if err != nil {
		return token, err
	}

	var user models.User
	err = initializers.DB.Preload("UserProfile.Status").First(&user, "id = ?", tokenClaims.UserID).Error

	if err == gorm.ErrRecordNotFound {
		return token, fmt.Errorf("the user belonging to this token no logger exists")
	}

	// generate new tokens
	accessTokenDetails, err := helpers.CreateToken(fmt.Sprint(user.ID), config.AccessTokenExpiresIn, config.AccessTokenPrivateKey)
	if err != nil {
		return token, err
	}

	refreshTokenDetails, err := helpers.CreateToken(fmt.Sprint(user.ID), config.RefreshTokenExpiresIn, config.RefreshTokenPrivateKey)
	if err != nil {
		return token, err
	}

	token.AccessToken = *accessTokenDetails.Token
	token.RefreshToken = *refreshTokenDetails.Token
	token.ExpiresIn = *accessTokenDetails.ExpiresIn
	token.TokenType = accessTokenDetails.TokenType

	return token, nil
}

// Create implements models.UserRepository.
func (r *UserRepository) Create(ctx context.Context, md models.User) error {
	panic("unimplemented")
}

// FindUserById implements models.UserRepository.
func (r *UserRepository) FindUserById(ctx context.Context, id uint) (models.User, error) {
	var user models.User
	err := r.Conn.Preload("UserProfile.Status").Where("ID=?", id).Find(&user).Error
	if err != nil {
		return models.User{}, err
	}
	return user, nil
}

// FindUserByUsername implements models.UserRepository.
func (r *UserRepository) FindUserByUsername(ctx context.Context, username string) (models.User, error) {
	var user models.User
	err := r.Conn.Where("username=?", strings.ToLower(username)).Or("email=?", strings.ToLower(username)).Find(&user).Error
	if err != nil {
		return user, err
	}
	return user, nil
}

// Login implements models.UserRepository.
func (r *UserRepository) Login(ctx context.Context, payload models.LoginInput) (models.Token, error) {
	var token models.Token

	user, err := r.FindUserByUsername(ctx, payload.Username)
	if err != nil {
		fmt.Println(err)
		return token, err
	}

	err = user.ValidatePassword(payload.Password)
	if err != nil {
		return token, fmt.Errorf("invalid email/username or password")
	}

	// token, err = helpers.GenerateJWT(user)
	// if err != nil {
	// 	return token, err
	// }

	config, _ := initializers.LoadConfig(".")

	accessTokenDetails, err := helpers.CreateToken(fmt.Sprint(user.ID), config.AccessTokenExpiresIn, config.AccessTokenPrivateKey)
	if err != nil {
		return token, err
	}

	refreshTokenDetails, err := helpers.CreateToken(fmt.Sprint(user.ID), config.RefreshTokenExpiresIn, config.RefreshTokenPrivateKey)
	if err != nil {
		return token, err
	}

	token.AccessToken = *accessTokenDetails.Token
	token.RefreshToken = *refreshTokenDetails.Token
	token.ExpiresIn = *accessTokenDetails.ExpiresIn
	token.TokenType = accessTokenDetails.TokenType

	return token, nil
}

// Register implements models.UserRepository.
func (r *UserRepository) Register(ctx context.Context, payload models.RegisterInput) (models.User, error) {

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

	err := r.Conn.Create(&user).Error
	if err != nil {
		if strings.Contains(err.Error(), "duplicate key value violates unique") {
			// fmt.Println(err)
			return user, fmt.Errorf("user with that email/username already exists")
		}

		return user, err
	}

	// Send Email if register successfully
	// emailData := utils.EmailData{
	// 	URL:       config.Origin + "/api/v1/verify-email/" + code,
	// 	FirstName: firstName,
	// 	Subject:   "Your account verification",
	// }

	// // TODO: send email with celery task
	// utils.SendEmail(user, &emailData, "verificationCode.html")

	return user, nil
}

// NewMysqlArticleRepository will create an object that represent the article.Repository interface
func NewUserRepository(Conn *gorm.DB) models.UserRepository {
	return &UserRepository{Conn}
}
