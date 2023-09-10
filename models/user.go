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
	FirstName   string     `json:"first_name" binding:"required"`
	LastName    string     `json:"last_name" binding:"required"`
	Username    string     `json:"username" gorm:"not null;unique"`
	Email       string     `json:"email" binding:"required" gorm:"unique"`
	Password    string     `json:"-" binding:"required"`
	Verified    *bool      `gorm:"not null;default:false"`
	IsSuperuser *bool      `json:"is_superuser" gorm:"default:false"`
	IsStaff     *bool      `json:"is_staff" gorm:"default:false"`
	LastLogin   *time.Time `json:"last_login"`
	VerifiedAt  *time.Time `json:"verified_at"`

	UserProfile UserProfile `gorm:"foreignkey:UserID"`
}

func (user *User) BeforeSave(*gorm.DB) error {
	passwordHash, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	user.Password = string(passwordHash)
	user.Username = html.EscapeString(strings.TrimSpace(user.Username))
	user.Email = strings.ToLower(user.Email)

	// if user.Username == "admin" {
	// 	*user.Verified = true
	// 	*user.IsSuperuser = true
	// 	*user.IsStaff = true
	// 	*user.VerifiedAt = time.Now()
	// 	user.UserProfile.Role = "admin"
	// }
	return nil
}

func (user *User) ValidatePassword(password string) error {
	return bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
}

type UserEntity interface {
	Register(ctx context.Context, payload RegisterInput) (User, error)
	Login(ctx context.Context, payload AuthenticationInput) (Token, error)

	GetMe(ctx context.Context) (User, error)

	Create(ctx context.Context, md User) error
	FindUserByUsername(ctx context.Context, username string) (User, error)
	FindUserById(ctx context.Context, id uint) (User, error)
}

type UserRepository interface {
	Register(ctx context.Context, payload RegisterInput) (User, error)
	Login(ctx context.Context, payload AuthenticationInput) (Token, error)

	GetMe(ctx context.Context) (User, error)

	Create(ctx context.Context, md User) error
	FindUserByUsername(ctx context.Context, username string) (User, error)
	FindUserById(ctx context.Context, id uint) (User, error)
}

// func (user *User) Save() (*User, error) {
// 	err := database.Database.Create(&user).Error
// 	if err != nil {
// 		return &User{}, err
// 	}
// 	return user, nil
// }

// func FindUserByUsername(username string) (User, error) {
// 	var user User
// 	err := database.Database.Where("username=?", username).Find(&user).Error
// 	if err != nil {
// 		return User{}, err
// 	}
// 	return user, nil
// }

// func FindUserById(id uint) (User, error) {
// 	var user User
// 	err := database.Database.Preload("Entries").Where("ID=?", id).Find(&user).Error
// 	if err != nil {
// 		return User{}, err
// 	}
// 	return user, nil
// }
