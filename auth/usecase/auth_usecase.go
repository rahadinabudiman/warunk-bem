package usecase

import (
	"context"
	"errors"
	"strconv"
	"time"
	"warunk-bem/auth/middlewares"
	"warunk-bem/domain"
	"warunk-bem/domain/dtos"
	"warunk-bem/user/usecase/helpers"

	"github.com/labstack/echo"
	"github.com/spf13/viper"
	"golang.org/x/crypto/bcrypt"
)

type authUsecase struct {
	UserRepository domain.UserRepository
	contextTimeout time.Duration
	Config         *viper.Viper
}

func NewAuthUsecase(us domain.UserRepository, t time.Duration, config *viper.Viper) domain.AuthUsecase {
	return &authUsecase{
		UserRepository: us,
		contextTimeout: t,
		Config:         config,
	}
}

func (u *authUsecase) LoginUser(c echo.Context, ctx context.Context, req *dtos.LoginUserRequest) (*dtos.LoginUserResponse, error) {
	var res *dtos.LoginUserResponse
	ctx, cancel := context.WithTimeout(ctx, u.contextTimeout)
	defer cancel()

	user, err := u.UserRepository.FindUsername(ctx, req.Username)
	if err != nil {
		return nil, err
	}

	if !user.Verified {
		return nil, errors.New("please verify your account first")
	}

	err = helpers.ComparePassword(req.Password, user.Password)
	if err != nil && err == bcrypt.ErrMismatchedHashAndPassword {
		return nil, errors.New("username or password is incorrect")
	}

	lifetime, err := strconv.ParseInt(u.Config.GetString("LIFETIME"), 10, 64)
	if err != nil {
		lifetime = 60
	}

	Role := user.Role

	secret := u.Config.GetString("SECRET_JWT")
	token, err := middlewares.CreateJwtToken(user.ID.Hex(), Role, lifetime, secret)
	if err != nil {
		return res, errors.New("something went wrong")
	}

	req = &dtos.LoginUserRequest{
		Username: req.Username,
		Password: user.Password,
	}

	credential, err := u.UserRepository.GetByCredential(ctx, req)
	if err != nil {
		return nil, err
	}

	middlewares.CreateCookie(c, token)

	res = &dtos.LoginUserResponse{
		Username: credential.Username,
		Token:    token,
	}

	return res, nil
}

func (u *authUsecase) LogoutUser(c echo.Context) (res *dtos.LogoutUserResponse, err error) {
	err = middlewares.DeleteCookie(c)
	if err != nil {
		return nil, errors.New("cannot logout")
	}

	return res, err
}