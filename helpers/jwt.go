package helpers

import (
	"errors"
	"fmt"
	"go-gin/models"
	"os"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
)

var privateKey = []byte(os.Getenv("JWT_SECRET"))

func GenerateJWT(user models.User) (models.Token, error) {
	tokenTTL, _ := time.ParseDuration(os.Getenv("JWT_TOKEN_TTL"))
	// fmt.Println("tokenTTL: ", tokenTTL)
	now := time.Now()

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"id":  user.ID,
		"iat": now.Unix(),
		"eat": now.Add(tokenTTL).Unix(),
	})

	jwt := models.Token{}

	tokenString, err := token.SignedString(privateKey)
	if err != nil {
		return jwt, fmt.Errorf("generating JWT Token failed: %w", err)
	}

	jwt.AccessToken = tokenString
	jwt.TokenType = "Bearer"
	jwt.ExpiresIn = int64(tokenTTL.Seconds())

	return jwt, nil
}

func ValidateJWT(context *gin.Context) error {
	token, err := getToken(context)
	if err != nil {
		return err
	}
	_, ok := token.Claims.(jwt.MapClaims)
	if ok && token.Valid {
		return nil
	}
	return errors.New("invalid token provided")
}

func CurrentUser(context *gin.Context) (models.User, error) {
	err := ValidateJWT(context)
	if err != nil {
		return models.User{}, err
	}
	token, _ := getToken(context)
	claims, _ := token.Claims.(jwt.MapClaims)
	userId := uint(claims["id"].(float64))

	user, err := models.FindUserById(userId)
	if err != nil {
		return models.User{}, err
	}
	return user, nil
}

func getToken(context *gin.Context) (*jwt.Token, error) {
	tokenString := getTokenFromRequest(context)
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		return privateKey, nil
	})
	return token, err
}

func getTokenFromRequest(context *gin.Context) string {
	bearerToken := context.Request.Header.Get("Authorization")
	splitToken := strings.Split(bearerToken, " ")
	if len(splitToken) == 2 {
		return splitToken[1]
	}
	return ""
}
