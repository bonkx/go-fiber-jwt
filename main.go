package main

import (
	"log"
	"myapp/pkg/configs"
	"myapp/pkg/middleware"
	"myapp/routes"

	"github.com/gofiber/fiber/v2"
)

func init() {
	envConfig, err := configs.LoadConfig(".")
	if err != nil {
		log.Fatalln("Failed to load environment variables! \n", err.Error())
	}
	configs.ConnectDB(&envConfig)
	configs.ConnectRedis(&envConfig)
	configs.MigrateDB()
}

// @title Fiber Example API
// @version 1.0
// @description This is a sample swagger for Fiber
// @termsOfService http://swagger.io/terms/
// @contact.name API Support
// @contact.email fiber@swagger.io
// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html
// @host localhost:8000
// @BasePath /api/
// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @description Type "Bearer" followed by a space and JWT token.
func main() {
	// Define Fiber config.
	config := configs.FiberConfig()

	// Define a new Fiber app with config.
	app := fiber.New(config)

	// Register Fiber's middleware for app.
	middleware.FiberMiddleware(app)

	// get DB
	db := configs.GetDBConnection()

	// Routes.
	routes.PublicRoutes(app, db) // Register a public routes for app.
	routes.APIRoutes(app, db)    // Register a API routes for app.
	routes.SwaggerRoute(app)     // Register a route for API Docs (Swagger).
	// place at end of routes
	routes.NotFoundRoute(app) // Register route for 404 Error.

	// Start server (with graceful shutdown).
	configs.StartServer(app)
}
