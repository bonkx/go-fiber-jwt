package middleware

import (
	"myapp/pkg/utils"
	"strconv"

	"github.com/gofiber/fiber/v2"
)

func PaginationMiddleware() fiber.Handler {
	return func(c *fiber.Ctx) error {
		per_page, err := strconv.ParseInt(c.Query("per_page"), 10, 0)
		if err != nil {
			per_page = 10
		}

		page, err := strconv.ParseInt(c.Query("page"), 10, 0)
		if err != nil {
			page = 1
		}

		c.Set("per_page", utils.IntToString(per_page))
		c.Set("page", utils.IntToString(page))

		return c.Next()
	}
}
