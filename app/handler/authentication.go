package handler

import (
	"fmt"
	middleware "myapp/middlewares"
	"myapp/models"
	"net/http"

	"github.com/gin-gonic/gin"
)

type AuthHandler struct {
	entity models.UserEntity
}

func NewAuthHandler(r *gin.RouterGroup, us models.UserEntity) {
	handler := &AuthHandler{
		entity: us,
	}

	auth := r.Group("/auth")
	auth.POST("/register", handler.Register)
	auth.POST("/login", handler.Login)
	auth.POST("/refresh", handler.RefreshAccessToken)

	r.GET("/me", middleware.JWTAuthMiddleware(), handler.GetMe)
}

func (h *AuthHandler) Register(ctx *gin.Context) {
	var payload models.RegisterInput

	if err := ctx.ShouldBindJSON(&payload); err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	if payload.Password != payload.PasswordConfirm {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "Passwords do not match!"})
		return
	}

	savedUser, err := h.entity.Register(ctx.Request.Context(), payload)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	message := "We sent an email with a verification link to " + savedUser.Email
	ctx.JSON(http.StatusCreated, gin.H{
		"user":    savedUser,
		"message": message,
	})
}

func (h *AuthHandler) Login(ctx *gin.Context) {
	var payload models.LoginInput

	if err := ctx.ShouldBindJSON(&payload); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	// errors := models.ValidateStruct(payload)
	// if errors != nil {
	// 	return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": "fail", "errors": errors})
	// }

	token, err := h.entity.Login(ctx.Request.Context(), payload)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	// config, _ := initializers.LoadConfig(".")
	// // set cookie
	// ctx.SetCookie("access_token", token.AccessToken, config.AccessTokenMaxAge*60, "/", "localhost", false, true)
	// ctx.SetCookie("refresh_token", token.RefreshToken, config.RefreshTokenMaxAge*60, "/", "localhost", false, true)
	// ctx.SetCookie("logged_in", "true", config.AccessTokenMaxAge*60, "/", "localhost", false, false)

	ctx.JSON(http.StatusOK, &token)
}

func (h *AuthHandler) GetMe(ctx *gin.Context) {
	// A *model.User will eventually be added to context in middleware
	user, exists := ctx.Get("user")
	if !exists {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": "Unable to extract user from request context for unknown reason",
		})
		return
	}

	ctx.JSON(http.StatusOK, user)
}
func (h *AuthHandler) RefreshAccessToken(ctx *gin.Context) {
	fmt.Println("do RefreshAccessToken ======================")
	var payload models.RefreshTokenInput

	if err := ctx.ShouldBindJSON(&payload); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}
	// fmt.Printf("RefreshTokenInput : %v\n", payload.RefreshToken)

	token, err := h.entity.RefreshToken(ctx.Request.Context(), payload)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, &token)
}

// func LogoutUser(c *fiber.Ctx) error {
// 	message := "Token is invalid or session has expired"

// 	refresh_token := c.Cookies("refresh_token")

// 	if refresh_token == "" {
// 		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"status": "fail", "message": message})
// 	}

// 	config, _ := initializers.LoadConfig(".")
// 	ctx := context.TODO()

// 	tokenClaims, err := utils.ValidateToken(refresh_token, config.RefreshTokenPublicKey)
// 	if err != nil {
// 		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"status": "fail", "message": err.Error()})
// 	}

// 	access_token_uuid := c.Locals("access_token_uuid").(string)
// 	_, err = initializers.RedisClient.Del(ctx, tokenClaims.TokenUuid, access_token_uuid).Result()
// 	if err != nil {
// 		return c.Status(fiber.StatusBadGateway).JSON(fiber.Map{"status": "fail", "message": err.Error()})
// 	}

// 	expired := time.Now().Add(-time.Hour * 24)
// 	c.Cookie(&fiber.Cookie{
// 		Name:    "access_token",
// 		Value:   "",
// 		Expires: expired,
// 	})
// 	c.Cookie(&fiber.Cookie{
// 		Name:    "refresh_token",
// 		Value:   "",
// 		Expires: expired,
// 	})
// 	c.Cookie(&fiber.Cookie{
// 		Name:    "logged_in",
// 		Value:   "",
// 		Expires: expired,
// 	})
// 	return c.Status(fiber.StatusOK).JSON(fiber.Map{"status": "success"})
// }
