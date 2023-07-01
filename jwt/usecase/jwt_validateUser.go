package usecase

import (
	"net/http"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
)

func (h *JwtUsecase) SetJwtUser(g *gin.RouterGroup) {
	secret := h.Config.GetString("SECRET_JWT")

	// Validate JWT token
	g.Use(func(c *gin.Context) {
		tokenString := c.GetHeader("Authorization")

		// Verify token
		_, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
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

		// Call the validateJwtUser function
		validateJwtUser(c)
	})
}

func validateJwtUser(c *gin.Context) {
	_, exists := c.Get("user")
	if !exists {
		c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "Forbidden"})
		return
	}

	// Custom logic for user JWT validation
	// You can add your own checks here

	c.Next()
}
