package middlewares

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

// CreateCookie creates a JWT cookie
func CreateCookie(c *gin.Context, token string) {
	cookie := &http.Cookie{
		Name:     "Warunk-BEM",
		Value:    token,
		Expires:  time.Now().Add(1 * time.Hour),
		Path:     "/",
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteStrictMode,
	}
	c.SetCookie(cookie.Name, cookie.Value, cookie.MaxAge, cookie.Path, cookie.Domain, cookie.Secure, cookie.HttpOnly)
}

// DeleteCookie deletes the JWT cookie
func DeleteCookie(c *gin.Context) error {
	cookie := &http.Cookie{
		Name:     "Warunk-BEM",
		Value:    "",
		Expires:  time.Now().Add(-1 * time.Hour),
		Path:     "/",
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteStrictMode,
	}
	c.SetCookie(cookie.Name, cookie.Value, cookie.MaxAge, cookie.Path, cookie.Domain, cookie.Secure, cookie.HttpOnly)
	c.Status(http.StatusNoContent)
	return nil
}
