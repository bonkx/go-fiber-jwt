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
		// siteData, _ := configs.GetSiteData(".")
		// var context = fiber.Map{
		// 	"SiteData": siteData,
		// 	"Subject":  "Golang Fiber JWT",
		// }
		// return c.Render("index", context)
	})

	// metrics route
	a.Get("/metrics", monitor.New(monitor.Config{Title: "App Metrics Page"}))
	// app.Get("/metrics", middleware.AdminProtectedAuth(), monitor.New())

	// // Routes for GET method:
	// route.Get("/books", controllers.GetBooks)              // get list of all books
	// route.Get("/book/:id", controllers.GetBook)            // get one book by ID
	// route.Get("/token/new", controllers.GetNewAccessToken) // create a new access tokens
}
