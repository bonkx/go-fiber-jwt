package helpers

import (
	"encoding/base64"
	"errors"
	"fmt"
	"myapp/pkg/configs"
	"myapp/src/models"
	"strconv"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v4"
	uuid "github.com/satori/go.uuid"
)

func CreateToken(userid uint) (*models.TokenDetails, error) {
	config, _ := configs.LoadConfig(".")

	now := time.Now().UTC()
	td := &models.TokenDetails{
		TokenType: "Bearer",
	}
	td.AtExpires = now.Add(config.AccessTokenExpiresIn).Unix()
	td.AccessUuid = uuid.NewV4().String()
	td.RtExpires = now.Add(config.RefreshTokenExpiresIn).Unix()
	td.RefreshUuid = td.AccessUuid + "++" + strconv.Itoa(int(userid))

	// Creating Access Token
	atDecodedPrivateKey, err := base64.StdEncoding.DecodeString(config.AccessTokenPrivateKey)
	if err != nil {
		return nil, fmt.Errorf("could not decode token private key: %w", err)
	}
	atKey, err := jwt.ParseRSAPrivateKeyFromPEM(atDecodedPrivateKey)

	if err != nil {
		return nil, fmt.Errorf("create: parse token private key: %w", err)
	}

	atClaims := make(jwt.MapClaims)
	atClaims["authorized"] = true
	atClaims["sub"] = userid
	atClaims["token_uuid"] = td.AccessUuid
	atClaims["exp"] = td.AtExpires
	td.AccessToken, err = jwt.NewWithClaims(jwt.SigningMethodRS256, atClaims).SignedString(atKey)
	if err != nil {
		return nil, fmt.Errorf("create: sign access token: %w", err)
	}

	// Creating Refresh Token
	rtDecodedPrivateKey, err := base64.StdEncoding.DecodeString(config.RefreshTokenPrivateKey)
	if err != nil {
		return nil, fmt.Errorf("could not decode token private key: %w", err)
	}
	rtKey, err := jwt.ParseRSAPrivateKeyFromPEM(rtDecodedPrivateKey)

	if err != nil {
		return nil, fmt.Errorf("create: parse token private key: %w", err)
	}

	rtClaims := make(jwt.MapClaims)
	rtClaims["sub"] = userid
	rtClaims["token_uuid"] = td.RefreshUuid
	rtClaims["exp"] = td.RtExpires
	td.RefreshToken, err = jwt.NewWithClaims(jwt.SigningMethodRS256, rtClaims).SignedString(rtKey)
	if err != nil {
		return nil, fmt.Errorf("create: sign token: %w", err)
	}

	return td, nil
}

func ValidateToken(token string, publicKey string) (*models.AccessDetails, error) {
	decodedPublicKey, err := base64.StdEncoding.DecodeString(publicKey)
	if err != nil {
		return nil, fmt.Errorf("could not decode: %w", err)
	}

	key, err := jwt.ParseRSAPublicKeyFromPEM(decodedPublicKey)

	if err != nil {
		return nil, fmt.Errorf("validate: parse key: %w", err)
	}

	parsedToken, err := jwt.Parse(token, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodRSA); !ok {
			return nil, fmt.Errorf("unexpected method: %s", t.Header["alg"])
		}
		return key, nil
	})

	if err != nil {
		return nil, fmt.Errorf("validate: %w", err)
	}

	claims, ok := parsedToken.Claims.(jwt.MapClaims)
	if !ok || !parsedToken.Valid {
		return nil, fmt.Errorf("validate: invalid token")
	}

	tokenUuid, ok := claims["token_uuid"].(string)
	if !ok {
		return nil, err
	}
	userId, err := strconv.ParseUint(fmt.Sprintf("%.f", claims["sub"]), 10, 64)
	if err != nil {
		return nil, err
	}

	return &models.AccessDetails{
		TokenUuid: tokenUuid,
		UserID:    uint(userId),
	}, nil
}

func ExtractToken(c *fiber.Ctx) string {
	bearToken := c.Get("Authorization")

	// Normally Authorization HTTP header.
	onlyToken := strings.Split(bearToken, " ")
	if len(onlyToken) == 2 {
		return onlyToken[1]
	}

	return ""
}

func ExtractTokenMetadata(c *fiber.Ctx) (*models.AccessDetails, error) {
	config, _ := configs.LoadConfig(".")

	tokenString := ExtractToken(c)
	if tokenString == "" {
		return nil, errors.New("Unauthorized! No credentials provided.")
	}
	token, err := ValidateToken(tokenString, config.AccessTokenPublicKey)
	if err != nil {
		return nil, err
	}
	return token, nil
}
