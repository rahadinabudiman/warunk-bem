package usecase

import (
	"context"
	"errors"
	"time"
	"warunk-bem/domain"

	"github.com/spf13/viper"
)

type JwtUsecase struct {
	UserRepo       domain.UserRepository
	ContextTimeout time.Duration
	Config         *viper.Viper
}

func NewJwtUsecase(u domain.UserRepository, to time.Duration, config *viper.Viper) domain.JwtUsecase {
	return &JwtUsecase{
		UserRepo:       u,
		ContextTimeout: to,
		Config:         config,
	}
}

func (h *JwtUsecase) getOneUser(c context.Context, id string) (*domain.User, error) {

	ctx, cancel := context.WithTimeout(c, h.ContextTimeout)
	defer cancel()

	res, err := h.UserRepo.FindOne(ctx, id)
	if err != nil {
		return res, errors.New("couldn't find user")
	}

	return res, nil
}
