package middlewares

import (
	"fmt"
	"net/http"
	"time"
	"warunk-bem/utils"

	"github.com/gin-gonic/gin"
)

// GoMiddleware represents the middleware handler
type GoMiddleware struct {
	// Some other fields that may be needed by middleware
}

// Log handles the logging middleware
func (m *GoMiddleware) Log() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()

		c.Next()

		end := time.Now()
		latency := end.Sub(start)

		msg := fmt.Sprintf("[%s] %d %s %s%s %s", end.Format(time.RFC3339), c.Writer.Status(), c.Request.Method, c.Request.Host, c.Request.URL.Path, latency.String())
		fmt.Println(msg)
	}
}

// JwtAuthMiddleware handles the JWT authentication middleware
func (m *GoMiddleware) JwtAuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		err := utils.TokenValid(c)
		if err != nil {
			c.String(http.StatusUnauthorized, "Unauthorized")
			c.Abort()
			return
		}
		c.Next()
	}
}

// InitMiddleware initializes the middleware
func InitMiddleware() *GoMiddleware {
	return &GoMiddleware{}
}
