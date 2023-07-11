package http

import (
	"net/http"
	"warunk-bem/domain"
	"warunk-bem/dtos"
	"warunk-bem/middlewares"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

type KeranjangHandler struct {
	KeranjangUsecase domain.KeranjangUsecase
}

func NewKeranjangHandler(protected *gin.RouterGroup, protectedAdmin *gin.RouterGroup, uu domain.KeranjangUsecase) {
	handler := &KeranjangHandler{
		KeranjangUsecase: uu,
	}

	protected = protected.Group("/keranjang")
	// protectedAdmin = protectedAdmin.Group("/keranjang")

	protected.POST("", handler.InsertOne)
}

func isRequestValid(m *dtos.InsertKeranjangRequest) (bool, error) {
	validate := validator.New()
	err := validate.Struct(m)
	if err != nil {
		return false, err
	}
	return true, nil
}

func (tc *KeranjangHandler) InsertOne(c *gin.Context) {
	idUser, err := middlewares.IsUser(c)
	if err != nil {
		c.JSON(
			http.StatusUnauthorized,
			dtos.NewErrorResponse(
				http.StatusUnauthorized,
				"Unauthorized",
				dtos.GetErrorData(err),
			),
		)
		return
	}

	var keranjang dtos.InsertKeranjangRequest
	err = c.ShouldBindJSON(&keranjang)
	if err != nil {
		c.JSON(
			http.StatusUnprocessableEntity,
			dtos.NewErrorResponse(
				http.StatusUnprocessableEntity,
				"Invalid request",
				dtos.GetErrorData(err),
			),
		)
		return
	}

	if ok, err := isRequestValid(&keranjang); !ok {
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

	keranjang.UserID = idUser
	res, err := tc.KeranjangUsecase.InsertOne(c, &keranjang)
	if err != nil {
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

	c.JSON(
		http.StatusCreated,
		dtos.NewResponse(
			http.StatusCreated,
			"Success",
			res,
		),
	)
}
