package models

type RegisterInput struct {
	Username        string `json:"username" validate:"required"`
	Password        string `json:"password" validate:"required"`
	PasswordConfirm string `json:"password_confirm"`
	FirstName       string `json:"first_name" validate:"required"`
	LastName        string `json:"last_name" validate:"required"`
	Email           string `json:"email" validate:"required"`
	Phone           string `json:"phone" validate:"required"`
}

type LoginInput struct {
	Username string `json:"username" validate:"required"`
	Password string `json:"password" validate:"required"`
}
type RefreshTokenInput struct {
	RefreshToken string `json:"refresh_token" validate:"required"`
}
