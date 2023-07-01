package usecase

import (
	"net/http"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
)

func (h *JwtUsecase) SetJwtGeneral(g *gin.RouterGroup) {
	secret := h.Config.GetString("SECRET_JWT")

	// Validate JWT token
	g.Use(func(c *gin.Context) {
		tokenString := c.GetHeader("Authorization")

		// Verify token
		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			// Validate signing method and secret
			if token.Method != jwt.SigningMethodHS512 {
				c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid signing method"})
				return nil, nil
			}
			return []byte(secret), nil
		})
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
			return
		}

		// Custom logic for general JWT validation
		// You can add your own checks here

		c.Set("user", token)
	})
}

// ValidateGeneralJwt is a middleware for validating general JWT
func (h *JwtUsecase) ValidateGeneralJwt(c *gin.Context) {
	tokenInterface, exists := c.Get("user")
	if !exists {
		c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "Forbidden"})
		return
	}

	token, ok := tokenInterface.(*jwt.Token)
	if !ok || !token.Valid {
		c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "Forbidden"})
		return
	}

	// Custom logic for general JWT validation
	// You can add your own checks here

	c.Next()
}
