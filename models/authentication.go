package models

type RegisterInput struct {
	Username        string `json:"username" binding:"required"`
	Password        string `json:"password" binding:"required"`
	PasswordConfirm string `json:"password_confirm"`
	FirstName       string `json:"first_name" binding:"required"`
	LastName        string `json:"last_name" binding:"required"`
	Email           string `json:"email" binding:"required"`
	Phone           string `json:"phone" binding:"required"`
}

type AuthenticationInput struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

// type AuthEntity interface {
// 	Register(ctx context.Context, payload RegisterInput) (User, error)
// 	Login(ctx context.Context, payload AuthenticationInput) (Token, error)
// }

// type AuthRepository interface {
// 	Register(ctx context.Context, payload RegisterInput) (User, error)
// 	Login(ctx context.Context, payload AuthenticationInput) (Token, error)
// }
