package repository

import (
	"context"
	"errors"
	"fmt"
	"myapp/pkg/configs"
	"myapp/pkg/helpers"
	"myapp/src/models"
	"strings"
	"time"

	"github.com/thanhpk/randstr"
	"gorm.io/gorm"
)

type UserRepository struct {
	DB *gorm.DB
}

// VerificationEmail implements models.UserRepository.
func (r *UserRepository) VerificationEmail(ctx context.Context, code string) error {
	verification_code := helpers.Encode(code)

	var user models.User
	result := r.DB.First(&user, "verification_code = ?", verification_code)
	if result.Error != nil {
		return errors.New("Invalid verification code or user doesn't exists.")
	}

	if user.Verified {
		return errors.New("User already verified.")
	}

	now := time.Now()
	user.VerificationCode = ""
	user.Verified = true
	user.VerifiedAt = &now

	// run update userprofile.statud_id to 1 (Active)
	r.DB.Model(&models.UserProfile{}).Where("user_id = ?", user.ID).
		Update("status_id", 1)

	err := r.DB.Save(&user).Error
	if err != nil {
		return err
	}

	return nil
}

// EmailExists implements models.UserRepository.
func (r *UserRepository) EmailExists(email string) error {
	var user models.User
	err := r.DB.Where("email = ?", strings.ToLower(email)).Find(&user).Error
	if err != nil {
		return errors.New("Email already registered, please use another one!")
	}
	return nil
}

// UsernameExists implements models.UserRepository.
func (r *UserRepository) UsernameExists(username string) error {
	var user models.User
	err := r.DB.Where("username = ?", strings.ToLower(username)).Find(&user).Error
	if err != nil {
		return errors.New("Username already registered, please use another one!")
	}
	return nil
}

// RefreshToken implements models.UserRepository.
func (r *UserRepository) RefreshToken(ctx context.Context, payload models.RefreshTokenInput) (models.Token, error) {
	var token models.Token

	refresh_token := payload.RefreshToken

	config, _ := configs.LoadConfig(".")

	// validate refrefresh_token
	tokenClaims, err := helpers.ValidateToken(refresh_token, config.RefreshTokenPublicKey)
	if err != nil {
		return token, err
	}

	var user models.User
	err = r.DB.Preload("UserProfile.Status").First(&user, "id = ?", tokenClaims.UserID).Error

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
	err := r.DB.Preload("UserProfile.Status").Where("ID=?", id).Find(&user).Error
	if err != nil {
		return models.User{}, err
	}
	return user, nil
}

// FindUserByUsername implements models.UserRepository.
func (r *UserRepository) FindUserByUsername(ctx context.Context, username string) (models.User, error) {
	var user models.User
	err := r.DB.Where("username=?", strings.ToLower(username)).Or("email=?", strings.ToLower(username)).Find(&user).Error
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

	config, _ := configs.LoadConfig(".")

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
func (r *UserRepository) Register(ctx context.Context, user models.User) (models.User, error) {

	// Generate Verification Code
	code := randstr.String(64)

	verification_code := helpers.Encode(code)

	// fill User verification code
	user.VerificationCode = verification_code

	err := r.DB.Create(&user).Error
	if err != nil {
		if strings.Contains(err.Error(), "duplicate key value violates unique") {
			// fmt.Println(err)
			return user, fmt.Errorf("user with that email/username already exists")
		}

		return user, err
	}

	var accountName = user.FirstName

	if accountName != "" {
		if strings.Contains(accountName, " ") {
			accountName = strings.Split(accountName, " ")[1]
		}
	} else {
		accountName = user.Email
	}

	siteData, _ := configs.GetSiteData(".")
	// Send Email if register successfully
	emailData := helpers.EmailData{
		URL:          siteData.ClientOrigin + "/verify-email/" + code,
		FirstName:    accountName,
		Subject:      "Your account verification",
		TypeOfAction: "Register",
		SiteData:     siteData,
	}

	// send email with goroutine
	go helpers.SendEmail(user, &emailData, "verificationCode.html")

	return user, nil
}

// NewMysqlArticleRepository will create an object that represent the article.Repository interface
func NewUserRepository(Conn *gorm.DB) models.UserRepository {
	return &UserRepository{Conn}
}
