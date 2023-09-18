package routes

import (
	"github.com/gofiber/fiber/v2"
)

// NotFoundRoute func for describe 404 Error route.
func NotFoundRoute(a *fiber.App) {
	// Register new special route.
	a.Use(
		// Anonymous function.
		func(c *fiber.Ctx) error {
			// Return HTTP 404 status and JSON response.
			errD := fiber.NewError(fiber.StatusNotFound, "Sorry, Endpoint is not found")
			return c.Status(errD.Code).JSON(errD)
		},
	)
}
