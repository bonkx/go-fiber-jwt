package routes

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/monitor"
	"gorm.io/gorm"
)

// PublicRoutes func for describe group of public routes.
func PublicRoutes(a *fiber.App, _ *gorm.DB) {
	// home route
	a.Get("/", func(c *fiber.Ctx) error {
		// return c.SendString("Golang Fiber JWT 👋!")
		return c.Status(200).JSON(fiber.Map{
			"message": "Golang Fiber JWT 👋!",
		})
	})

	// metrics route
	a.Get("/metrics", monitor.New(monitor.Config{Title: "App Metrics Page"}))
	// app.Get("/metrics", middleware.AdminProtectedAuth(), monitor.New())
}
