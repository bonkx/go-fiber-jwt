package handler

import (
	middleware "myapp/pkg/middleware"
	"myapp/src/models"

	"github.com/gofiber/fiber/v2"
)

type AccountHandler struct {
	userUsecase models.UserUsecase
}

func NewAccountHandler(r fiber.Router, uc models.UserUsecase) {
	handler := &AccountHandler{
		userUsecase: uc,
	}

	// ROUTES
	acc := r.Group("/accounts")

	// private API
	acc.Get("/me", middleware.JWTAuthMiddleware(), handler.GetMe)
	acc.Post("/change-password", middleware.JWTAuthMiddleware(), handler.ChangePassword)
	acc.Put("/update", middleware.JWTAuthMiddleware(), handler.UpdateProfile)
	acc.Post("/photo", middleware.JWTAuthMiddleware(), handler.UploadPhotoProfile)

	acc.Post("/delete", middleware.JWTAuthMiddleware(), handler.RequestDeleteAccount)
	acc.Delete("/delete", middleware.JWTAuthMiddleware(), handler.DeleteAccount)
}

// GetMe
// @Summary      GetMe
// @Description  Get User Data
// @Tags         Accounts
// @Accept       json
// @Produce      json
// @Success      200  {object}  models.User
// @Failure      500  {object}  models.ResponseError
// @Security 	 BearerAuth
// @Router		 /v1/accounts/me [get]
func (h *AccountHandler) GetMe(c *fiber.Ctx) error {
	// A *model.User will eventually be added to context in middleware
	user, errLocal := c.Locals("user").(models.User)
	if !errLocal {
		res := models.ResponseHTTP{
			Code:    fiber.ErrInternalServerError.Code,
			Message: "Unable to extract user from request context for unknown reason",
		}
		return c.Status(res.Code).JSON(res)
	}

	return c.Status(fiber.StatusOK).JSON(user)
}

// ChangePassword
// @Summary      Change Password
// @Description  Change your old password to new password
// @Tags         Accounts
// @Accept       json
// @Produce      json
// @Param 		 body body models.ChangePasswordInput true "Body"
// @Success      200  {object}  models.ResponseSuccess
// @Failure      400  {object}  models.ResponseError
// @Failure      422  {object}  models.ResponseHTTP
// @Failure      500  {object}  models.ResponseError
// @Security 	 BearerAuth
// @Router       /v1/accounts/change-password [post]
func (h *AccountHandler) ChangePassword(c *fiber.Ctx) error {
	var payload models.ChangePasswordInput
	res := models.ResponseHTTP{
		Code:    fiber.StatusOK,
		Message: "Your password has been changed successfully",
	}

	if err := c.BodyParser(&payload); err != nil {
		res.Code = fiber.StatusBadRequest
		res.Message = err.Error()
		return c.Status(res.Code).JSON(res)
	}

	if payload.Password != payload.PasswordConfirm {
		res.Code = fiber.StatusBadRequest
		res.Message = "Passwords do not match!"
		return c.Status(res.Code).JSON(res)
	}

	// form POST validations
	errD := models.ValidateStruct(payload)
	if errD.Errors != nil {
		return c.Status(errD.Code).JSON(errD)
	}

	user, errLocal := c.Locals("user").(models.User)
	if !errLocal {
		res.Code = fiber.StatusInternalServerError
		res.Message = "Unable to extract user from request context for unknown reason"
		return c.Status(res.Code).JSON(res)
	}

	err := h.userUsecase.ChangePassword(c.Context(), user, payload)
	if err != nil {
		res.Code = err.Code
		res.Message = err.Message
		return c.Status(res.Code).JSON(res)
	}

	return c.Status(res.Code).JSON(res)
}

// UpdateProfile
// @Summary      Update Profile
// @Description  Update your profile
// @Tags         Accounts
// @Accept       json
// @Accept       multipart/form-data
// @Produce      json
// @Param 		 body formData models.UpdateProfileInput true "Body"
// @Param 		 file formData file false "File to upload" format(multipart/form-data)
// @Success      200  {object}  models.User
// @Failure      400  {object}  models.ResponseError
// @Failure      422  {object}  models.ResponseHTTP
// @Failure      500  {object}  models.ResponseError
// @Security 	 BearerAuth
// @Router       /v1/accounts/update [put]
func (h *AccountHandler) UpdateProfile(c *fiber.Ctx) error {
	var payload models.UpdateProfileInput
	res := models.ResponseHTTP{
		Code:    fiber.StatusOK,
		Message: "Request has been processed successfully",
	}

	if err := c.BodyParser(&payload); err != nil {
		res.Code = fiber.StatusBadRequest
		res.Message = err.Error()
		return c.Status(res.Code).JSON(res)
	}

	// form POST validations
	errD := models.ValidateStruct(payload)
	if errD.Errors != nil {
		return c.Status(errD.Code).JSON(errD)
	}

	user, err := h.userUsecase.UpdateProfile(c, payload)
	if err != nil {
		res.Code = err.Code
		res.Message = err.Message
		return c.Status(res.Code).JSON(res)
	}

	return c.Status(res.Code).JSON(user)
}

// UploadPhotoProfile
// @Summary      Upload Photo Profile
// @Description  Change your photo profile
// @Tags         Accounts
// @Accept       json
// @Accept       multipart/form-data
// @Produce      json
// @Param 		 file formData file true "File to upload" format(multipart/form-data)
// @Success      200  {object}  models.ResponseSuccess
// @Failure      500  {object}  models.ResponseError
// @Security 	 BearerAuth
// @Router       /v1/accounts/photo [post]
func (h *AccountHandler) UploadPhotoProfile(c *fiber.Ctx) error {
	res := models.ResponseHTTP{
		Code:    fiber.StatusOK,
		Message: "Request has been processed successfully",
	}

	user, errLocal := c.Locals("user").(models.User)
	if !errLocal {
		res.Code = fiber.StatusInternalServerError
		res.Message = "Unable to extract user from request context for unknown reason"
		return c.Status(res.Code).JSON(res)
	}

	err := h.userUsecase.UploadPhotoProfile(c, user)
	if err != nil {
		res.Code = err.Code
		res.Message = err.Message
		return c.Status(res.Code).JSON(res)
	}

	return c.Status(res.Code).JSON(res)
}

// RequestDeleteAccount
// @Summary      Request Delete Account
// @Description  Request OTP for account deletion
// @Tags         Accounts
// @Accept       json
// @Produce      json
// @Success      200  {object}  models.ResponseSuccess
// @Failure      500  {object}  models.ResponseError
// @Security 	 BearerAuth
// @Router       /v1/accounts/delete [post]
func (h *AccountHandler) RequestDeleteAccount(c *fiber.Ctx) error {
	res := models.ResponseHTTP{
		Code:    fiber.StatusOK,
		Message: "Request has been processed successfully",
	}

	user, errLocal := c.Locals("user").(models.User)
	if !errLocal {
		res.Code = fiber.StatusInternalServerError
		res.Message = "Unable to extract user from request context for unknown reason"
		return c.Status(res.Code).JSON(res)
	}

	err := h.userUsecase.RequestDeleteAccount(c, user)
	if err != nil {
		res.Code = err.Code
		res.Message = err.Message
		return c.Status(res.Code).JSON(res)
	}

	res.Message = "We sent an email with a OTP code to your email. Check your inbox."
	return c.Status(res.Code).JSON(res)
}

// DeleteAccount
// @Summary      Delete Account
// @Description  Delete your account
// @Tags         Accounts
// @Accept       json
// @Produce      json
// @Param 		 body body models.OTPInput true "Body"
// @Success      200  {object}  models.ResponseSuccess
// @Failure      400  {object}  models.ResponseError
// @Failure      422  {object}  models.ResponseHTTP
// @Failure      500  {object}  models.ResponseError
// @Security 	 BearerAuth
// @Router       /v1/accounts/delete [delete]
func (h *AccountHandler) DeleteAccount(c *fiber.Ctx) error {
	var payload models.OTPInput
	res := models.ResponseHTTP{
		Code:    fiber.StatusOK,
		Message: "Request has been processed successfully",
	}

	if err := c.BodyParser(&payload); err != nil {
		res.Code = fiber.StatusBadRequest
		res.Message = err.Error()
		return c.Status(res.Code).JSON(res)
	}

	// form POST validations
	errD := models.ValidateStruct(payload)
	if errD.Errors != nil {
		return c.Status(errD.Code).JSON(errD)
	}

	err := h.userUsecase.DeleteAccount(c, payload.Otp)
	if err != nil {
		res.Code = err.Code
		res.Message = err.Message
		return c.Status(res.Code).JSON(res)
	}

	return c.Status(res.Code).JSON(res)
}
