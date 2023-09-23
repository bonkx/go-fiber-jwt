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

		tokenClaims, err := helpers.ExtractTokenMetadata(c)
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

		SaveUserLogs(c, user)

		c.Locals("user", user)
		c.Locals("token_uuid", tokenClaims.TokenUuid)

		return c.Next()
	}
}

func AdminAuthMiddleware() fiber.Handler {
	return func(c *fiber.Ctx) error {

		tokenClaims, err := helpers.ExtractTokenMetadata(c)
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

		/// ==================== is staff protected ======================
		if !user.IsStaff {
			return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
				"code":    fiber.ErrForbidden.Code,
				"error":   fiber.ErrForbidden.Message,
				"message": "Oops, You Are Not Allowed to Access it!",
			})
		}

		SaveUserLogs(c, user)

		c.Locals("user", user)
		c.Locals("token_uuid", tokenClaims.TokenUuid)

		return c.Next()
	}
}
