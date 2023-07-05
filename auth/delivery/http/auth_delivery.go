package http

import (
	"context"
	"net/http"
	"warunk-bem/domain"
	"warunk-bem/dtos"
	"warunk-bem/helpers"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/spf13/viper"
)

type AuthHandler struct {
	AuthUsecase domain.AuthUsecase
	Config      *viper.Viper
}

func NewAuthHandler(api *gin.RouterGroup, protected *gin.RouterGroup, lu domain.AuthUsecase, config *viper.Viper) {
	handler := &AuthHandler{
		AuthUsecase: lu,
		Config:      config,
	}

	api.POST("/login", handler.LoginUser)
	protected.GET("/logout", handler.LogoutUser)
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

func (delivery *AuthHandler) LoginUser(c *gin.Context) {
	var (
		err          error
		loginPayload domain.Login
	)

	err = c.ShouldBindJSON(&loginPayload)
	if err != nil {
		c.JSON(
			http.StatusUnprocessableEntity,
			dtos.NewErrorResponse(
				http.StatusUnprocessableEntity,
				"Field Cannot Be Empty",
				dtos.GetErrorData(err),
			),
		)
		return
	}

	if ok, err := isRequestValid(&loginPayload); !ok {
		c.JSON(
			http.StatusBadRequest,
			dtos.NewErrorResponse(
				http.StatusBadRequest,
				"Bad Request",
				dtos.GetErrorData(err),
			),
		)
		return
	}

	ctx := c.Request.Context()
	if ctx == nil {
		ctx = context.Background()
	}

	loginRequest := dtos.LoginUserRequest{
		Email:    loginPayload.Email,
		Password: loginPayload.Password,
	}

	res, err := delivery.AuthUsecase.LoginUser(c, ctx, &loginRequest)
	if err != nil {
		c.JSON(
			http.StatusUnauthorized,
			dtos.NewErrorResponse(
				http.StatusUnauthorized,
				"Cannot login",
				err.Error(),
			),
		)
		return
	}

	c.JSON(
		http.StatusOK,
		dtos.NewResponse(
			http.StatusOK,
			"success",
			res,
		),
	)
}

func (delivery *AuthHandler) LogoutUser(c *gin.Context) {
	_, err := delivery.AuthUsecase.LogoutUser(c)
	if err != nil {
		c.JSON(
			http.StatusBadRequest,
			dtos.NewErrorResponse(
				http.StatusBadRequest,
				"Cannot logout",
				err.Error(),
			),
		)
		return
	}

	c.JSON(
		http.StatusOK,
		dtos.NewResponseMessage(
			http.StatusOK,
			"Success Logout",
		),
	)
}
