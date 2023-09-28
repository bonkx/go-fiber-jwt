package usecase

import (
	"context"
	"myapp/pkg/response"
	"myapp/pkg/utils"
	"myapp/src/models"
	"strconv"
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

// ListUser implements models.UserUsecase.
func (uc *UserUsecase) ListUser(c *fiber.Ctx) (*response.Pagination, []*models.User, *fiber.Error) {
	// 	Parse the query parameters
	search := c.Query("search")
	sortBy := c.Query("sort", "id|desc")
	page := c.Query("page", "1")
	limit := c.Query("per_page", "10")

	// Convert the page and limit to integers
	pageInt, _ := strconv.Atoi(page)
	limitInt, _ := strconv.Atoi(limit)

	sortQuery, errSort := utils.ValidateAndReturnSortQuery(sortBy)
	// log.Print(sortQuery)
	if errSort != nil {
		errD := fiber.NewError(fiber.StatusInternalServerError, errSort.Error())
		return nil, nil, errD
	}

	// make param pagination struct
	pagParam := response.ParamsPagination{
		Page:      pageInt,
		Limit:     limitInt,
		SortQuery: sortQuery,
		Search:    search,
		NoPage:    c.Query("no_page"),
	}

	pagination, data, err := uc.userRepo.ListUser(pagParam)
	if err != nil {
		return nil, nil, err
	}
	return pagination, data, nil
}

// PermanentDeleteUser implements models.UserUsecase.
func (uc *UserUsecase) PermanentDeleteUser(c *fiber.Ctx, id uint) *fiber.Error {
	// get user data
	user, err := uc.userRepo.FindUserById(id)
	if err != nil {
		return err
	}

	// permanent deleted user
	err = uc.userRepo.PermanentDelete(user)
	if err != nil {
		return err
	}

	return nil
}

// DeleteUser implements models.UserUsecase.
func (uc *UserUsecase) DeleteUser(c *fiber.Ctx, id uint) *fiber.Error {
	// get user data
	user, err := uc.userRepo.FindUserById(id)
	if err != nil {
		return err
	}

	// deleted user
	err = uc.userRepo.Delete(user)
	if err != nil {
		return err
	}

	return nil
}

// RestoreUser implements models.UserUsecase.
func (uc *UserUsecase) RestoreUser(c *fiber.Ctx, email string) *fiber.Error {
	// find user from email
	user, err := uc.userRepo.FindDeletedUserByEmail(email)
	if err != nil {
		return err
	}

	// check user is deleted
	if !user.DeletedAt.Valid {
		return fiber.NewError(422, "Unable to process, this account exists in the database")
	}

	// restore deleted user by ID
	err = uc.userRepo.RestoreUser(user.ID)
	if err != nil {
		return err
	}

	return nil
}

// DeleteAccount implements models.UserUsecase.
func (uc *UserUsecase) DeleteAccount(c *fiber.Ctx, otp string) *fiber.Error {
	user, errLocal := c.Locals("user").(models.User)
	if !errLocal {
		return fiber.NewError(500, utils.ERR_CURRENT_USER_NOT_FOUND)
	}

	// find OTP Request
	_, err := uc.userRepo.FindOTPRequest(otp)
	if err != nil {
		return err
	}

	err = uc.userRepo.Delete(user)
	if err != nil {
		return err
	}

	return nil
}

// RequestDeleteAccount implements models.UserUsecase.
func (uc *UserUsecase) RequestDeleteAccount(c *fiber.Ctx, user models.User) *fiber.Error {
	message := "Here, your OTP code for delete the account:"

	// do request OTP
	if err := uc.userRepo.RequestOTPEmail(user, message); err != nil {
		return err
	}

	return nil
}

// UploadPhotoProfile implements models.UserUsecase.
func (uc *UserUsecase) UploadPhotoProfile(c *fiber.Ctx, user models.User) *fiber.Error {
	// MultipartForm POST
	if form, err := c.MultipartForm(); err == nil {
		files := form.File["file"]

		// Loop through files:
		for _, file := range files {
			// fmt.Println(file.Filename, file.Size, file.Header["Content-Type"][0])
			// => "avatar1.jpeg" 6472 "image/jpeg"

			// // Save the files to disk:
			imageUrl, errFile := utils.ImageUpload(c, file, "users")
			if errFile != nil {
				return fiber.NewError(500, errFile.Error())
			}

			// update user photo path
			user.UserProfile.Photo = &imageUrl
		}
	}

	// do update user
	_, err := uc.userRepo.Update(user)
	if err != nil {
		return err
	}

	return nil
}

// UpdateProfile implements models.UserUsecase.
func (uc *UserUsecase) UpdateProfile(c *fiber.Ctx, payload models.UpdateProfileInput) (models.User, *fiber.Error) {
	// user := models.User{}
	user, errLocal := c.Locals("user").(models.User)
	if !errLocal {
		return user, fiber.NewError(500, utils.ERR_CURRENT_USER_NOT_FOUND)
	}

	// MultipartForm POST
	if form, err := c.MultipartForm(); err == nil {
		files := form.File["file"]

		// Loop through files:
		for _, file := range files {
			// fmt.Println(file.Filename, file.Size, file.Header["Content-Type"][0])
			// => "avatar1.jpeg" 6472 "image/jpeg"

			// // Save the files to disk:
			imageUrl, errFile := utils.ImageUpload(c, file, "users")
			if errFile != nil {
				return user, fiber.NewError(500, errFile.Error())
			}

			// update user photo path
			user.UserProfile.Photo = &imageUrl
		}
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
	user, err := uc.userRepo.Update(user)
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

	message := "Here, your OTP code for reset your password:"
	// do request OTP
	if err := uc.userRepo.RequestOTPEmail(user, message); err != nil {
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
func (uc *UserUsecase) Register(ctx context.Context, payload models.RegisterInput) *fiber.Error {

	// cek email of user
	if err := uc.userRepo.EmailExists(payload.Email); err != nil {
		return err
	}

	// cek username of user
	if err := uc.userRepo.UsernameExists(payload.Username); err != nil {
		return err
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

	err := uc.userRepo.Register(user)
	if err != nil {
		return err
	}
	return nil
}
