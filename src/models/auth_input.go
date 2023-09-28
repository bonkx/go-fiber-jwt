package models

import "strings"

type RegisterInput struct {
	Username        string `json:"username" validate:"required,gte=4"`
	Password        string `json:"password" validate:"required,gte=4"`
	PasswordConfirm string `json:"password_confirm"`
	FirstName       string `json:"first_name" validate:"required"`
	LastName        string `json:"last_name" validate:"required"`
	Email           string `json:"email" validate:"required,email,gte=4"`
	Phone           string `json:"phone" validate:"required"`
}

func (f *RegisterInput) Sanitize() {
	// https://github.com/leebenson/conform
	f.Username = strings.TrimSpace(f.FirstName)
	f.FirstName = strings.TrimSpace(f.FirstName)
	f.LastName = strings.TrimSpace(f.LastName)
	f.Email = strings.TrimSpace(f.Email)
	f.Password = strings.TrimSpace(f.Password)
	f.PasswordConfirm = strings.TrimSpace(f.PasswordConfirm)
	f.Phone = strings.TrimSpace(f.Phone)
}

type LoginInput struct {
	Email    string `json:"email" validate:"required,gte=4"`
	Password string `json:"password" validate:"required,gte=4"`
}

type RefreshTokenInput struct {
	RefreshToken string `json:"refresh_token" validate:"required"`
}
