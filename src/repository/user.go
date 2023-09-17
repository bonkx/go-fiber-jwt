package repository

import (
	"context"
	"fmt"
	"myapp/pkg/configs"
	"myapp/pkg/helpers"
	"myapp/pkg/utils"
	"myapp/src/models"
	"strconv"
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

// NewUserRepository will create an object that represent the models.UserRepository interface
func NewUserRepository(Conn *gorm.DB) models.UserRepository {
	return &UserRepository{Conn}
}

// Update implements models.UserRepository.
func (r *UserRepository) Update(user models.User) (models.User, *fiber.Error) {

	// update user and profile data as well
	err := r.DB.Session(&gorm.Session{FullSaveAssociations: true}).Save(&user).Error
	if err != nil {
		return user, fiber.NewError(500, err.Error())
	}

	return user, nil
}

// ChangePassword implements models.UserRepository.
func (r *UserRepository) ChangePassword(user models.User) *fiber.Error {
	err := r.DB.Save(&user).Error
	if err != nil {
		return fiber.NewError(500, err.Error())
	}

	return nil
}

func (r *UserRepository) deleteAllOTPRequestByEmail(email string) {
	fmt.Println("deleteAllOTPRequestByEmail ======================")
	r.DB.Unscoped().Where("email=?", email).Delete(&models.OTPRequest{})
}

// ResetPassword implements models.UserRepository.
func (r *UserRepository) ResetPassword(user models.User) *fiber.Error {

	err := r.DB.Save(&user).Error
	if err != nil {
		return fiber.NewError(500, err.Error())
	}

	// goroutine - delete all otpotpR by email
	go r.deleteAllOTPRequestByEmail(user.Email)

	return nil
}

// VerifyOTP implements models.UserRepository.
func (r *UserRepository) VerifyOTP(otpR models.OTPRequest) (string, *fiber.Error) {
	// generate random number (20)
	randomN, err := utils.GenerateRandomNumber(20)
	if err != nil {
		return "", fiber.NewError(500, err.Error())
	}

	refNo := strconv.Itoa(randomN)
	// update refNo otpR
	err = r.DB.Model(&otpR).Update("reference_no", refNo).Error
	if err != nil {
		return "", fiber.NewError(500, err.Error())
	}

	return refNo, nil
}

// FindReferenceOTPRequest implements models.UserRepository.
func (r *UserRepository) FindReferenceOTPRequest(refNo string) (models.OTPRequest, *fiber.Error) {
	otpR := models.OTPRequest{}

	result := r.DB.First(&otpR, "reference_no=?", refNo)
	if result.RowsAffected == 0 {
		return otpR, fiber.NewError(422, "ReferenceNo doesn't exists.")
	}
	return otpR, nil
}

// FindOTPRequest implements models.UserRepository.
func (r *UserRepository) FindOTPRequest(otp string) (models.OTPRequest, *fiber.Error) {
	otpR := models.OTPRequest{}

	result := r.DB.First(&otpR, "otp=?", otp)
	if result.RowsAffected == 0 {
		return otpR, fiber.NewError(422, "OTP code doesn't exists.")
	}
	return otpR, nil
}

// RequestOTPEmail implements models.UserRepository.
func (r *UserRepository) RequestOTPEmail(user models.User) *fiber.Error {
	otpR := models.OTPRequest{Email: user.Email}

	// create OTP Request
	err := r.DB.Create(&otpR).Error
	if err != nil {
		return fiber.NewError(500, err.Error())
	}

	siteData, _ := configs.GetSiteData(".")
	emailData := helpers.EmailData{
		URL:          otpR.Otp,
		FirstName:    user.Email,
		Subject:      "Your OTP code",
		TypeOfAction: "Forgot Password",
		SiteData:     siteData,
	}

	// send email with goroutine
	go helpers.SendEmail(user, &emailData, "otp_code.html")

	return nil
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
	// Send Email
	emailData := helpers.EmailData{
		URL:          siteData.ClientOrigin + "/verify-email/" + code,
		FirstName:    accountName,
		Subject:      "Your account verification",
		TypeOfAction: "Register",
		SiteData:     siteData,
	}

	// send email with goroutine
	go helpers.SendEmail(user, &emailData, "verification_code.html")

	return nil
}

// ResendVerificationCode implements models.UserRepository.
func (r *UserRepository) ResendVerificationCode(user models.User) *fiber.Error {
	if user.Verified {
		return fiber.NewError(422, "User already verified. You can login now.")
	}

	// Generate Verification Code
	code := randstr.String(64)

	verification_code := utils.Encode(code)

	// Update User in Database
	user.VerificationCode = verification_code
	r.DB.Save(&user)

	// Send verification email
	r.SendVerificationEmail(user, code)

	return nil
}

// VerificationEmail implements models.UserRepository.
func (r *UserRepository) VerificationEmail(code string) *fiber.Error {
	verification_code := utils.Encode(code)

	var user models.User
	result := r.DB.First(&user, "verification_code=?", verification_code)
	if result.Error != nil {
		return fiber.NewError(404, "Invalid verification code or user doesn't exists.")
	}

	now := time.Now()
	user.VerificationCode = ""
	user.Verified = true
	user.VerifiedAt = &now

	// run update userprofile.statud_id to 1 (Active)
	r.DB.Model(&models.UserProfile{}).Where("user_id=?", user.ID).
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
	result := r.DB.Where("email=?", strings.ToLower(email)).Find(&user)
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
	result := r.DB.Where("username=?", strings.ToLower(username)).Find(&user)
	if result.Error != nil {
		return fiber.NewError(500, result.Error.Error())
	}

	if result.RowsAffected != 0 {
		return fiber.NewError(422, "Username already registered, please use another one!")
	}
	return nil
}

// Create implements models.UserRepository.
func (r *UserRepository) Create(md models.User) *fiber.Error {
	panic("unimplemented")
}

// FindUserById implements models.UserRepository.
func (r *UserRepository) FindUserById(id uint) (models.User, *fiber.Error) {
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
func (r *UserRepository) FindUserByEmail(email string) (models.User, *fiber.Error) {
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
func (r *UserRepository) FindUserByIdentity(identity string) (models.User, *fiber.Error) {
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
func (r *UserRepository) RefreshToken(payload models.RefreshTokenInput) (models.Token, *fiber.Error) {
	token := models.Token{}
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
	err = r.DB.First(&user, "id=?", tokenClaims.UserID).Error

	if err == gorm.ErrRecordNotFound {
		return token, fiber.NewError(404, "the user belonging to this token no logger exists")
	}

	// generate new tokens
	token, err = r.GeneratePairToken(tokenClaims.UserID)
	if err != nil {
		return token, fiber.NewError(404, err.Error())
	}

	return token, nil
}

// GeneratePairToken implements models.UserRepository.
func (r *UserRepository) GeneratePairToken(userID uint) (models.Token, error) {
	token := models.Token{}

	td, err := helpers.CreateToken(userID)
	if err != nil {
		return token, err
	}

	ctxTodo := context.TODO()
	now := time.Now()

	// fmt.Println("td.AccessUuid : ", td.AccessUuid)
	// Save Access Token in Redis
	errAccess := configs.RedisClient.Set(ctxTodo, td.AccessUuid, userID, time.Unix(td.AtExpires, 0).Sub(now)).Err()
	if errAccess != nil {
		return token, errAccess
	}

	// fmt.Println("td.RefreshUuid : ", td.RefreshUuid)
	// Save Access Refresh in Redis
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
func (r *UserRepository) Login(user models.User) (models.Token, *fiber.Error) {
	token, err := r.GeneratePairToken(user.ID)
	if err != nil {
		return token, fiber.NewError(500, err.Error())
	}
	return token, nil
}

// Register implements models.UserRepository.
func (r *UserRepository) Register(user models.User) (models.User, *fiber.Error) {
	// Generate Verification Code
	code := randstr.String(64)

	verification_code := utils.Encode(code)

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
