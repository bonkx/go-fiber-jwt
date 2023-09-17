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

func NewAuthHandler(r fiber.Router, userUsecase models.UserUsecase) {
	handler := &AuthHandler{
		userUsecase: userUsecase,
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

func (h *AuthHandler) Register(c *fiber.Ctx) error {
	var payload models.RegisterInput

	if err := c.BodyParser(&payload); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(
			fiber.NewError(fiber.StatusBadRequest, err.Error()),
		)
	}

	if payload.Password != payload.PasswordConfirm {
		return c.Status(fiber.StatusBadRequest).JSON(
			fiber.NewError(400, "Passwords do not match!"),
		)
	}

	// form POST validation
	errors := models.ValidateStruct(payload)
	if errors != nil {
		errD := models.ErrorDetailsResponse{
			Code:    fiber.ErrUnprocessableEntity.Code,
			Message: fiber.ErrUnprocessableEntity.Message,
			Errors:  errors,
		}
		return c.Status(errD.Code).JSON(errD)
	}

	if payload.Phone != "" {
		phone_number_validated := utils.FormatPhoneNumber(payload.Phone)
		errors := models.ValidatePhoneNumber(phone_number_validated)
		if errors != nil {
			errD := models.ErrorDetailsResponse{
				Code:    fiber.StatusUnprocessableEntity,
				Message: fiber.ErrUnprocessableEntity.Message,
				Errors:  errors,
			}
			return c.Status(errD.Code).JSON(errD)
		}

		payload.Phone = phone_number_validated
	}

	savedUser, err := h.userUsecase.Register(c.Context(), payload)
	if err != nil {
		return c.Status(err.Code).JSON(err)
	}

	message := "We sent an email with a verification link to " + savedUser.Email
	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"code":    fiber.StatusCreated,
		"user":    savedUser,
		"message": message,
	})
}

func (h *AuthHandler) RequestVerifyCode(c *fiber.Ctx) error {
	var payload models.EmailInput

	if err := c.BodyParser(&payload); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(
			fiber.NewError(fiber.StatusBadRequest, err.Error()),
		)
	}

	// form POST validations
	errors := models.ValidateStruct(payload)
	if errors != nil {
		errD := models.ErrorDetailsResponse{
			Code:    fiber.ErrUnprocessableEntity.Code,
			Message: fiber.ErrUnprocessableEntity.Message,
			Errors:  errors,
		}
		return c.Status(errD.Code).JSON(errD)
	}

	err := h.userUsecase.ResendVerificationCode(c.Context(), payload.Email)
	if err != nil {
		return c.Status(err.Code).JSON(err)
	}

	message := "We sent an email with a verification link to " + payload.Email
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"code":    fiber.StatusOK,
		"message": message,
	})
}

func (h *AuthHandler) Login(c *fiber.Ctx) error {
	var payload models.LoginInput

	if err := c.BodyParser(&payload); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(
			fiber.NewError(fiber.StatusBadRequest, err.Error()),
		)
	}

	// form POST validation
	errors := models.ValidateStruct(payload)
	if errors != nil {
		errD := models.ErrorDetailsResponse{
			Code:    fiber.ErrUnprocessableEntity.Code,
			Message: fiber.ErrUnprocessableEntity.Message,
			Errors:  errors,
		}
		return c.Status(errD.Code).JSON(errD)
	}

	token, err := h.userUsecase.Login(c.Context(), payload)
	if err != nil {
		return c.Status(err.Code).JSON(err)
	}

	// config, _ := configs.LoadConfig(".")
	// // set cookie
	// ctx.SetCookie("access_token", token.AccessToken, config.AccessTokenMaxAge*60, "/", "localhost", false, true)
	// ctx.SetCookie("refresh_token", token.RefreshToken, config.RefreshTokenMaxAge*60, "/", "localhost", false, true)
	// ctx.SetCookie("logged_in", "true", config.AccessTokenMaxAge*60, "/", "localhost", false, false)

	return c.Status(fiber.StatusOK).JSON(&token)
}

func (h *AuthHandler) RefreshAccessToken(c *fiber.Ctx) error {
	fmt.Println("do RefreshAccessToken ======================")
	var payload models.RefreshTokenInput

	if err := c.BodyParser(&payload); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(
			fiber.NewError(fiber.StatusBadRequest, err.Error()),
		)
	}

	// form POST validations
	errors := models.ValidateStruct(payload)
	if errors != nil {
		errD := models.ErrorDetailsResponse{
			Code:    fiber.ErrUnprocessableEntity.Code,
			Message: fiber.ErrUnprocessableEntity.Message,
			Errors:  errors,
		}
		return c.Status(errD.Code).JSON(errD)
	}

	token, err := h.userUsecase.RefreshToken(c.Context(), payload)
	if err != nil {
		return c.Status(err.Code).JSON(err)
	}

	return c.Status(fiber.StatusOK).JSON(token)
}

func (h *AuthHandler) Logout(c *fiber.Ctx) error {
	token, err := helpers.ExtractTokenMetadata(c)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(
			fiber.NewError(fiber.StatusUnauthorized, err.Error()),
		)
	}

	if err := h.userUsecase.Logout(token); err != nil {
		return c.Status(err.Code).JSON(err)
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"code":    fiber.StatusOK,
		"message": "Logout successful",
	})
}

func (h *AuthHandler) ForgotPassword(c *fiber.Ctx) error {
	var payload models.EmailInput

	if err := c.BodyParser(&payload); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(
			fiber.NewError(fiber.StatusBadRequest, err.Error()),
		)
	}

	// form POST validations
	errors := models.ValidateStruct(payload)
	if errors != nil {
		errD := models.ErrorDetailsResponse{
			Code:    fiber.ErrUnprocessableEntity.Code,
			Message: fiber.ErrUnprocessableEntity.Message,
			Errors:  errors,
		}
		return c.Status(errD.Code).JSON(errD)
	}

	err := h.userUsecase.ForgotPassword(c.Context(), payload)
	if err != nil {
		return c.Status(err.Code).JSON(err)
	}

	message := "We sent an email with a OTP code to " + payload.Email
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"code":    fiber.StatusOK,
		"message": message,
	})
}

func (h *AuthHandler) ForgotPasswordOTP(c *fiber.Ctx) error {
	var payload models.OTPInput

	if err := c.BodyParser(&payload); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(
			fiber.NewError(fiber.StatusBadRequest, err.Error()),
		)
	}

	// form POST validations
	errors := models.ValidateStruct(payload)
	if errors != nil {
		errD := models.ErrorDetailsResponse{
			Code:    fiber.ErrUnprocessableEntity.Code,
			Message: fiber.ErrUnprocessableEntity.Message,
			Errors:  errors,
		}
		return c.Status(errD.Code).JSON(errD)
	}

	refNo, err := h.userUsecase.ForgotPasswordOTP(c.Context(), payload)
	if err != nil {
		return c.Status(err.Code).JSON(err)
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"code":         fiber.StatusOK,
		"message":      "OTP verification successful",
		"reference_no": refNo,
	})
}

func (h *AuthHandler) ResetPassword(c *fiber.Ctx) error {
	var payload models.ResetPasswordInput

	if err := c.BodyParser(&payload); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(
			fiber.NewError(fiber.StatusBadRequest, err.Error()),
		)
	}

	if payload.Password != payload.PasswordConfirm {
		return c.Status(fiber.StatusBadRequest).JSON(
			fiber.NewError(400, "Passwords do not match!"),
		)
	}

	// form POST validations
	errors := models.ValidateStruct(payload)
	if errors != nil {
		errD := models.ErrorDetailsResponse{
			Code:    fiber.ErrUnprocessableEntity.Code,
			Message: fiber.ErrUnprocessableEntity.Message,
			Errors:  errors,
		}
		return c.Status(errD.Code).JSON(errD)
	}

	err := h.userUsecase.ResetPassword(c.Context(), payload)
	if err != nil {
		return c.Status(err.Code).JSON(err)
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"code":    fiber.StatusOK,
		"message": "Your password has been successfully reset",
	})
}
