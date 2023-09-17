package repository

import (
	"context"
	"fmt"
	"myapp/pkg/configs"
	"myapp/pkg/helpers"
	"myapp/src/models"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/redis/go-redis/v9"
	"github.com/thanhpk/randstr"
	"gorm.io/gorm"
)

type UserRepository struct {
	DB *gorm.DB
}

// DeleteAuth implements models.UserRepository.
func (*UserRepository) DeleteAuthRedis(givenUuid string) (int64, error) {
	ctx := context.TODO()
	deleted, err := configs.RedisClient.Del(ctx, givenUuid).Result()
	if err != nil {
		return 0, err
	}
	// fmt.Println("DeleteAuthRedis: ", deleted)
	return deleted, nil

}

// DeleteToken implements models.UserRepository.
func (r *UserRepository) DeleteToken(authD *models.AccessDetails) *fiber.Error {
	//get the refresh uuid
	refreshUuid := fmt.Sprintf("%s++%d", authD.TokenUuid, authD.UserID)

	//delete access token
	deletedAt, err := r.DeleteAuthRedis(authD.TokenUuid)
	if err != nil || deletedAt == 0 {
		return fiber.NewError(500, err.Error())
	}

	//delete refresh token
	deletedRt, err := r.DeleteAuthRedis(refreshUuid)
	if err != nil || deletedRt == 0 {
		return fiber.NewError(500, err.Error())
	}

	//When the record is deleted, the return value is 1
	if deletedAt != 1 || deletedRt != 1 {
		return fiber.NewError(500, "something went wrong")
	}
	return nil
}

// SendVerificationEmail implements models.UserRepository.
func (*UserRepository) SendVerificationEmail(user models.User, code string) error {
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
	go helpers.SendEmail(user, &emailData, "verificationCode")

	return nil
}

// ResendVerificationCode implements models.UserRepository.
func (r *UserRepository) ResendVerificationCode(user models.User) *fiber.Error {
	if user.Verified {
		return fiber.NewError(422, "User already verified. You can login now.")
	}

	// Generate Verification Code
	code := randstr.String(64)

	verification_code := helpers.Encode(code)

	// Update User in Database
	user.VerificationCode = verification_code
	r.DB.Save(&user)

	// Send verification email
	r.SendVerificationEmail(user, code)

	return nil
}

// VerificationEmail implements models.UserRepository.
func (r *UserRepository) VerificationEmail(ctx context.Context, code string) *fiber.Error {
	verification_code := helpers.Encode(code)

	var user models.User
	result := r.DB.First(&user, "verification_code = ?", verification_code)
	if result.Error != nil {
		return fiber.NewError(404, "Invalid verification code or user doesn't exists.")
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
		return fiber.NewError(404, err.Error())
	}

	return nil
}

// EmailExists implements models.UserRepository.
func (r *UserRepository) EmailExists(email string) *fiber.Error {
	var user models.User
	result := r.DB.Where("email = ?", strings.ToLower(email)).Find(&user)
	if result.Error != nil {
		return fiber.NewError(500, result.Error.Error())
	}

	if result.RowsAffected != 0 {
		return fiber.NewError(422, "Email already registered, please use another one!")
	}
	return nil
}

// UsernameExists implements models.UserRepository.
func (r *UserRepository) UsernameExists(username string) *fiber.Error {
	var user models.User
	result := r.DB.Where("username = ?", strings.ToLower(username)).Find(&user)
	if result.Error != nil {
		return fiber.NewError(500, result.Error.Error())
	}

	if result.RowsAffected != 0 {
		return fiber.NewError(422, "Username already registered, please use another one!")
	}
	return nil
}

// Create implements models.UserRepository.
func (r *UserRepository) Create(ctx context.Context, md models.User) *fiber.Error {
	panic("unimplemented")
}

// FindUserById implements models.UserRepository.
func (r *UserRepository) FindUserById(ctx context.Context, id uint) (models.User, *fiber.Error) {
	var user models.User
	result := r.DB.Preload("UserProfile.Status").Where("ID=?", id).Find(&user)
	if result.Error != nil {
		return user, fiber.NewError(500, result.Error.Error())
	}

	if result.RowsAffected == 0 {
		return user, fiber.NewError(422, "Invalid email or account doesn't exists.")
	}
	return user, nil
}

// FindUserByEmail implements models.UserRepository.
func (r *UserRepository) FindUserByEmail(ctx context.Context, email string) (models.User, *fiber.Error) {
	var user models.User
	result := r.DB.Where("email=?", strings.ToLower(email)).Find(&user)
	if result.Error != nil {
		return user, fiber.NewError(500, result.Error.Error())
	}

	if result.RowsAffected == 0 {
		return user, fiber.NewError(422, "Invalid email or account doesn't exists.")
	}

	return user, nil
}

// FindUserByIdentity implements models.UserRepository.
func (r *UserRepository) FindUserByIdentity(ctx context.Context, identity string) (models.User, *fiber.Error) {
	var user models.User
	result := r.DB.Where("username=?", strings.ToLower(identity)).Or("email=?", strings.ToLower(identity)).Find(&user)
	if result.Error != nil {
		return user, fiber.NewError(500, result.Error.Error())
	}

	if result.RowsAffected == 0 {
		return user, fiber.NewError(422, "Invalid Email or Account doesn't exists.")
	}

	return user, nil
}

// RefreshToken implements models.UserRepository.
func (r *UserRepository) RefreshToken(ctx context.Context, payload models.RefreshTokenInput) (models.Token, *fiber.Error) {
	var token models.Token
	config, _ := configs.LoadConfig(".")

	// validate refrefresh_token
	tokenClaims, err := helpers.ValidateToken(payload.RefreshToken, config.RefreshTokenPublicKey)
	if err != nil {
		return token, fiber.ErrUnauthorized
	}

	refreshUuid := tokenClaims.TokenUuid

	// check refreshUuid in Redis
	ctxTodo := context.TODO()
	_, err = configs.RedisClient.Get(ctxTodo, refreshUuid).Result()
	if err == redis.Nil {
		return token, fiber.NewError(401, "Token is invalid or session has expired")
	}

	var user models.User
	err = r.DB.First(&user, "id = ?", tokenClaims.UserID).Error

	if err == gorm.ErrRecordNotFound {
		return token, fiber.NewError(404, "the user belonging to this token no logger exists")
	}

	// generate new tokens
	token, err = r.GeneratePairToken(user.ID)
	if err != nil {
		return token, fiber.NewError(404, err.Error())
	}
	return token, fiber.NewError(404, err.Error())
}

// GeneratePairToken implements models.UserRepository.
func (*UserRepository) GeneratePairToken(userID uint) (models.Token, error) {
	var token models.Token

	td, err := helpers.CreateToken(userID)
	if err != nil {
		return token, err
	}

	// Save Token in Redis
	ctxTodo := context.TODO()
	now := time.Now()

	fmt.Println("td.AccessUuid : ", td.AccessUuid)
	errAccess := configs.RedisClient.Set(ctxTodo, td.AccessUuid, userID, time.Unix(td.AtExpires, 0).Sub(now)).Err()
	if errAccess != nil {
		return token, errAccess
	}

	fmt.Println("td.RefreshUuid : ", td.RefreshUuid)
	errRefresh := configs.RedisClient.Set(ctxTodo, td.RefreshUuid, userID, time.Unix(td.RtExpires, 0).Sub(now)).Err()
	if errRefresh != nil {
		return token, errRefresh
	}

	token.AccessToken = td.AccessToken
	token.RefreshToken = td.RefreshToken
	token.ExpiresIn = td.AtExpires
	token.TokenType = td.TokenType

	return token, nil
}

// Login implements models.UserRepository.
func (r *UserRepository) Login(ctx context.Context, user models.User) (models.Token, *fiber.Error) {
	token, err := r.GeneratePairToken(user.ID)
	if err != nil {
		return token, fiber.NewError(500, err.Error())
	}
	return token, nil
}

// Register implements models.UserRepository.
func (r *UserRepository) Register(ctx context.Context, user models.User) (models.User, *fiber.Error) {
	// Generate Verification Code
	code := randstr.String(64)

	verification_code := helpers.Encode(code)

	// fill User verification code
	user.VerificationCode = verification_code

	err := r.DB.Create(&user).Error
	if err != nil {
		if strings.Contains(err.Error(), "duplicate key value violates unique") {
			// fmt.Println(err)
			return user, fiber.NewError(422, "user with that email/username already exists")
		}

		return user, fiber.NewError(500, err.Error())
	}

	// Send verification email
	r.SendVerificationEmail(user, code)

	return user, nil
}

// NewMysqlArticleRepository will create an object that represent the article.Repository interface
func NewUserRepository(Conn *gorm.DB) models.UserRepository {
	return &UserRepository{Conn}
}
