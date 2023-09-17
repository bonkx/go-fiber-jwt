package models

import (
	"myapp/pkg/utils"
	"strconv"
	"time"

	"gorm.io/gorm"
)

type OTPRequest struct {
	gorm.Model
	Email       string    `json:"email"`
	Otp         string    `json:"otp"`
	ReferenceNo string    `json:"reference_no"`
	ExpiredAt   time.Time `json:"expired_at"`
}

func (m *OTPRequest) BeforeCreate(*gorm.DB) error {
	// generate random OTP code
	randomCode, err := utils.GenerateRandomNumber(6)
	if err != nil {
		return err
	}

	m.Otp = strconv.Itoa(randomCode)
	m.ExpiredAt = time.Now().Local().Add(time.Hour * 24)
	return nil
}
