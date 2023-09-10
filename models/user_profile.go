package models

import (
	"time"

	"gorm.io/gorm"
)

type UserProfile struct {
	gorm.Model
	UserID            uint
	StatusID          uint
	Phone             string     `json:"phone" binding:"required" gorm:"size:20;"`
	Photo             string     `json:"photo"`
	Role              string     `json:"role" gorm:"size:100;"`
	LastLoginIp       string     `json:"last_login_ip"`
	LoginWithSosmed   *bool      `json:"login_with_sosmed" gorm:"default:false"`
	LoginWithSosmedAt *time.Time `json:"login_with_sosmed_at"`
	Birthday          *time.Time `json:"birthday"`
	IsPhoneVerified   bool       `json:"is_phone_verified" gorm:"default:false"`
	PhoneVerifiedAt   *time.Time `json:"phone_verified_at"`
	PhoneVerifiedOtp  string     `json:"phone_verified_otp"`
	Status            Status     `gorm:"foreignkey:StatusID"`
}
