package usecase

import (
	"context"
	"errors"
	"time"
	"warunk-bem/domain"
	"warunk-bem/domain/dtos"
	"warunk-bem/user/usecase/helpers"

	"golang.org/x/crypto/bcrypt"
)

type loginUsecase struct {
	UserRepository domain.UserRepository
	contextTimeout time.Duration
}

func NewLoginUsecase(us domain.UserRepository, t time.Duration) domain.LoginUsecase {
	return &loginUsecase{
		UserRepository: us,
		contextTimeout: t,
	}
}

func (u *loginUsecase) GetUser(ctx context.Context, req *dtos.LoginUserRequest) (*domain.User, error) {
	ctx, cancel := context.WithTimeout(ctx, u.contextTimeout)
	defer cancel()

	user, err := u.UserRepository.FindUsername(ctx, req.Username)
	if err != nil {
		return nil, err
	}

	err = helpers.ComparePassword(req.Password, user.Password)
	if err != nil && err == bcrypt.ErrMismatchedHashAndPassword {
		return nil, errors.New("username or password is incorrect")

	}

	req = &dtos.LoginUserRequest{
		Username: req.Username,
		Password: user.Password,
	}

	res, err := u.UserRepository.GetByCredential(ctx, req)
	if err != nil {
		return nil, err
	}

	return res, nil
}
