package middleware

import (
	"myapp/helpers"
	"myapp/initializers"
	"myapp/models"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func JWTAuthMiddleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var access_token string
		authorization := ctx.Request.Header.Get("Authorization")

		if strings.HasPrefix(authorization, "Bearer ") {
			access_token = strings.TrimPrefix(authorization, "Bearer ")
		}

		if access_token == "" {
			ctx.JSON(http.StatusUnauthorized, gin.H{"message": "Authentication required"})
			ctx.Abort()
			return
		}

		config, _ := initializers.LoadConfig(".")

		tokenClaims, err := helpers.ValidateToken(access_token, config.AccessTokenPublicKey)
		if err != nil {
			ctx.JSON(http.StatusForbidden, gin.H{"message": err.Error()})
			ctx.Abort()
			return
		}

		var user models.User
		err = initializers.DB.Preload("UserProfile.Status").First(&user, "id = ?", tokenClaims.UserID).Error

		if err == gorm.ErrRecordNotFound {
			ctx.JSON(http.StatusForbidden, gin.H{"message": "the user belonging to this token no logger exists"})
			return
		}

		ctx.Set("user", user)
		ctx.Set("access_token_uuid", tokenClaims.TokenUuid)

		ctx.Next()
	}
}

// func JWTAuthMiddleware() gin.HandlerFunc {
// 	return func(context *gin.Context) {
// 		err := helpers.ValidateJWT(context)
// 		if err != nil {
// 			context.JSON(http.StatusUnauthorized, gin.H{"message": "Authentication required"})
// 			context.Abort()
// 			return
// 		}
// 		context.Next()
// 	}
// }
