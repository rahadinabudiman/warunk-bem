package http

import (
	"net/http"
	"warunk-bem/domain"
	"warunk-bem/dtos"
	"warunk-bem/middlewares"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type TransaksiHandler struct {
	TransaksiUsecase domain.TransaksiUsecase
}

func NewUserHandler(protected *gin.RouterGroup, protectedAdmin *gin.RouterGroup, uu domain.TransaksiUsecase) {
	handler := &TransaksiHandler{
		TransaksiUsecase: uu,
	}

	// Main API
	protected = protected.Group("/transaksi")
	protectedAdmin = protectedAdmin.Group("/transaksi")

	protected.POST("", handler.InsertOne)
	protected.POST("/keranjang", handler.InsertByKeranjang)
}

func isRequestValid(m *dtos.InsertTransaksiRequest) (bool, error) {
	validate := validator.New()
	err := validate.Struct(m)
	if err != nil {
		return false, err
	}
	return true, nil
}

func (tc *TransaksiHandler) InsertOne(c *gin.Context) {
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

	var (
		usr dtos.InsertTransaksiRequest
	)

	err = c.ShouldBindJSON(&usr)
	if err != nil {
		c.JSON(
			http.StatusUnprocessableEntity,
			dtos.NewErrorResponse(
				http.StatusUnprocessableEntity,
				"Filed Cannot Be Empty",
				dtos.GetErrorData(err),
			),
		)
		return
	}

	if ok, err := isRequestValid(&usr); !ok {
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

	objectID, err := primitive.ObjectIDFromHex(idUser)
	if err != nil {
		c.JSON(
			http.StatusInternalServerError,
			dtos.NewErrorResponse(
				http.StatusInternalServerError,
				"Cannot Convert ObjectID",
				dtos.GetErrorData(err),
			),
		)
		return
	}

	usr.UserID = objectID
	res, err := tc.TransaksiUsecase.InsertOne(c, &usr)
	if err != nil {
		c.JSON(
			http.StatusInternalServerError,
			dtos.NewErrorResponse(
				http.StatusInternalServerError,
				"Cannot Insert Transaksi",
				dtos.GetErrorData(err),
			),
		)
		return
	}

	c.JSON(
		http.StatusCreated,
		dtos.NewResponse(
			http.StatusCreated,
			"Transaksi Berhasil",
			res,
		),
	)
}

func (tc *TransaksiHandler) InsertByKeranjang(c *gin.Context) {
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

	var (
		req dtos.InsertTransaksiKeranjangRequest
	)

	err = c.ShouldBindJSON(&req)
	if err != nil {
		c.JSON(
			http.StatusUnprocessableEntity,
			dtos.NewErrorResponse(
				http.StatusUnprocessableEntity,
				"Filed Cannot Be Empty",
				dtos.GetErrorData(err),
			),
		)
		return
	}

	objectID, err := primitive.ObjectIDFromHex(idUser)
	if err != nil {
		c.JSON(
			http.StatusInternalServerError,
			dtos.NewErrorResponse(
				http.StatusInternalServerError,
				"Cannot Convert ObjectID",
				dtos.GetErrorData(err),
			),
		)
		return
	}

	req.UserID = objectID

	res, err := tc.TransaksiUsecase.InsertByKeranjang(c, &req)
	if err != nil {
		c.JSON(
			http.StatusInternalServerError,
			dtos.NewErrorResponse(
				http.StatusInternalServerError,
				"Cannot Insert Transaksi",
				dtos.GetErrorData(err),
			),
		)
		return
	}

	c.JSON(
		http.StatusCreated,
		dtos.NewResponse(
			http.StatusCreated,
			"Transaksi Berhasil",
			res,
		),
	)
}
