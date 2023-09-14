package main

import (
	"log"
	"myapp/pkg/configs"
	"myapp/pkg/middleware"
	"myapp/src/routes"

	"github.com/gofiber/fiber/v2"
)

func init() {
	envConfig, err := configs.LoadConfig(".")
	if err != nil {
		log.Fatalln("Failed to load environment variables! \n", err.Error())
	}
	configs.ConnectDB(&envConfig)
	// configs.ConnectRedis(&envConfig)
}

func main() {
	// Define Fiber config.
	config := configs.FiberConfig()

	// Define a new Fiber app with config.
	app := fiber.New(config)

	// Middlewares.
	middleware.FiberMiddleware(app) // Register Fiber's middleware for app.

	// get DB
	db := configs.GetDBConnection()

	// Routes.
	routes.SwaggerRoute(app)      // Register a route for API Docs (Swagger).
	routes.PublicRoutes(app, db)  // Register a public routes for app.
	routes.PrivateRoutes(app, db) // Register a private routes for app.
	routes.NotFoundRoute(app)     // Register route for 404 Error.

	// Start server (with graceful shutdown).
	configs.StartServer(app)
}
