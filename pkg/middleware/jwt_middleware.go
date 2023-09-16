package middleware

import (
	"context"
	"myapp/pkg/configs"
	"myapp/pkg/helpers"
	"myapp/src/models"

	"github.com/gofiber/fiber/v2"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

func JWTAuthMiddleware() fiber.Handler {
	return func(c *fiber.Ctx) error {
		access_token := helpers.ExtractToken(c)

		if access_token == "" {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"code":    fiber.ErrUnauthorized.Code,
				"error":   fiber.ErrUnauthorized.Message,
				"message": "Unauthorized! No credentials provided.",
			})
		}

		config, _ := configs.LoadConfig(".")

		tokenClaims, err := helpers.ValidateToken(access_token, config.AccessTokenPublicKey)
		if err != nil {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"code":    fiber.ErrUnauthorized.Code,
				"error":   fiber.ErrUnauthorized.Message,
				"message": err.Error(),
			})
		}

		ctxTodo := context.TODO()
		userid, err := configs.RedisClient.Get(ctxTodo, tokenClaims.TokenUuid).Result()
		if err == redis.Nil {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"code":    fiber.ErrUnauthorized.Code,
				"error":   fiber.ErrUnauthorized.Message,
				"message": "Token is invalid or session has expired",
			})
		}

		var user models.User
		err = configs.DB.Preload("UserProfile.Status").First(&user, "id = ?", userid).Error

		if err == gorm.ErrRecordNotFound {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"code":    fiber.ErrUnauthorized.Code,
				"error":   fiber.ErrUnauthorized.Message,
				"message": "the user belonging to this token no logger exists",
			})
		}

		c.Locals("user", user)
		c.Locals("access_token_uuid", tokenClaims.TokenUuid)

		return c.Next()
	}
}
