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

type WarunkHandler struct {
	WarunkUsecase domain.WarunkUsecase
	ProdukUsecase domain.ProdukUsecase
}

func NewWarunkHandler(protectedAdmin *gin.RouterGroup, wu domain.WarunkUsecase, pu domain.ProdukUsecase) {
	handler := &WarunkHandler{
		WarunkUsecase: wu,
		ProdukUsecase: pu,
	}

	protectedAdmin = protectedAdmin.Group("/warunk")
	protectedAdmin.POST("", handler.InsertOne)
}

func isRequestValid(m *domain.InsertWarunkRequest) (bool, error) {
	validate := validator.New()
	err := validate.Struct(m)
	if err != nil {
		return false, err
	}
	return true, nil
}

func (fh *WarunkHandler) InsertOne(c *gin.Context) {
	IDAdmin, err := middlewares.IsAdmin(c)
	if err != nil {
		c.JSON(
			http.StatusUnauthorized,
			dtos.NewErrorResponse(
				http.StatusUnauthorized,
				"Please login first to access this pages",
				dtos.GetErrorData(err),
			),
		)
		return
	}

	var Warunk domain.InsertWarunkRequest
	err = c.ShouldBindJSON(&Warunk)
	if err != nil {
		c.JSON(
			http.StatusUnprocessableEntity,
			dtos.NewErrorResponse(
				http.StatusUnprocessableEntity,
				"Invalid Request",
				dtos.GetErrorData(err),
			),
		)
		return
	}

	if ok, err := isRequestValid(&Warunk); !ok {
		c.JSON(
			http.StatusBadRequest,
			dtos.NewErrorResponse(
				http.StatusBadRequest,
				"Invalid Request",
				dtos.GetErrorData(err),
			),
		)
		return
	}

	ctx := c.Request.Context()
	if ctx == nil {
		ctx = context.Background()
	}

	Warunk.UserID = IDAdmin
	res, err := fh.WarunkUsecase.InsertOne(ctx, &Warunk)
	if err != nil {
		c.JSON(
			http.StatusBadRequest,
			dtos.NewErrorResponse(
				http.StatusBadRequest,
				"Failed to insert data",
				dtos.GetErrorData(err),
			),
		)
		return
	}

	c.JSON(
		http.StatusCreated,
		dtos.NewResponse(
			http.StatusCreated,
			"Success insert data",
			res,
		),
	)
}
