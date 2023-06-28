package domain

import (
	"context"
	"warunk-bem/domain/dtos"
)

type Login struct {
	Username string `json:"username" validate:"required"`
	Password string `json:"password" validate:"required"`
}

type LoginUsecase interface {
	GetUser(ctx context.Context, req *dtos.LoginUserRequest) (*User, error)
}
