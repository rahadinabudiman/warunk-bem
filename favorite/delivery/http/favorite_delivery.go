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

type FavoriteHandler struct {
	FavoriteUsecase domain.FavoriteUsecase
	ProdukUsecase   domain.ProdukUsecase
}

func NewFavoriteHandler(protected *gin.RouterGroup, protectedAdmin *gin.RouterGroup, fu domain.FavoriteUsecase, pu domain.ProdukUsecase) {
	handler := &FavoriteHandler{
		FavoriteUsecase: fu,
		ProdukUsecase:   pu,
	}

	protected = protected.Group("/favorite")

	protected.POST("", handler.InsertOne)
}

func isRequestValid(m *domain.InsertFavoriteRequest) (bool, error) {
	validate := validator.New()
	err := validate.Struct(m)
	if err != nil {
		return false, err
	}
	return true, nil
}

func (fd *FavoriteHandler) InsertOne(c *gin.Context) {
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

	var favorite domain.InsertFavoriteRequest
	err = c.ShouldBindJSON(&favorite)
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

	if ok, err := isRequestValid(&favorite); !ok {
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

	check, _ := fd.FavoriteUsecase.FindOne(c, IDUser)
	if check == nil {
		favorite.UserID = IDUser
		res, err := fd.FavoriteUsecase.InsertOne(c, &favorite)
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
		produkBaru, err := fd.ProdukUsecase.FindOne(c, favorite.ProdukID)
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

		produkID, err := primitive.ObjectIDFromHex(favorite.ProdukID)
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
		FavoriteBaruBanget := &domain.FavoriteProduk{
			UserID: UserID,
			Produk: []domain.Produk{produkBaruLagi},
		}

		res, err := fd.FavoriteUsecase.UpdateOne(c, IDUser, FavoriteBaruBanget)
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
