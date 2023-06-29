package http

import (
	"context"
	"net/http"
	"warunk-bem/domain"
	"warunk-bem/domain/dtos"
	"warunk-bem/user/usecase/helpers"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo"
	"github.com/spf13/viper"
)

type AuthHandler struct {
	AuthUsecase domain.AuthUsecase
	Config      *viper.Viper
}

func NewAuthHandler(api *echo.Group, generalJwt *echo.Group, lu domain.AuthUsecase, config *viper.Viper) {
	handler := &AuthHandler{
		AuthUsecase: lu,
		Config:      config,
	}

	api.POST("/login", handler.LoginUser)
	generalJwt.GET("/logout", handler.LogoutUser)
}

func isRequestValid(m *domain.Login) (bool, error) {
	validate := validator.New()
	cv := &helpers.CustomValidator{Validators: validate}
	err := cv.Validate(m)
	if err != nil {
		return false, err
	}
	return true, nil
}

func (delivery *AuthHandler) LoginUser(c echo.Context) error {
	var (
		err          error
		loginPayload domain.Login
	)

	err = c.Bind(&loginPayload)
	if err != nil {
		return c.JSON(http.StatusUnprocessableEntity, err.Error())
	}

	if ok, err := isRequestValid(&loginPayload); !ok {
		return c.JSON(
			http.StatusBadRequest,
			dtos.NewErrorResponse(
				http.StatusBadRequest,
				"username or password cannot be blank",
				err.Error(),
			),
		)
	}

	ctx := c.Request().Context()
	if ctx == nil {
		ctx = context.Background()
	}

	loginRequest := dtos.LoginUserRequest{
		Username: loginPayload.Username,
		Password: loginPayload.Password,
	}

	res, err := delivery.AuthUsecase.LoginUser(c, ctx, &loginRequest)
	if err != nil {
		return c.JSON(
			http.StatusUnauthorized,
			dtos.NewErrorResponse(
				http.StatusUnauthorized,
				"cannot login",
				err.Error(),
			),
		)
	}

	return c.JSON(
		http.StatusOK,
		dtos.NewResponse(
			http.StatusOK,
			"success",
			res,
		),
	)

}

func (delivery *AuthHandler) LogoutUser(c echo.Context) error {
	_, err := delivery.AuthUsecase.LogoutUser(c)
	if err != nil {
		return c.JSON(
			http.StatusBadRequest,
			dtos.NewErrorResponse(
				http.StatusBadRequest,
				"cannot logout",
				err.Error(),
			),
		)
	}

	return c.JSON(
		http.StatusOK,
		dtos.NewResponseMessage(
			http.StatusOK,
			"Success Logout",
		),
	)
}
