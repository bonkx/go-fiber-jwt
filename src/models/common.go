package models

type EmailInput struct {
	Email string `json:"email" validate:"required,email"`
}

type OTPInput struct {
	Otp string `json:"otp" validate:"required,min=6,max=6"`
}

type ResetPasswordInput struct {
	ReferenceNo     string `json:"reference_no" validate:"required"`
	Password        string `json:"password" validate:"required,gte=4"`
	PasswordConfirm string `json:"password_confirm"`
}

type ChangePasswordInput struct {
	Password        string `json:"password" validate:"required,gte=4"`
	PasswordConfirm string `json:"password_confirm"`
}
