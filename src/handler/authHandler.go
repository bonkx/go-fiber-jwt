package handler

import (
	"fmt"
	"myapp/pkg/helpers"
	middleware "myapp/pkg/middleware"
	"myapp/pkg/utils"
	"myapp/src/models"

	"github.com/gofiber/fiber/v2"
)

type AuthHandler struct {
	userUsecase models.UserUsecase
}

func NewAuthHandler(r fiber.Router, uc models.UserUsecase) {
	handler := &AuthHandler{
		userUsecase: uc,
	}

	// ROUTES
	auth := r.Group("/auth")
	auth.Post("/register", handler.Register)
	auth.Post("/request-verify-code", handler.RequestVerifyCode)
	auth.Post("/login", handler.Login)
	auth.Post("/refresh", handler.RefreshAccessToken)
	auth.Post("/forgot-password", handler.ForgotPassword)
	auth.Post("/forgot-password-otp", handler.ForgotPasswordOTP)
	auth.Post("/reset-password", handler.ResetPassword)

	auth.Post("/logout", middleware.JWTAuthMiddleware(), handler.Logout)

}

// Register
// @Summary      Register new Account
// @Description  Register new Account
// @Tags         Auth
// @Accept       json
// @Produce      json
// @Param 		 body body models.RegisterInput true "Body"
// @Success      201  {object}  models.ResponseSuccess
// @Failure      400  {object}  models.ResponseError
// @Failure      422  {object}  models.ResponseHTTP
// @Failure      500  {object}  models.ResponseError
// @Router       /v1/auth/register [post]
func (h *AuthHandler) Register(c *fiber.Ctx) error {
	var payload models.RegisterInput
	res := models.ResponseHTTP{
		Code:    fiber.StatusOK,
		Message: "Request has been processed successfully",
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

	// form POST validation
	errors := models.ValidateStruct(payload)
	if errors != nil {
		res.Code = fiber.StatusUnprocessableEntity
		res.Message = fiber.ErrUnprocessableEntity.Message
		res.Errors = errors
		return c.Status(res.Code).JSON(res)
	}

	if payload.Phone != "" {
		phone_number_validated := utils.FormatPhoneNumber(payload.Phone)
		errors := models.ValidatePhoneNumber(phone_number_validated)
		if errors != nil {
			res.Code = fiber.StatusUnprocessableEntity
			res.Message = fiber.ErrUnprocessableEntity.Message
			res.Errors = errors
			return c.Status(res.Code).JSON(res)
		}

		payload.Phone = phone_number_validated
	}

	err := h.userUsecase.Register(c.Context(), payload)
	if err != nil {
		res.Code = err.Code
		res.Message = err.Message
		return c.Status(res.Code).JSON(res)
	}

	res.Code = fiber.StatusCreated
	res.Message = "We sent an email with a verification link to " + payload.Email + ". Check your inbox."
	return c.Status(res.Code).JSON(res)
}

// RequestVerifyCode
// @Summary      Request Verification Code
// @Description  Request new Verification Code
// @Tags         Auth
// @Accept       json
// @Produce      json
// @Param 		 body body models.EmailInput true "Body"
// @Success      201  {object}  models.ResponseSuccess
// @Failure      400  {object}  models.ResponseError
// @Failure      422  {object}  models.ResponseHTTP
// @Failure      500  {object}  models.ResponseError
// @Router       /v1/auth/request-verify-code [post]
func (h *AuthHandler) RequestVerifyCode(c *fiber.Ctx) error {
	var payload models.EmailInput
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
	errors := models.ValidateStruct(payload)
	if errors != nil {
		res.Code = fiber.StatusUnprocessableEntity
		res.Message = fiber.ErrUnprocessableEntity.Message
		res.Errors = errors
		return c.Status(res.Code).JSON(res)
	}

	err := h.userUsecase.ResendVerificationCode(c.Context(), payload.Email)
	if err != nil {
		res.Code = err.Code
		res.Message = err.Message
		return c.Status(res.Code).JSON(res)
	}

	res.Message = "We sent an email with a verification link to " + payload.Email + ". Check your inbox."
	return c.Status(res.Code).JSON(res)
}

// Login
// @Summary      Login
// @Description  Get your token
// @Tags         Auth
// @Accept       json
// @Produce      json
// @Param 		 body body models.LoginInput true "Body"
// @Success      200  {object}  models.Token
// @Failure      400  {object}  models.ResponseError
// @Failure      422  {object}  models.ResponseHTTP
// @Failure      500  {object}  models.ResponseError
// @Router       /v1/auth/login [post]
func (h *AuthHandler) Login(c *fiber.Ctx) error {
	var payload models.LoginInput
	res := models.ResponseHTTP{
		Code:    fiber.StatusOK,
		Message: "Request has been processed successfully",
	}

	if err := c.BodyParser(&payload); err != nil {
		res.Code = fiber.StatusBadRequest
		res.Message = err.Error()
		return c.Status(res.Code).JSON(res)
	}

	// form POST validation
	errors := models.ValidateStruct(payload)
	if errors != nil {
		res.Code = fiber.ErrUnprocessableEntity.Code
		res.Message = fiber.ErrUnprocessableEntity.Message
		res.Errors = errors
		return c.Status(res.Code).JSON(res)
	}

	token, err := h.userUsecase.Login(c.Context(), payload)
	if err != nil {
		res.Code = err.Code
		res.Message = err.Message
		return c.Status(res.Code).JSON(res)
	}

	// config, _ := configs.LoadConfig(".")
	// // set cookie
	// ctx.SetCookie("access_token", token.AccessToken, config.AccessTokenMaxAge*60, "/", "localhost", false, true)
	// ctx.SetCookie("refresh_token", token.RefreshToken, config.RefreshTokenMaxAge*60, "/", "localhost", false, true)
	// ctx.SetCookie("logged_in", "true", config.AccessTokenMaxAge*60, "/", "localhost", false, false)

	return c.Status(res.Code).JSON(&token)
}

// RefreshAccessToken
// @Summary      Refresh Access Token
// @Description  Refresh your access token
// @Tags         Auth
// @Accept       json
// @Produce      json
// @Param 		 body body models.RefreshTokenInput true "Body"
// @Success      200  {object}  models.Token
// @Failure      400  {object}  models.ResponseError
// @Failure      422  {object}  models.ResponseHTTP
// @Failure      500  {object}  models.ResponseError
// @Router       /v1/auth/refresh [post]
func (h *AuthHandler) RefreshAccessToken(c *fiber.Ctx) error {
	fmt.Println("do RefreshAccessToken ======================")
	var payload models.RefreshTokenInput
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
	errors := models.ValidateStruct(payload)
	if errors != nil {
		res.Code = fiber.StatusUnprocessableEntity
		res.Message = fiber.ErrUnprocessableEntity.Message
		res.Errors = errors
		return c.Status(res.Code).JSON(res)
	}

	token, err := h.userUsecase.RefreshToken(c.Context(), payload)
	if err != nil {
		res.Code = err.Code
		res.Message = err.Message
		return c.Status(res.Code).JSON(res)
	}

	return c.Status(res.Code).JSON(token)
}

// Logout
// @Summary      Logout
// @Description  Revoke your token
// @Tags         Auth
// @Accept       json
// @Produce      json
// @Success      200  {object}  models.ResponseSuccess
// @Failure      400  {object}  models.ResponseError
// @Failure      401  {object}  models.ResponseError
// @Failure      422  {object}  models.ResponseHTTP
// @Failure      500  {object}  models.ResponseError
// @Security 	 BearerAuth
// @Router       /v1/auth/logout [post]
func (h *AuthHandler) Logout(c *fiber.Ctx) error {
	res := models.ResponseHTTP{
		Code:    fiber.StatusOK,
		Message: "Request has been processed successfully",
	}

	token, err := helpers.ExtractTokenMetadata(c)
	if err != nil {
		res.Code = fiber.StatusUnauthorized
		res.Message = err.Error()
		return c.Status(res.Code).JSON(res)
	}

	if err := h.userUsecase.Logout(token); err != nil {
		res.Code = err.Code
		res.Message = err.Message
		return c.Status(res.Code).JSON(res)
	}

	return c.Status(res.Code).JSON(res)
}

// ForgotPassword
// @Summary      Forgot Password
// @Description  Request email with OTP for reset your password
// @Tags         Auth
// @Accept       json
// @Produce      json
// @Param 		 body body models.EmailInput true "Body"
// @Success      200  {object}  models.ResponseSuccess
// @Failure      400  {object}  models.ResponseError
// @Failure      422  {object}  models.ResponseHTTP
// @Failure      500  {object}  models.ResponseError
// @Router       /v1/auth/forgot-password [post]
func (h *AuthHandler) ForgotPassword(c *fiber.Ctx) error {
	var payload models.EmailInput
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
	errors := models.ValidateStruct(payload)
	if errors != nil {
		res.Code = fiber.StatusUnprocessableEntity
		res.Message = fiber.ErrUnprocessableEntity.Message
		res.Errors = errors
		return c.Status(res.Code).JSON(res)
	}

	err := h.userUsecase.ForgotPassword(c.Context(), payload)
	if err != nil {
		res.Code = err.Code
		res.Message = err.Message
		return c.Status(res.Code).JSON(res)
	}

	res.Message = "We sent an email with a OTP code to " + payload.Email + ". Check your inbox."
	return c.Status(res.Code).JSON(res)
}

// ForgotPasswordOTP
// @Summary      Forgot Password OTP
// @Description  Verify OTP your reset password
// @Tags         Auth
// @Accept       json
// @Produce      json
// @Param 		 body body models.OTPInput true "Body"
// @Success      200  {object}  models.ResponseSuccess
// @Failure      400  {object}  models.ResponseError
// @Failure      422  {object}  models.ResponseHTTP
// @Failure      500  {object}  models.ResponseError
// @Router       /v1/auth/forgot-password-otp [post]
func (h *AuthHandler) ForgotPasswordOTP(c *fiber.Ctx) error {
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
	errors := models.ValidateStruct(payload)
	if errors != nil {
		res.Code = fiber.StatusUnprocessableEntity
		res.Message = fiber.ErrUnprocessableEntity.Message
		res.Errors = errors
		return c.Status(res.Code).JSON(res)
	}

	refNo, err := h.userUsecase.ForgotPasswordOTP(c.Context(), payload)
	if err != nil {
		res.Code = err.Code
		res.Message = err.Message
		return c.Status(res.Code).JSON(res)
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"code":         fiber.StatusOK,
		"message":      "OTP verification successful",
		"reference_no": refNo,
	})
}

// ResetPassword
// @Summary      Reset Password
// @Description  Save your new password
// @Tags         Auth
// @Accept       json
// @Produce      json
// @Param 		 body body models.OTPInput true "Body"
// @Success      200  {object}  models.ResponseSuccess
// @Failure      400  {object}  models.ResponseError
// @Failure      422  {object}  models.ResponseHTTP
// @Failure      500  {object}  models.ResponseError
// @Router       /v1/auth/reset-password [post]
func (h *AuthHandler) ResetPassword(c *fiber.Ctx) error {
	var payload models.ResetPasswordInput
	res := models.ResponseHTTP{
		Code:    fiber.StatusOK,
		Message: "Your password has been successfully reset",
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
	errors := models.ValidateStruct(payload)
	if errors != nil {
		res.Code = fiber.ErrUnprocessableEntity.Code
		res.Message = fiber.ErrUnprocessableEntity.Message
		res.Errors = errors
		return c.Status(res.Code).JSON(res)
	}

	err := h.userUsecase.ResetPassword(c.Context(), payload)
	if err != nil {
		res.Code = err.Code
		res.Message = err.Message
		return c.Status(res.Code).JSON(res)
	}

	return c.Status(res.Code).JSON(res)
}
