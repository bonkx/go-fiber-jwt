package models

type EmailInput struct {
	Email string `json:"email" validate:"required,email"`
}

type OTPInput struct {
	Otp string `json:"otp" validate:"required,min=6,max=6"`
}
