package routes

import (
	"myapp/pkg/configs"
	"myapp/pkg/helpers"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/monitor"
	"gorm.io/gorm"
)

// PublicRoutes func for describe group of public routes.
func PublicRoutes(a *fiber.App, db *gorm.DB) {
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

	a.Get("/view/register-email", func(c *fiber.Ctx) error {
		siteData, _ := configs.GetSiteData(".")
		// Send Email if register successfully
		// var user models.User
		// err := db.Find(&user).Error
		// if err != nil {
		// 	return c.SendString(err.Error())
		// }
		emailData := helpers.EmailData{
			URL:          siteData.ClientOrigin + "/api/v1/verify-email/" + "code",
			FirstName:    "Admin",
			Subject:      "Your account verification",
			TypeOfAction: "Register",
			SiteData:     siteData,
		}

		return c.Render("emails/verificationCode", emailData)
	})

	// // Routes for GET method:
	// route.Get("/books", controllers.GetBooks)              // get list of all books
	// route.Get("/book/:id", controllers.GetBook)            // get one book by ID
	// route.Get("/token/new", controllers.GetNewAccessToken) // create a new access tokens
}
