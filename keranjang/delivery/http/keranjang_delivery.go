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

type KeranjangHandler struct {
	KeranjangUsecase domain.KeranjangUsecase
	ProdukUsecase    domain.ProdukUsecase
}

func NewKeranjangHandler(protected *gin.RouterGroup, protectedAdmin *gin.RouterGroup, uu domain.KeranjangUsecase, pu domain.ProdukUsecase) {
	handler := &KeranjangHandler{
		KeranjangUsecase: uu,
		ProdukUsecase:    pu,
	}

	protected = protected.Group("/keranjang")
	// protectedAdmin = protectedAdmin.Group("/keranjang")

	protected.POST("", handler.InsertOne)
}

func isRequestValid(m *domain.InsertKeranjangRequest) (bool, error) {
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

	var keranjang domain.InsertKeranjangRequest
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

	check, _ := tc.KeranjangUsecase.FindOne(c, idUser)
	if check == nil {
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
		return
	}

	// keranjangExisting, err := tc.KeranjangUsecase.FindOne(c, idUser)
	// if err != nil {
	// 	c.JSON(
	// 		http.StatusBadRequest,
	// 		dtos.NewErrorResponse(
	// 			http.StatusBadRequest,
	// 			"Invalid request",
	// 			dtos.GetErrorData(err),
	// 		),
	// 	)
	// 	return
	// }

	ProdukBaru, err := tc.ProdukUsecase.FindOne(c, keranjang.ProdukID)
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
	ProdukID, err := primitive.ObjectIDFromHex(keranjang.ProdukID)
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
	ProdukBaruLagi := []domain.Produk{
		{
			ID:       ProdukID,
			Slug:     ProdukBaru.Slug,
			Name:     ProdukBaru.Name,
			Price:    ProdukBaru.Price,
			Stock:    int64(keranjang.Total),
			Image:    ProdukBaru.Image,
			Category: ProdukBaru.Category,
		},
	}

	// Update existing keranjang
	keranjangBaruBanget := &domain.Keranjang{
		Produk: ProdukBaruLagi,
	}

	res, err := tc.KeranjangUsecase.UpdateOne(c, idUser, keranjangBaruBanget)
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
