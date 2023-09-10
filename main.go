package main

import (
	"fmt"
	"io"
	"log"
	"myapp/initializers"
	"net/http"
	"os"

	"github.com/gin-gonic/contrib/cors"
	"github.com/gin-gonic/gin"

	_entity "myapp/app/entity"
	_handler "myapp/app/handler"
	_repo "myapp/app/repository"
)

func init() {
	config, err := initializers.LoadConfig(".")
	if err != nil {
		log.Fatalln("Failed to load environment variables! \n", err.Error())
	}
	initializers.ConnectDB(&config)
	// initializers.ConnectRedis(&config)
}

func main() {
	// Disable Console Color, you don't need console color when writing the logs to file.
	gin.DisableConsoleColor()

	// Logging to a file.
	f, _ := os.Create("logs/gin.log")
	gin.DefaultWriter = io.MultiWriter(f)

	// Use the following code if you need to write the logs to file and console at the same time.
	// gin.DefaultWriter = io.MultiWriter(f, os.Stdout)

	r := gin.Default()

	r.Use(cors.Default())
	r.Use(gin.Logger())

	// Recovery middleware recovers from any panics and writes a 500 if there was one.
	r.Use(gin.Recovery())

	r.GET("/", healthCheckHandler)
	r.GET("/ping", func(ctx *gin.Context) {
		ctx.JSON(200, gin.H{
			"message": "pong",
		})
	})
	r.NoRoute(func(ctx *gin.Context) {
		ctx.JSON(http.StatusNotFound, gin.H{"code": http.StatusNotFound, "message": "404 page not found"})
	})

	r.NoMethod(func(ctx *gin.Context) {
		ctx.JSON(http.StatusMethodNotAllowed, gin.H{"code": http.StatusMethodNotAllowed, "message": "405 method not allowed"})
	})

	apiGroup := r.Group("/api")
	v1Group := apiGroup.Group("/v1")

	db := initializers.GetDBConnection()

	repoUser := _repo.NewUserRepository(db)
	entityUser := _entity.NewUserEntity(repoUser)

	// ROUTES
	_handler.NewAuthHandler(v1Group, entityUser)

	r.Run(":8080")
	fmt.Println("Server running on port 8080")
}

func healthCheckHandler(ctx *gin.Context) {
	// var superSecretKey string = goDotEnvVariable("SUPER_SECRET_KEY")

	// if superSecretKey == "" {
	// 	c.JSON(http.StatusInternalServerError, gin.H{
	// 		"error": "Missing environment variables",
	// 	})
	// 	return
	// }

	ctx.JSON(http.StatusOK, gin.H{
		"healthy": true,
		"message": "Golang API with Gin Gonic...",
	})
}
