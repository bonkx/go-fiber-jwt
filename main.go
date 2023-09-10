package main

import (
	"fmt"
	"go-gin/controllers"
	"go-gin/database"
	"go-gin/models"
	"log"
	"net/http"

	"github.com/gin-gonic/contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func loadEnv() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal("Error loading .env file")
	}
}

func loadDatabase() {
	database.InitDatabase()
	database.Database.AutoMigrate(&models.User{})
}

func main() {
	loadEnv()
	loadDatabase()
	serveApplication()
}

func healthCheckHandler(c *gin.Context) {
	// var superSecretKey string = goDotEnvVariable("SUPER_SECRET_KEY")

	// if superSecretKey == "" {
	// 	c.JSON(http.StatusInternalServerError, gin.H{
	// 		"error": "Missing environment variables",
	// 	})
	// 	return
	// }

	c.JSON(http.StatusOK, gin.H{
		"healthy": true,
		"message": "Golang API with Gin Gonic...",
	})
}

func serveApplication() {
	r := gin.Default()
	r.Use(cors.Default())

	r.GET("/", healthCheckHandler)
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})

	apiGroup := r.Group("/api")
	v1Group := apiGroup.Group("/v1")

	publicRoutes := v1Group.Group("/auth")
	publicRoutes.POST("/register", controllers.Register)
	publicRoutes.POST("/login", controllers.Login)

	r.Run(":8080")
	fmt.Println("Server running on port 8080")
}
