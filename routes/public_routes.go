package routes

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/monitor"
	"gorm.io/gorm"
)

// PublicRoutes func for describe group of public routes.
func PublicRoutes(a *fiber.App, _ *gorm.DB) {
	// home route
	a.Get("/", HealthCheck)

	// metrics route
	a.Get("/metrics", monitor.New(monitor.Config{Title: "App Metrics Page"}))
	// app.Get("/metrics", middleware.AdminProtectedAuth(), monitor.New())
}

func HealthCheck(c *fiber.Ctx) error {
	res := map[string]interface{}{
		"message": "Golang Fiber JWT 👋!",
		"status":  "Server is up and running",
	}

	if err := c.JSON(res); err != nil {
		return err
	}

	return nil
}
