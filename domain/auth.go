package domain

import (
	"context"
	"warunk-bem/domain/dtos"

	"github.com/labstack/echo"
)

type Login struct {
	Username string `json:"username" validate:"required"`
	Password string `json:"password" validate:"required"`
}

type AuthUsecase interface {
	LoginUser(c echo.Context, ctx context.Context, req *dtos.LoginUserRequest) (*dtos.LoginUserResponse, error)
	LogoutUser(c echo.Context) (res *dtos.LogoutUserResponse, err error)
}
