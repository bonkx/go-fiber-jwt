package models

import (
	"encoding/json"
	"fmt"
	"os"
	"time"

	"gorm.io/gorm"
)

type UserProfile struct {
	gorm.Model
	UserID            uint
	StatusID          uint
	Phone             string     `json:"phone" gorm:"size:20;"`
	Photo             *string    `json:"photo"`
	Role              string     `json:"role" gorm:"size:100;"`
	LoginWithSosmed   bool       `json:"login_with_sosmed" gorm:"default:false"`
	LoginWithSosmedAt *time.Time `json:"login_with_sosmed_at"`
	Birthday          *time.Time `json:"birthday"`
	IsPhoneVerified   bool       `json:"is_phone_verified" gorm:"default:false"`
	PhoneVerifiedAt   *time.Time `json:"phone_verified_at"`
	PhoneVerifiedOtp  string     `json:"phone_verified_otp"`
	Status            Status     `gorm:"foreignkey:StatusID"`
}

func (md UserProfile) MarshalJSON() ([]byte, error) {
	type Alias UserProfile
	var photo *string = md.Photo
	if photo != nil {
		*photo = fmt.Sprintf("%s/%s", os.Getenv("CLIENT_ORIGIN"), *md.Photo)
	}

	aux := struct {
		Alias
		Photo *string `json:"photo"`
	}{
		Alias: (Alias)(md),
		Photo: photo,
	}
	return json.Marshal(aux)
}
