package usecase

import (
	"context"
	"myapp/pkg/utils"
	"myapp/src/models"
	"time"

	"github.com/gofiber/fiber/v2"
)

type UserUsecase struct {
	userRepo models.UserRepository
}

// NewUserUsecase will create an object that represent the models.UserUsecase interface
func NewUserUsecase(userRepo models.UserRepository) models.UserUsecase {
	return &UserUsecase{
		userRepo: userRepo,
	}
}

// Update implements models.UserUsecase.
func (uc *UserUsecase) Update(c *fiber.Ctx, payload models.UpdateProfileInput) (models.User, *fiber.Error) {
	user := models.User{}
	userLocal, errLocal := c.Locals("user").(models.User)
	if !errLocal {
		return user, fiber.NewError(500, "Unable to extract user from request context for unknown reason")
	}

	// get latest user data
	user, err := uc.userRepo.FindUserById(userLocal.ID)
	if err != nil {
		return user, err
	}

	dateBirthday, errFormat := time.Parse(time.DateOnly, payload.Birthday)
	if errFormat != nil {
		return user, fiber.NewError(500, errFormat.Error())
	}

	// fill user updates
	user.FirstName = payload.FirstName
	user.LastName = payload.LastName
	user.UserProfile.Phone = payload.Phone
	user.UserProfile.Birthday = &dateBirthday

	// do update user
	user, err = uc.userRepo.Update(user)
	if err != nil {
		return user, err
	}

	return user, nil
}

// ChangePassword implements models.UserUsecase.
func (uc *UserUsecase) ChangePassword(ctx context.Context, user models.User, payload models.ChangePasswordInput) *fiber.Error {
	passwordHash, errHash := user.HashPassword(payload.Password)
	if errHash != nil {
		return fiber.NewError(500, errHash.Error())
	}
	// update user password with passwordHash
	user.Password = passwordHash

	// do change password
	err := uc.userRepo.ChangePassword(user)
	if err != nil {
		return err
	}

	return nil
}

// ResetPassword implements models.UserUsecase.
func (uc *UserUsecase) ResetPassword(ctx context.Context, payload models.ResetPasswordInput) *fiber.Error {
	// find OTP Request
	otpR, err := uc.userRepo.FindReferenceOTPRequest(payload.ReferenceNo)
	if err != nil {
		return err
	}

	// find user from ototpR email
	user := models.User{}
	user, err = uc.userRepo.FindUserByEmail(otpR.Email)
	if err != nil {
		return err
	}

	passwordHash, errHash := user.HashPassword(payload.Password)
	if errHash != nil {
		return fiber.NewError(500, errHash.Error())
	}
	// update user password
	user.Password = passwordHash

	// do reset password
	err = uc.userRepo.ResetPassword(user)
	if err != nil {
		return err
	}

	return nil
}

// ForgotPasswordOTP implements models.UserUsecase.
func (uc *UserUsecase) ForgotPasswordOTP(ctx context.Context, payload models.OTPInput) (string, *fiber.Error) {
	// find OTP Request
	otpR, err := uc.userRepo.FindOTPRequest(payload.Otp)
	if err != nil {
		return "", err
	}

	// check otpotpR expired?
	expired, errEx := utils.Expired(time.Now().UTC(), otpR.ExpiredAt.UTC().Format(time.RFC3339))
	if errEx != nil {
		return "", fiber.ErrInternalServerError
	}

	if expired {
		return "", fiber.NewError(422, "OTP code has expired")
	}

	// verify OTP and generate refNo
	refNo, err := uc.userRepo.VerifyOTP(otpR)
	if err != nil {
		return "", err
	}
	return refNo, nil
}

// ForgotPassword implements models.UserUsecase.
func (uc *UserUsecase) ForgotPassword(ctx context.Context, payload models.EmailInput) *fiber.Error {
	// find user based on email
	user, err := uc.userRepo.FindUserByEmail(payload.Email)
	if err != nil {
		return err
	}

	// do request OTP
	if err := uc.userRepo.RequestOTPEmail(user); err != nil {
		return err
	}
	return nil
}

// Logout implements models.UserUsecase.
func (uc *UserUsecase) Logout(authD *models.AccessDetails) *fiber.Error {
	if err := uc.userRepo.DeleteToken(authD); err != nil {
		return err
	}
	return nil
}

// ResendVerificationCode implements models.UserUsecase.
func (uc *UserUsecase) ResendVerificationCode(ctx context.Context, email string) *fiber.Error {
	// get user based on email param
	user, err := uc.userRepo.FindUserByEmail(email)
	if err != nil {
		return err
	}

	err = uc.userRepo.ResendVerificationCode(user)
	if err != nil {
		return err
	}

	return nil
}

// VerificationEmail implements models.UserUsecase.
func (uc *UserUsecase) VerificationEmail(ctx context.Context, code string) *fiber.Error {
	err := uc.userRepo.VerificationEmail(code)
	if err != nil {
		return err
	}
	return nil
}

// RefreshToken implements models.UserUsecase.
func (uc *UserUsecase) RefreshToken(ctx context.Context, payload models.RefreshTokenInput) (models.Token, *fiber.Error) {
	data, err := uc.userRepo.RefreshToken(payload)
	if err != nil {
		return data, err
	}
	return data, nil
}

// Login implements models.UserUsecase.
func (uc *UserUsecase) Login(ctx context.Context, payload models.LoginInput) (models.Token, *fiber.Error) {
	// check email or username exists
	user, err := uc.userRepo.FindUserByIdentity(payload.Email)
	if err != nil {
		return models.Token{}, err
	}

	if !user.Verified {
		return models.Token{}, fiber.NewError(400, "Your account is not active yet, please verify your email.")
	}

	if err := user.ValidatePassword(payload.Password); err != nil {
		return models.Token{}, fiber.NewError(400, "Invalid Email or Password.")
	}

	data, err := uc.userRepo.Login(user)
	if err != nil {
		return models.Token{}, err
	}
	return data, nil
}

// Register implements models.UserUsecase.
func (uc *UserUsecase) Register(ctx context.Context, payload models.RegisterInput) (models.User, *fiber.Error) {

	// cek email of user
	if err := uc.userRepo.EmailExists(payload.Email); err != nil {
		return models.User{}, err
	}

	// cek username of user
	if err := uc.userRepo.UsernameExists(payload.Username); err != nil {
		return models.User{}, err
	}

	user := models.User{
		Username:  payload.Username,
		Password:  payload.Password,
		Email:     payload.Email,
		FirstName: payload.FirstName,
		LastName:  payload.LastName,
		UserProfile: models.UserProfile{
			Phone:    payload.Phone,
			StatusID: 3, // pending
		},
	}

	user, err := uc.userRepo.Register(user)
	if err != nil {
		return user, err
	}
	return user, nil
}
