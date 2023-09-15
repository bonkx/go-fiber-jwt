package configs

import (
	"fmt"
	"log"
	"os"

	"github.com/gofiber/fiber/v2"
)

// StartServer func for starting a simple server.
func StartServer(a *fiber.App) {
	// Run server.
	port := os.Getenv("PORT")
	serverPort := fmt.Sprintf("0.0.0.0:%s", port)

	if err := a.Listen(serverPort); err != nil {
		log.Printf("Oops... Server is not running! Reason: %v", err)
	}
}
