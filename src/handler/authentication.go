package handler

import (
	"fmt"
	"myapp/pkg/helpers"
	middleware "myapp/pkg/middleware"
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
	auth.Post("/login", handler.Login)
	auth.Post("/refresh", handler.RefreshAccessToken)

	r.Get("/me", middleware.JWTAuthMiddleware(), handler.GetMe)
}

func (h *AuthHandler) Register(c *fiber.Ctx) error {
	var payload models.RegisterInput

	if err := c.BodyParser(&payload); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": err.Error(),
		})
	}

	if payload.Password != payload.PasswordConfirm {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Passwords do not match!",
		})
	}

	// form POST validation
	errors := models.ValidateStruct(&payload)
	if errors != nil {
		errD := models.ErrorDetailsResponse{
			Code:    fiber.ErrUnprocessableEntity.Code,
			Message: fiber.ErrUnprocessableEntity.Message,
			Errors:  errors,
		}
		return c.Status(fiber.StatusUnprocessableEntity).JSON(errD)
	}

	if payload.Phone != "" {
		phone_number_validated := helpers.FormatPhoneNumber(payload.Phone)
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
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": err.Error(),
		})
	}

	message := "We sent an email with a verification link to " + savedUser.Email

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"user":    savedUser,
		"message": message,
	})
}

func (h *AuthHandler) Login(c *fiber.Ctx) error {
	var payload models.LoginInput

	if err := c.BodyParser(&payload); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": err.Error(),
		})
	}

	// form POST validation
	errors := models.ValidateStruct(&payload)
	if errors != nil {
		errD := models.ErrorDetailsResponse{
			Code:    fiber.ErrUnprocessableEntity.Code,
			Message: fiber.ErrUnprocessableEntity.Message,
			Errors:  errors,
		}
		return c.Status(fiber.StatusUnprocessableEntity).JSON(errD)
	}

	token, err := h.userUsecase.Login(c.Context(), payload)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": err.Error(),
		})
	}

	// config, _ := initializers.LoadConfig(".")
	// // set cookie
	// ctx.SetCookie("access_token", token.AccessToken, config.AccessTokenMaxAge*60, "/", "localhost", false, true)
	// ctx.SetCookie("refresh_token", token.RefreshToken, config.RefreshTokenMaxAge*60, "/", "localhost", false, true)
	// ctx.SetCookie("logged_in", "true", config.AccessTokenMaxAge*60, "/", "localhost", false, false)

	return c.Status(fiber.StatusOK).JSON(&token)
}

func (h *AuthHandler) GetMe(c *fiber.Ctx) error {
	// A *model.User will eventually be added to context in middleware
	user, err := c.Locals("user").(models.User)
	if !err {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Unable to extract user from request context for unknown reason",
		})
	}

	return c.Status(fiber.StatusOK).JSON(user)
}

func (h *AuthHandler) RefreshAccessToken(c *fiber.Ctx) error {
	fmt.Println("do RefreshAccessToken ======================")
	var payload models.RefreshTokenInput

	if err := c.BodyParser(&payload); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": err.Error(),
		})
	}

	token, err := h.userUsecase.RefreshToken(c.Context(), payload)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(token)
}

// func LogoutUser(c *fiber.Ctx) error {
// 	message := "Token is invalid or session has expired"

// 	refresh_token := c.Cookies("refresh_token")

// 	if refresh_token == "" {
// 		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"status": "fail", "message": message})
// 	}

// 	config, _ := initializers.LoadConfig(".")
// 	ctx := context.TODO()

// 	tokenClaims, err := utils.ValidateToken(refresh_token, config.RefreshTokenPublicKey)
// 	if err != nil {
// 		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"status": "fail", "message": err.Error()})
// 	}

// 	access_token_uuid := c.Locals("access_token_uuid").(string)
// 	_, err = initializers.RedisClient.Del(ctx, tokenClaims.TokenUuid, access_token_uuid).Result()
// 	if err != nil {
// 		return c.Status(fiber.StatusBadGateway).JSON(fiber.Map{"status": "fail", "message": err.Error()})
// 	}

// 	expired := time.Now().Add(-time.Hour * 24)
// 	c.Cookie(&fiber.Cookie{
// 		Name:    "access_token",
// 		Value:   "",
// 		Expires: expired,
// 	})
// 	c.Cookie(&fiber.Cookie{
// 		Name:    "refresh_token",
// 		Value:   "",
// 		Expires: expired,
// 	})
// 	c.Cookie(&fiber.Cookie{
// 		Name:    "logged_in",
// 		Value:   "",
// 		Expires: expired,
// 	})
// 	return c.Status(fiber.StatusOK).JSON(fiber.Map{"status": "success"})
// }
