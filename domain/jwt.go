package domain

import (
	"github.com/gin-gonic/gin"
)

type Jwt struct{}

type JwtUsecase interface {
	SetJwtAdmin(g *gin.RouterGroup)
	SetJwtUser(g *gin.RouterGroup)
	SetJwtGeneral(g *gin.RouterGroup)
}
