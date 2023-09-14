package configs

import (
	"encoding/json"
	"errors"
	"os"
	"strconv"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/template/html/v2"
)

// FiberConfig func for configuration Fiber app.
// See: https://docs.gofiber.io/api/fiber#config
func FiberConfig() fiber.Config {
	// Initialize standard Go html template engine
	engine := html.New("templates", ".html")

	// Define server settings.
	readTimeoutSecondsCount, _ := strconv.Atoi(os.Getenv("SERVER_READ_TIMEOUT"))

	// Return Fiber configuration.
	return fiber.Config{
		ReadTimeout: time.Second * time.Duration(readTimeoutSecondsCount),
		JSONEncoder: json.Marshal,
		JSONDecoder: json.Unmarshal,
		Views:       engine,
		BodyLimit:   25 * 1024 * 1024, // the default limit of 4MB
		// Override default error handler
		ErrorHandler: func(ctx *fiber.Ctx, err error) error {
			// Status code defaults to 500
			code := fiber.StatusInternalServerError

			// Retrieve the custom status code if it's a *fiber.Error
			var e *fiber.Error
			if errors.As(err, &e) {
				code = e.Code
			}

			// Send custom error page
			// err = ctx.Status(code).SendFile(fmt.Sprintf("./%d.html", code))
			// if err != nil {
			// 	// In case the SendFile fails
			// 	return ctx.Status(fiber.StatusInternalServerError).SendString("Internal Server Error")
			// }

			// Set Content-Type: text/plain; charset=utf-8
			ctx.Set(fiber.HeaderContentType, fiber.MIMETextPlainCharsetUTF8)

			// Return status code with error message
			// return ctx.Status(code).SendString(err.Error())
			return ctx.Status(code).JSON(fiber.Map{
				"code":    code,
				"message": err.Error(),
			})
		},
	}
}
