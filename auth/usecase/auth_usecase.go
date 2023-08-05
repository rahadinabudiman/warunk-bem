package usecase

import (
	"context"
	"errors"
	"strconv"
	"time"
	"warunk-bem/domain"
	"warunk-bem/dtos"
	"warunk-bem/helpers"
	"warunk-bem/middlewares"
	"warunk-bem/utils"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	"golang.org/x/crypto/bcrypt"
)

type authUsecase struct {
	UserRepository domain.UserRepository
	RedisClient    *redis.Client
	contextTimeout time.Duration
}

func NewAuthUsecase(us domain.UserRepository, RedisClient *redis.Client, t time.Duration) domain.AuthUsecase {
	return &authUsecase{
		UserRepository: us,
		RedisClient:    RedisClient,
		contextTimeout: t,
	}
}

// UserLogin godoc
// @Summary      Login User with Username and Password
// @Description  Login an account
// @Tags         User - Auth
// @Accept       json
// @Produce      json
// @Param        request body dtos.LoginUserRequest true "Payload Body [RAW]"
// @Success      200 {object} dtos.LoginStatusOKResponse
// @Failure      400 {object} dtos.BadRequestResponse
// @Failure      401 {object} dtos.UnauthorizedResponse
// @Failure      403 {object} dtos.ForbiddenResponse
// @Failure      404 {object} dtos.NotFoundResponse
// @Failure      500 {object} dtos.InternalServerErrorResponse
// @Router       /login [post]
func (u *authUsecase) LoginUser(c *gin.Context, ctx context.Context, req *dtos.LoginUserRequest) (*dtos.LoginUserResponse, error) {
	var res *dtos.LoginUserResponse
	ctx, cancel := context.WithTimeout(ctx, u.contextTimeout)
	defer cancel()

	user, err := u.UserRepository.FindEmail(ctx, req.Email)
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

	Role := user.Role

	token, err := utils.GenerateToken(user.ID.Hex(), Role)
	if err != nil {
		return res, errors.New("something went wrong")
	}

	req = &dtos.LoginUserRequest{
		Email:    req.Email,
		Password: user.Password,
	}

	credential, err := u.UserRepository.GetByCredential(ctx, req)
	if err != nil {
		return nil, err
	}

	middlewares.CreateCookie(c, token)

	otp := helpers.GenerateRandomOTP(6)
	NewOTP, err := strconv.Atoi(otp)
	if err != nil {
		return res, errors.New("failed to Generate OTP")
	}

	credential.LoginVerif = 1
	credential.VerificationCode = NewOTP
	credential.UpdatedAt = time.Now()

	_, err = u.UserRepository.UpdateOne(ctx, credential, credential.ID.Hex())
	if err != nil {
		return nil, err
	}

	emailData := utils.EmailData{
		Code:      NewOTP,
		FirstName: credential.Name,
		Subject:   "Your Verification Login Code",
	}

	utils.SendEmail(credential, &emailData)

	res = &dtos.LoginUserResponse{
		Username: credential.Username,
		Token:    token,
		Message:  "Please check your email for verification code",
	}

	return res, nil
}

// LogoutUser godoc
// @Summary      Logout User
// @Description  Logout User
// @Tags         User - Auth
// @Accept       json
// @Produce      json
// @Success      200 {object} dtos.LogoutUserResponse
// @Failure      400 {object} dtos.BadRequestResponse
// @Failure      401 {object} dtos.UnauthorizedResponse
// @Failure      403 {object} dtos.ForbiddenResponse
// @Failure      404 {object} dtos.NotFoundResponse
// @Failure      500 {object} dtos.InternalServerErrorResponse
// @Router       /user/logout [get]
// @Security BearerAuth
func (u *authUsecase) LogoutUser(c *gin.Context) (res *dtos.LogoutUserResponse, err error) {
	err = middlewares.DeleteCookie(c)
	if err != nil {
		return nil, errors.New("cannot logout")
	}

	return res, err
}
