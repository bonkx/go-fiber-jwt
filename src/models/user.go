package models

import (
	"context"
	"encoding/json"
	"fmt"
	"html"
	"myapp/pkg/response"
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
	LastLoginAt      *time.Time `json:"last_login_at"`
	LastLoginIp      string     `json:"last_login_ip"`
	VerificationCode string     `json:"verification_code"`
	VerifiedAt       *time.Time `json:"verified_at"`

	UserProfile UserProfile `gorm:"foreignkey:UserID;constraint:OnDelete:CASCADE;" json:"user_profile,omitempty"`
	Products    []Product   `gorm:"foreignkey:UserID;constraint:OnDelete:CASCADE;" json:"products,omitempty"`
}

func (md User) MarshalJSON() ([]byte, error) {
	type Alias User // NOTE this will not copy methods present in Struct
	aux := struct {
		Name string
		Alias
	}{
		Alias: (Alias)(md),
		Name:  fmt.Sprintf("%s %s", md.FirstName, md.LastName),
	}
	return json.Marshal(aux)
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

func (u *User) BeforeDelete(tx *gorm.DB) (err error) {
	if u.UserProfile.Role == "admin" {
		// return errors.New("admin user not allowed to delete")
		return fmt.Errorf("admin user not allowed to delete")
	}
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
	// USECASE
	Register(ctx context.Context, payload RegisterInput) *fiber.Error
	Login(ctx context.Context, payload LoginInput) (Token, *fiber.Error)
	RefreshToken(ctx context.Context, payload RefreshTokenInput) (Token, *fiber.Error)
	VerificationEmail(ctx context.Context, code string) *fiber.Error
	ResendVerificationCode(ctx context.Context, email string) *fiber.Error
	Logout(authD *AccessDetails) *fiber.Error

	ForgotPassword(ctx context.Context, payload EmailInput) *fiber.Error
	ForgotPasswordOTP(ctx context.Context, payload OTPInput) (string, *fiber.Error)
	ResetPassword(ctx context.Context, payload ResetPasswordInput) *fiber.Error
	ChangePassword(ctx context.Context, md User, payload ChangePasswordInput) *fiber.Error
	UpdateProfile(c *fiber.Ctx, payload UpdateProfileInput) (User, *fiber.Error)
	// Delete(ctx context.Context, md User) *fiber.Error
	UploadPhotoProfile(c *fiber.Ctx, md User) *fiber.Error
	RequestDeleteAccount(c *fiber.Ctx, md User) *fiber.Error
	DeleteAccount(c *fiber.Ctx, otp string) *fiber.Error
	ListUser(c *fiber.Ctx) (*response.Pagination, []*User, *fiber.Error)

	// ADMIN ROLE
	RestoreUser(c *fiber.Ctx, email string) *fiber.Error
	DeleteUser(c *fiber.Ctx, id uint) *fiber.Error
	PermanentDeleteUser(c *fiber.Ctx, id uint) *fiber.Error
}

type UserRepository interface {
	// FUNTIONS
	DeleteAuthRedis(givenUuid string) (int64, error)
	GeneratePairToken(userID uint) (Token, error)
	SendVerificationEmail(obj User, code string) error

	// REPOS
	Register(obj User) *fiber.Error
	Login(obj User) (Token, *fiber.Error)
	RefreshToken(payload RefreshTokenInput) (Token, *fiber.Error)
	DeleteToken(authD *AccessDetails) *fiber.Error
	VerificationEmail(code string) *fiber.Error
	ResendVerificationCode(obj User) *fiber.Error
	RequestOTPEmail(obj User, message string) *fiber.Error
	FindOTPRequest(otp string) (OTPRequest, *fiber.Error)
	FindReferenceOTPRequest(refNo string) (OTPRequest, *fiber.Error)
	VerifyOTP(obj OTPRequest) (string, *fiber.Error)
	ResetPassword(obj User) *fiber.Error
	ChangePassword(obj User) *fiber.Error

	EmailExists(email string) *fiber.Error
	UsernameExists(username string) *fiber.Error
	Create(obj User) *fiber.Error
	Update(obj User) (User, *fiber.Error)
	FindUserByIdentity(identity string) (User, *fiber.Error)
	FindUserByEmail(email string) (User, *fiber.Error)
	FindUserById(id uint) (User, *fiber.Error)
	ListUser(param response.ParamsPagination) (*response.Pagination, []*User, *fiber.Error)

	// ADMIN ROLE
	FindDeletedUserByEmail(email string) (User, *fiber.Error)
	RestoreUser(id uint) *fiber.Error
	Delete(obj User) *fiber.Error
	PermanentDelete(obj User) *fiber.Error
}
