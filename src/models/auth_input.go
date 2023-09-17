package models

type RegisterInput struct {
	Username        string `json:"username" validate:"required,gte=4"`
	Password        string `json:"password" validate:"required,gte=4"`
	PasswordConfirm string `json:"password_confirm"`
	FirstName       string `json:"first_name" validate:"required"`
	LastName        string `json:"last_name" validate:"required"`
	Email           string `json:"email" validate:"required,email,gte=4"`
	Phone           string `json:"phone" validate:"required"`
}

type LoginInput struct {
	Email    string `json:"email" validate:"required,gte=4"`
	Password string `json:"password" validate:"required,gte=4"`
}

type RefreshTokenInput struct {
	RefreshToken string `json:"refresh_token" validate:"required"`
}
