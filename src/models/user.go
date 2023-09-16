package models

import (
	"context"
	"html"
	"strings"
	"time"

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
		return err
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
	Register(ctx context.Context, payload RegisterInput) (User, error)
	Login(ctx context.Context, payload LoginInput) (Token, error)
	RefreshToken(ctx context.Context, payload RefreshTokenInput) (Token, error)
	VerificationEmail(ctx context.Context, code string) error
	ResendVerificationCode(ctx context.Context, email string) error

	// Create(ctx context.Context, md User) error
	// Update(ctx context.Context, md User) error
	// Delete(ctx context.Context, md User) error
}

type UserRepository interface {
	Register(ctx context.Context, md User) (User, error)
	Login(ctx context.Context, md User) (Token, error)
	RefreshToken(ctx context.Context, payload RefreshTokenInput) (Token, error)
	VerificationEmail(ctx context.Context, code string) error
	ResendVerificationCode(md User) error

	EmailExists(email string) error
	UsernameExists(username string) error
	Create(ctx context.Context, md User) error
	// Update(ctx context.Context, md User) error
	// Delete(ctx context.Context, md User) error
	FindUserByIdentity(ctx context.Context, identity string) (User, error)
	FindUserByEmail(ctx context.Context, email string) (User, error)
	FindUserById(ctx context.Context, id uint) (User, error)
}
