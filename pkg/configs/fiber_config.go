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
	engine := html.New("./templates", ".html")

	// Define server settings.
	readTimeoutSecondsCount, _ := strconv.Atoi(os.Getenv("SERVER_READ_TIMEOUT"))

	// Return Fiber configuration.
	return fiber.Config{
		ReadTimeout: time.Second * time.Duration(readTimeoutSecondsCount),
		JSONEncoder: json.Marshal,
		JSONDecoder: json.Unmarshal,
		Views:       engine,
		// ViewsLayout: "layouts/main",
		BodyLimit: 5 * 1024 * 1024, // the default limit: 4MB
		// Override default error handler
		ErrorHandler: func(c *fiber.Ctx, err error) error {
			// Status code defaults to 500
			code := fiber.StatusInternalServerError
			// fmt.Println(err)
			// Retrieve the custom status code if it's a *fiber.Error
			var e *fiber.Error
			if errors.As(err, &e) {
				code = e.Code
			}

			// Send custom error page
			// err = c.Status(code).SendFile(fmt.Sprintf("./%d.html", code))
			// if err != nil {
			// 	// In case the SendFile fails
			// 	return c.Status(fiber.StatusInternalServerError).SendString("Internal Server Error")
			// }

			// Set Content-Type: text/plain; charset=utf-8
			c.Set(fiber.HeaderContentType, fiber.MIMETextPlainCharsetUTF8)

			// Return status code with error message
			// return c.Status(code).SendString(err.Error())
			return c.Status(code).JSON(fiber.Map{
				"code":    code,
				"message": err.Error(),
			})
		},
	}
}
