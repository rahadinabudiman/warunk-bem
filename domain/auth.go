package domain

import (
	"context"
	"warunk-bem/domain/dtos"

	"github.com/gin-gonic/gin"
)

type Login struct {
	Email    string `json:"email" validate:"required"`
	Password string `json:"password" validate:"required"`
}

type AuthUsecase interface {
	LoginUser(c *gin.Context, ctx context.Context, req *dtos.LoginUserRequest) (*dtos.LoginUserResponse, error)
	LogoutUser(c *gin.Context) (res *dtos.LogoutUserResponse, err error)
}
