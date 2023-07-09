package http

import (
	"context"
	"net/http"
	"warunk-bem/domain"
	"warunk-bem/dtos"
	"warunk-bem/middlewares"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

type UserAmountHandler struct {
	UserAmountUsecase domain.UserAmountUsecase
}

func NewUserAmountHandler(protectedAdmin *gin.RouterGroup, uu domain.UserAmountUsecase) {
	handler := &UserAmountHandler{
		UserAmountUsecase: uu,
	}

	protectedAdmin = protectedAdmin.Group("/topup")
	protectedAdmin.POST("", handler.TopUpSaldo)
}

func isRequestValid(m *dtos.TopUpSaldoRequest) (bool, error) {
	validate := validator.New()
	err := validate.Struct(m)
	if err != nil {
		return false, err
	}
	return true, nil
}

func (uas *UserAmountHandler) TopUpSaldo(c *gin.Context) {
	var (
		useramount *dtos.TopUpSaldoRequest
		err        error
	)

	_, err = middlewares.IsAdmin(c)
	if err != nil {
		c.JSON(
			http.StatusForbidden,
			dtos.NewErrorResponse(
				http.StatusForbidden,
				"Only admin can top up saldo user",
				dtos.GetErrorData(err),
			),
		)
		return
	}

	err = c.ShouldBindJSON(&useramount)
	if err != nil {
		c.JSON(
			http.StatusUnprocessableEntity,
			dtos.NewErrorResponse(
				http.StatusUnprocessableEntity,
				"Field cannot be empty",
				dtos.GetErrorData(err),
			),
		)
		return
	}

	if ok, err := isRequestValid(useramount); !ok {
		c.JSON(
			http.StatusBadRequest,
			dtos.NewErrorResponse(
				http.StatusBadRequest,
				"Invalid request",
				dtos.GetErrorData(err),
			),
		)
		return
	}

	ctx := c.Request.Context()
	if ctx == nil {
		ctx = context.Background()
	}

	req, err := uas.UserAmountUsecase.TopUpSaldo(ctx, useramount)
	if err != nil {
		c.JSON(
			http.StatusBadRequest,
			dtos.NewErrorResponse(
				http.StatusBadRequest,
				"Cannot top up saldo user",
				dtos.GetErrorData(err),
			),
		)
		return
	}

	c.JSON(
		http.StatusOK,
		dtos.NewResponse(
			http.StatusOK,
			"Success top up saldo user",
			req,
		),
	)
}
