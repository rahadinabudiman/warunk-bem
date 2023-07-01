package middlewares

import (
	"errors"
	"strings"
	"warunk-bem/author"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
)

func IsUser(c *gin.Context) (string, error) {
	token := c.GetHeader("Authorization")
	if token == "" {
		return "", errors.New("missing Token")
	}

	// Extract the token from the "Bearer <token>" format
	tokenParts := strings.Split(token, " ")
	if len(tokenParts) != 2 || strings.ToLower(tokenParts[0]) != "bearer" {
		return "", errors.New("invalid Token Format")
	}

	// Parse and validate the JWT token
	jwtToken, err := jwt.Parse(tokenParts[1], func(token *jwt.Token) (interface{}, error) {
		// Replace "your-secret-key" with the actual secret key used to sign the tokens
		// You may need to retrieve the secret key from your configuration or environment variables
		return []byte(author.App.Config.GetString("SECRET_JWT")), nil
	})
	if err != nil {
		return "", errors.New("invalid Token")
	}

	claims, ok := jwtToken.Claims.(jwt.MapClaims)
	if !ok || !jwtToken.Valid {
		return "", errors.New("Unauthorized")
	}

	if claims["is_admin"] != false {
		return "", errors.New("Unauthorized")
	}

	// Extract the admin ID from the token's payload
	id, ok := claims["name"].(string)
	if !ok {
		return "", errors.New("failed to retrieve name from JWT token")
	}

	return id, nil
}
