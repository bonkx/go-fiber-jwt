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
	FirstName        string     `json:"first_name" binding:"required"`
	LastName         string     `json:"last_name" binding:"required"`
	Username         string     `json:"username" gorm:"not null;unique"`
	Email            string     `json:"email" binding:"required" gorm:"unique"`
	Password         string     `json:"-" binding:"required"`
	Verified         bool       `json:"verified" gorm:"not null;default:false"`
	IsSuperuser      bool       `json:"is_superuser" gorm:"default:false"`
	IsStaff          bool       `json:"is_staff" gorm:"default:false"`
	LastLogin        *time.Time `json:"last_login"`
	VerificationCode string     `json:"verification_code"`
	VerifiedAt       *time.Time `json:"verified_at"`

	UserProfile UserProfile `gorm:"foreignkey:UserID"`
}

func (user *User) BeforeCreate(*gorm.DB) error {
	passwordHash, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return fmt.Errorf("could not hash password %w", err)
	}
	user.Password = string(passwordHash)
	user.Username = html.EscapeString(strings.TrimSpace(user.Username))
	user.Email = strings.ToLower(user.Email)

	return nil
}

func (user *User) AfterCreate(tx *gorm.DB) (err error) {
	if user.Username == "admin" {
		tx.Model(user).Updates(User{IsSuperuser: true, IsStaff: true})
		tx.Model(user.UserProfile).Update("role", "admin")
	}
	return
}

func (user *User) ValidatePassword(password string) error {
	return bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
}

type UserUsecase interface {
	Register(ctx context.Context, payload RegisterInput) (User, *fiber.Error)
	Login(ctx context.Context, payload LoginInput) (Token, *fiber.Error)
	RefreshToken(ctx context.Context, payload RefreshTokenInput) (Token, *fiber.Error)
	VerificationEmail(ctx context.Context, code string) *fiber.Error
	ResendVerificationCode(ctx context.Context, email string) *fiber.Error
	Logout(authD *AccessDetails) *fiber.Error

	ForgotPassword(ctx context.Context, payload EmailInput) *fiber.Error
	// Update(ctx context.Context, md User) error
	// Delete(ctx context.Context, md User) error
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
	RequestOTPEmail(md User) *fiber.Error

	EmailExists(email string) *fiber.Error
	UsernameExists(username string) *fiber.Error
	Create(md User) *fiber.Error
	// Update( md User) *fiber.Error
	// Delete( md User) *fiber.Error
	FindUserByIdentity(identity string) (User, *fiber.Error)
	FindUserByEmail(email string) (User, *fiber.Error)
	FindUserById(id uint) (User, *fiber.Error)
}
