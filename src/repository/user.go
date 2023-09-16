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
	fmt.Println("DeleteAuthRedis: ", deleted)
	return deleted, nil

}

// DeleteToken implements models.UserRepository.
func (r *UserRepository) DeleteToken(authD *models.AccessDetails) error {
	//get the refresh uuid
	refreshUuid := fmt.Sprintf("%s++%d", authD.TokenUuid, authD.UserID)

	//delete access token
	fmt.Println("DeleteToken.AccessUuid : ", authD.TokenUuid)
	deletedAt, err := r.DeleteAuthRedis(authD.TokenUuid)
	if err != nil || deletedAt == 0 {
		return err
	}

	//delete refresh token
	fmt.Println("DeleteToken.refreshUuid : ", refreshUuid)
	deletedRt, err := r.DeleteAuthRedis(refreshUuid)
	if err != nil || deletedRt == 0 {
		return err
	}

	//When the record is deleted, the return value is 1
	if deletedAt != 1 || deletedRt != 1 {
		return errors.New("something went wrong")
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
func (r *UserRepository) ResendVerificationCode(user models.User) error {
	if user.Verified {
		return errors.New("User already verified. You can login now.")
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
func (r *UserRepository) VerificationEmail(ctx context.Context, code string) error {
	verification_code := helpers.Encode(code)

	var user models.User
	result := r.DB.First(&user, "verification_code = ?", verification_code)
	if result.Error != nil {
		return errors.New("Invalid verification code or user doesn't exists.")
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
	result := r.DB.Where("email = ?", strings.ToLower(email)).Find(&user)
	if result.Error != nil {
		return errors.New(result.Error.Error())
	}

	if result.RowsAffected != 0 {
		return errors.New("Email already registered, please use another one!")
	}
	return nil
}

// UsernameExists implements models.UserRepository.
func (r *UserRepository) UsernameExists(username string) error {
	var user models.User
	result := r.DB.Where("username = ?", strings.ToLower(username)).Find(&user)
	if result.Error != nil {
		return errors.New(result.Error.Error())
	}

	if result.RowsAffected != 0 {
		return errors.New("Username already registered, please use another one!")
	}
	return nil
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

// FindUserByEmail implements models.UserRepository.
func (r *UserRepository) FindUserByEmail(ctx context.Context, email string) (models.User, error) {
	var user models.User
	result := r.DB.Where("email=?", strings.ToLower(email)).Find(&user)
	if result.Error != nil {
		return user, errors.New(result.Error.Error())
	}

	if result.RowsAffected == 0 {
		return user, errors.New("Invalid email or account doesn't exists.")
	}

	return user, nil
}

// FindUserByIdentity implements models.UserRepository.
func (r *UserRepository) FindUserByIdentity(ctx context.Context, identity string) (models.User, error) {
	var user models.User
	result := r.DB.Where("username=?", strings.ToLower(identity)).Or("email=?", strings.ToLower(identity)).Find(&user)
	if result.Error != nil {
		return user, errors.New(result.Error.Error())
	}

	if result.RowsAffected == 0 {
		return user, errors.New("Invalid Email or Account doesn't exists.")
	}

	return user, nil
}

// RefreshToken implements models.UserRepository.
func (r *UserRepository) RefreshToken(ctx context.Context, payload models.RefreshTokenInput) (models.Token, error) {
	var token models.Token
	config, _ := configs.LoadConfig(".")

	// validate refrefresh_token
	tokenClaims, err := helpers.ValidateToken(payload.RefreshToken, config.RefreshTokenPublicKey)
	if err != nil {
		fmt.Println("tokenClaims")
		return token, err
	}
	fmt.Println("tokenClaims PASS")

	refreshUuid := tokenClaims.TokenUuid

	// check refreshUuid in Redis
	ctxTodo := context.TODO()
	_, err = configs.RedisClient.Get(ctxTodo, refreshUuid).Result()
	if err == redis.Nil {
		return token, errors.New("Token is invalid or session has expired")
	}

	// fmt.Println("refreshUuid: ", refreshUuid)
	// //Delete the previous Refresh Token
	// deletedRt, err := r.DeleteAuthRedis(refreshUuid)
	// if err != nil || deletedRt == 0 {
	// 	fmt.Println("deleted")
	// 	return token, err
	// }
	// fmt.Println("DeleteAuthRedis PASS")

	var user models.User
	err = r.DB.First(&user, "id = ?", tokenClaims.UserID).Error

	if err == gorm.ErrRecordNotFound {
		return token, fmt.Errorf("the user belonging to this token no logger exists")
	}

	// generate new tokens
	token, err = r.GeneratePairToken(user.ID)
	if err != nil {
		fmt.Println("GeneratePairToken")
		return token, err
	}
	fmt.Println("GeneratePairToken PASS")
	return token, nil
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
func (r *UserRepository) Login(ctx context.Context, user models.User) (models.Token, error) {
	token, err := r.GeneratePairToken(user.ID)
	if err != nil {
		return token, err
	}
	fmt.Println("Login PASS")
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

	// Send verification email
	r.SendVerificationEmail(user, code)

	return user, nil
}

// NewMysqlArticleRepository will create an object that represent the article.Repository interface
func NewUserRepository(Conn *gorm.DB) models.UserRepository {
	return &UserRepository{Conn}
}
