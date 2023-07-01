package usecase

import (
	"net/http"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
)

// SetJwtAdmin sets JWT middleware for admin routes
func (h *JwtUsecase) SetJwtAdmin(g *gin.RouterGroup) {
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

		// Call the validateJwtAdmin function
		validateJwtAdmin(c)
	})
}

// validateJwtAdmin is a middleware for validating access to admin-only resources
func validateJwtAdmin(c *gin.Context) {
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

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "Forbidden"})
		return
	}

	isAdmin, ok := claims["is_admin"].(bool)
	if !ok || !isAdmin {
		c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "Forbidden"})
		return
	}

	c.Next()
}
