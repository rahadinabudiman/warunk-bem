package usecase

import (
	"context"
	"errors"
	"time"
	"warunk-bem/domain"
	"warunk-bem/domain/dtos"
	"warunk-bem/middlewares"
	"warunk-bem/user/usecase/helpers"
	"warunk-bem/utils"

	"github.com/gin-gonic/gin"
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

func (u *authUsecase) LoginUser(c *gin.Context, ctx context.Context, req *dtos.LoginUserRequest) (*dtos.LoginUserResponse, error) {
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

	// Role := user.Role

	token, err := utils.GenerateToken(user.ID.Hex())
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

func (u *authUsecase) LogoutUser(c *gin.Context) (res *dtos.LogoutUserResponse, err error) {
	err = middlewares.DeleteCookie(c)
	if err != nil {
		return nil, errors.New("cannot logout")
	}

	return res, err
}
