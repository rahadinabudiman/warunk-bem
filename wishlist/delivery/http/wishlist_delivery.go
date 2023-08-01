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

type WishlistHandler struct {
	WishlistUsecase domain.WishlistUsecase
	ProdukUsecase   domain.ProdukUsecase
}

func NewWishlistHandler(protected *gin.RouterGroup, protectedAdmin *gin.RouterGroup, fu domain.WishlistUsecase, pu domain.ProdukUsecase) {
	handler := &WishlistHandler{
		WishlistUsecase: fu,
		ProdukUsecase:   pu,
	}

	protected = protected.Group("/wishlist")

	protected.POST("", handler.InsertOne)
}

func isRequestValid(m *domain.InsertWishlistRequest) (bool, error) {
	validate := validator.New()
	err := validate.Struct(m)
	if err != nil {
		return false, err
	}
	return true, nil
}

func (fd *WishlistHandler) InsertOne(c *gin.Context) {
	IDUser, err := middlewares.IsUser(c)
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

	var Wishlist domain.InsertWishlistRequest
	err = c.ShouldBindJSON(&Wishlist)
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

	if ok, err := isRequestValid(&Wishlist); !ok {
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

	check, _ := fd.WishlistUsecase.FindOne(c, IDUser)
	if check == nil {
		Wishlist.UserID = IDUser
		res, err := fd.WishlistUsecase.InsertOne(c, &Wishlist)
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
	} else {
		produkBaru, err := fd.ProdukUsecase.FindOne(c, Wishlist.ProdukID)
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

		produkID, err := primitive.ObjectIDFromHex(Wishlist.ProdukID)
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

		produkBaruLagi := domain.Produk{
			ID:       produkID,
			Slug:     produkBaru.Slug,
			Name:     produkBaru.Name,
			Price:    produkBaru.Price,
			Image:    produkBaru.Image,
			Category: produkBaru.Category,
		}

		UserID, err := primitive.ObjectIDFromHex(check.UserID)
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
		WishlistBaruBanget := &domain.WishlistProduk{
			UserID: UserID,
			Produk: []domain.Produk{produkBaruLagi},
		}

		res, err := fd.WishlistUsecase.UpdateOne(c, IDUser, WishlistBaruBanget)
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
}
