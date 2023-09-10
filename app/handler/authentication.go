package handler

import (
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
	var payload models.AuthenticationInput

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
