package middleware

import (
	"log"
	"os"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/favicon"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
)

// FiberMiddleware provide Fiber's built-in middlewares.
// See: https://docs.gofiber.io/api/middleware
func FiberMiddleware(a *fiber.App) {
	// LOG FILE WRITER
	file, err := os.OpenFile("logs/logfile.log", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalf("error opening file: %v", err)
	}

	loggerConfig := logger.Config{
		Output: file,
		Format: "${time} | ${status} |${latency} | ${ip} | ${method} | ${path}\n",
		// Format:     "${time} | ${status} |${latency} | ${ip} | ${method} | ${path} | ${resBody}\n",
		TimeFormat: time.RFC1123,
		TimeZone:   "Asia/Jakarta",
		// Done: func(c *fiber.Ctx, logString []byte) {
		// 	if c.Response().StatusCode() != fiber.StatusOK {
		// 		// reporter.SendToSlack(logString)
		// 		log.Println(logString)
		// 	}
		// },
	}

	faviconConfig := favicon.Config{
		File: "./static/favicon.ico",
	}

	// CORS config
	// corsConfig := cors.Config{
	// 	AllowOrigins: "http://localhost:8000",
	// 	// AllowCredentials: true,
	// 	AllowHeaders: "Origin, Content-Type, Accept",
	// }

	a.Use(
		// Add CORS to each route.
		cors.New(),
		// This panic will be caught by the middleware
		recover.New(),
		// Add simple logger.
		logger.New(loggerConfig),
		// Add favicon.
		favicon.New(faviconConfig),
	)

	// Serve static files from the "static" directory
	a.Static("/static", "./static", fiber.Static{
		Compress: true,
	})
	// Serve media files from the "media" directory
	a.Static("/media", "./media", fiber.Static{
		Compress: true,
	})

}
