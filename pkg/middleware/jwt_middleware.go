package middleware

import (
	"errors"
	"myapp/pkg/configs"
	"myapp/pkg/helpers"
	"myapp/src/models"
	"strings"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func JWTAuthMiddleware() fiber.Handler {
	return func(c *fiber.Ctx) error {
		var access_token string
		authorization := c.Get("Authorization")

		if strings.HasPrefix(authorization, "Bearer ") {
			access_token = strings.TrimPrefix(authorization, "Bearer ")
		}

		if access_token == "" {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"message": errors.New(fiber.ErrUnauthorized.Message),
				"error":   errors.New("Unauthorized! No credentials provided."),
			})
		}

		config, _ := configs.LoadConfig(".")

		tokenClaims, err := helpers.ValidateToken(access_token, config.AccessTokenPublicKey)
		if err != nil {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"message": errors.New(fiber.ErrForbidden.Message),
				"error":   err.Error(),
			})
		}

		var user models.User
		err = configs.DB.Preload("UserProfile.Status").First(&user, "id = ?", tokenClaims.UserID).Error

		if err == gorm.ErrRecordNotFound {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"message": errors.New(fiber.ErrForbidden.Message),
				"error":   "the user belonging to this token no logger exists",
			})
		}

		c.Locals("user", user)
		c.Locals("access_token_uuid", tokenClaims.TokenUuid)

		return c.Next()
	}
}
