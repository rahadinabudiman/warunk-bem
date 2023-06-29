package middlewares

import (
	"errors"
	"strings"
	"warunk-bem/author"

	"github.com/dgrijalva/jwt-go"
	"github.com/labstack/echo"
)

func IsUser(c echo.Context) (string, error) {
	token := c.Request().Header.Get("Authorization")
	if token == "" {
		return "", echo.NewHTTPError(400, "Missing Token")
	}

	// Extract the token from the "Bearer <token>" format
	tokenParts := strings.Split(token, " ")
	if len(tokenParts) != 2 || strings.ToLower(tokenParts[0]) != "bearer" {
		return "", echo.NewHTTPError(400, "Invalid Token Format")
	}

	// Parse and validate the JWT token
	jwtToken, err := jwt.Parse(tokenParts[1], func(token *jwt.Token) (interface{}, error) {
		// Replace "your-secret-key" with the actual secret key used to sign the tokens
		// You may need to retrieve the secret key from your configuration or environment variables
		return []byte(author.App.Config.GetString("SECRET_JWT")), nil
	})
	if err != nil {
		return "", echo.NewHTTPError(400, "Invalid Token")
	}

	claims, ok := jwtToken.Claims.(jwt.MapClaims)
	if !ok || !jwtToken.Valid {
		return "", echo.NewHTTPError(401, "Unauthorized")
	}

	if claims["is_admin"] != false {
		return "", echo.NewHTTPError(401, "Unauthorized")
	}

	// Extract the admin ID from the token's payload
	id, ok := claims["name"].(string)
	if !ok {
		return "", errors.New("failed to retrieve name from JWT token")
	}

	return id, nil
}
