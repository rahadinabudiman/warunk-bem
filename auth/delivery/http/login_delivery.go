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

type LoginHandler struct {
	LoginUsecase domain.LoginUsecase
	Config       *viper.Viper
}

func NewLoginHandler(e *echo.Echo, lu domain.LoginUsecase, config *viper.Viper) {
	handler := &LoginHandler{
		LoginUsecase: lu,
		Config:       config,
	}
	// Main API
	api := e.Group("/api/v1")
	login := api.Group("/login")

	login.POST("", handler.CreateJwtUser)
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

func (login *LoginHandler) CreateJwtUser(c echo.Context) error {

	var (
		err          error
		loginPayload domain.Login
	)

	err = c.Bind(&loginPayload)
	if err != nil {
		return c.JSON(http.StatusUnprocessableEntity, err.Error())
	}

	if ok, err := isRequestValid(&loginPayload); !ok {
		return c.JSON(http.StatusBadRequest, err.Error())
	}

	ctx := c.Request().Context()
	if ctx == nil {
		ctx = context.Background()
	}

	loginRequest := dtos.LoginUserRequest{
		Username: loginPayload.Username,
		Password: loginPayload.Password,
	}

	res, err := login.LoginUsecase.GetUser(c, ctx, &loginRequest)
	if err != nil {
		return c.String(http.StatusUnauthorized, "Your username or password were wrong")
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
