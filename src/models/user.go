package models

import (
	"context"
	"fmt"
	"html"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

// User defines the user in db
// User struct is used to store user information in the database
type User struct {
	gorm.Model
	FirstName        string     `json:"first_name"`
	LastName         string     `json:"last_name"`
	Username         string     `json:"username" gorm:"not null;unique"`
	Email            string     `json:"email" gorm:"unique"`
	Password         string     `json:"-"`
	Verified         bool       `json:"verified" gorm:"not null;default:false"`
	IsSuperuser      bool       `json:"is_superuser" gorm:"default:false"`
	IsStaff          bool       `json:"is_staff" gorm:"default:false"`
	LastLogin        *time.Time `json:"last_login"`
	VerificationCode string     `json:"verification_code"`
	VerifiedAt       *time.Time `json:"verified_at"`

	UserProfile UserProfile `gorm:"foreignkey:UserID"`
}

func (user *User) BeforeCreate(*gorm.DB) error {
	passwordHash, err := user.HashPassword(user.Password)
	if err != nil {
		return err
	}
	user.Password = passwordHash
	user.Username = html.EscapeString(strings.TrimSpace(user.Username))
	user.Email = strings.ToLower(user.Email)

	return nil
}

func (user *User) AfterCreate(tx *gorm.DB) (err error) {
	if user.Username == "admin" {
		tx.Model(user).Updates(User{IsSuperuser: true, IsStaff: true})
		tx.Model(user.UserProfile).Update("role", "admin")
	}
	tx.Model(user.UserProfile).Update("role", "user")
	return
}

func (user *User) ValidatePassword(password string) error {
	return bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
}

func (user *User) HashPassword(password string) (string, error) {
	passwordHash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", fmt.Errorf("could not hash password %w", err)
	}
	return string(passwordHash), nil
}

type UserUsecase interface {
	Register(ctx context.Context, payload RegisterInput) (User, *fiber.Error)
	Login(ctx context.Context, payload LoginInput) (Token, *fiber.Error)
	RefreshToken(ctx context.Context, payload RefreshTokenInput) (Token, *fiber.Error)
	VerificationEmail(ctx context.Context, code string) *fiber.Error
	ResendVerificationCode(ctx context.Context, email string) *fiber.Error
	Logout(authD *AccessDetails) *fiber.Error

	ForgotPassword(ctx context.Context, payload EmailInput) *fiber.Error
	ForgotPasswordOTP(ctx context.Context, payload OTPInput) (string, *fiber.Error)
	ResetPassword(ctx context.Context, payload ResetPasswordInput) *fiber.Error
	ChangePassword(ctx context.Context, md User, payload ChangePasswordInput) *fiber.Error
	Update(c *fiber.Ctx, payload UpdateProfileInput) (User, *fiber.Error)
	// Delete(ctx context.Context, md User) *fiber.Error
	UploadPhotoProfile(c *fiber.Ctx, md User) *fiber.Error
	RequestDeleteAccount(c *fiber.Ctx, md User) *fiber.Error
	DeleteAccount(c *fiber.Ctx, otp string) *fiber.Error

	// ADMIN ROLE
	RestoreUser(c *fiber.Ctx, email string) *fiber.Error
	DeleteUser(c *fiber.Ctx, id uint) *fiber.Error
}

type UserRepository interface {
	// funtions
	DeleteAuthRedis(givenUuid string) (int64, error)
	GeneratePairToken(userID uint) (Token, error)
	SendVerificationEmail(md User, code string) error

	Register(md User) (User, *fiber.Error)
	Login(md User) (Token, *fiber.Error)
	RefreshToken(payload RefreshTokenInput) (Token, *fiber.Error)
	DeleteToken(authD *AccessDetails) *fiber.Error
	VerificationEmail(code string) *fiber.Error
	ResendVerificationCode(md User) *fiber.Error
	RequestOTPEmail(md User, message string) *fiber.Error
	FindOTPRequest(otp string) (OTPRequest, *fiber.Error)
	FindReferenceOTPRequest(refNo string) (OTPRequest, *fiber.Error)
	VerifyOTP(md OTPRequest) (string, *fiber.Error)
	ResetPassword(md User) *fiber.Error
	ChangePassword(md User) *fiber.Error

	EmailExists(email string) *fiber.Error
	UsernameExists(username string) *fiber.Error
	Create(md User) *fiber.Error
	Update(md User) (User, *fiber.Error)
	Delete(md User) *fiber.Error
	FindUserByIdentity(identity string) (User, *fiber.Error)
	FindUserByEmail(email string) (User, *fiber.Error)
	FindUserById(id uint) (User, *fiber.Error)

	// ADMIN ROLE
	FindDeletedUserByEmail(email string) (User, *fiber.Error)
	RestoreUser(id uint) *fiber.Error
}
