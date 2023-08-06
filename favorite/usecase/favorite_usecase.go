package usecase

import (
	"context"
	"encoding/json"
	"errors"
	"time"
	"warunk-bem/domain"

	"github.com/go-redis/redis/v8"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type FavoriteUsecase struct {
	FavoriteRepo   domain.FavoriteRepository
	ProdukRepo     domain.ProdukRepository
	UserRepo       domain.UserRepository
	RedisClient    *redis.Client
	contextTimeout time.Duration
}

func NewFavoriteUsecase(FavoriteRepo domain.FavoriteRepository, ProdukRepo domain.ProdukRepository, UserRepo domain.UserRepository, RedisClient *redis.Client, contextTimeout time.Duration) domain.FavoriteUsecase {
	return &FavoriteUsecase{
		FavoriteRepo:   FavoriteRepo,
		ProdukRepo:     ProdukRepo,
		UserRepo:       UserRepo,
		RedisClient:    RedisClient,
		contextTimeout: contextTimeout,
	}
}

// AddFavorite godoc
// @Summary      Add Favorite
// @Description  Add Favorite
// @Tags         User - Favorite
// @Accept       json
// @Produce      json
// @Param        request body dtos.InsertFavoriteRequest true "Payload Body [RAW]"
// @Success      201 {object} dtos.FavoriteCreatedResponse
// @Failure      400 {object} dtos.BadRequestResponse
// @Failure      401 {object} dtos.UnauthorizedResponse
// @Failure      403 {object} dtos.ForbiddenResponse
// @Failure      404 {object} dtos.NotFoundResponse
// @Failure      500 {object} dtos.InternalServerErrorResponse
// @Router       /favorite [post]
// @Security BearerAuth
func (fu *FavoriteUsecase) InsertOne(ctx context.Context, req *domain.InsertFavoriteRequest) (*domain.InsertFavoriteResponse, error) {
	var res *domain.InsertFavoriteResponse

	ctx, cancel := context.WithTimeout(ctx, fu.contextTimeout)
	defer cancel()

	produk, err := fu.ProdukRepo.FindOne(ctx, req.ProdukID)
	if err != nil {
		return nil, errors.New("produk not found")
	}

	user, err := fu.UserRepo.FindOne(ctx, req.UserID)
	if err != nil {
		return nil, errors.New("user not found")
	}

	req.ID = primitive.NewObjectID()
	req.CreatedAt = time.Now()
	req.UpdatedAt = time.Now()

	_, err = fu.FavoriteRepo.InsertOne(ctx, &domain.FavoriteProduk{
		ID:        req.ID,
		CreatedAt: req.CreatedAt,
		UpdatedAt: req.UpdatedAt,
		UserID:    user.ID,
		Produk: []domain.Produk{
			{
				ID:       produk.ID,
				Slug:     produk.Slug,
				Name:     produk.Name,
				Price:    produk.Price,
				Image:    produk.Image,
				Category: produk.Category,
			},
		},
	})

	if err != nil {
		return nil, errors.New("cannot add produk to favorite")
	}

	cacheKey := "favorite:" + req.UserID
	cacheValue, err := json.Marshal(res)
	if err == nil {
		fu.RedisClient.Set(ctx, cacheKey, cacheValue, 10*time.Minute)
	}

	return res, nil
}

// GetFavorite godoc
// @Summary      Get Favorite by UserID
// @Description  Get Favorite by UserID
// @Tags         User - Favorite
// @Accept       json
// @Produce      json
// @Success      200 {object} dtos.FavoriteOKResponse
// @Failure      400 {object} dtos.BadRequestResponse
// @Failure      401 {object} dtos.UnauthorizedResponse
// @Failure      403 {object} dtos.ForbiddenResponse
// @Failure      404 {object} dtos.NotFoundResponse
// @Failure      500 {object} dtos.InternalServerErrorResponse
// @Router       /favorite [get]
// @Security BearerAuth
func (fu *FavoriteUsecase) FindOne(ctx context.Context, id string) (*domain.InsertFavoriteResponse, error) {
	var res *domain.InsertFavoriteResponse

	cacheKey := "favorite:" + id
	cacheValue, err := fu.RedisClient.Get(ctx, cacheKey).Result()
	if err == nil {
		res = &domain.InsertFavoriteResponse{}
		if err := json.Unmarshal([]byte(cacheValue), res); err != nil {
			return nil, errors.New("cannot unmarshal cache value")
		}
		return res, nil
	}

	ctx, cancel := context.WithTimeout(ctx, fu.contextTimeout)
	defer cancel()

	favorite, err := fu.FavoriteRepo.FindOne(ctx, id)
	if err != nil {
		return nil, errors.New("favorite not found")
	}

	res = &domain.InsertFavoriteResponse{
		ID:     favorite.ID.Hex(),
		UserID: favorite.UserID.Hex(),
		Produk: favorite.Produk,
	}

	cacheVal, err := json.Marshal(res)
	if err == nil {
		fu.RedisClient.Set(ctx, cacheKey, cacheVal, 10*time.Minute)
	}

	return res, nil
}

func (fu *FavoriteUsecase) UpdateOne(ctx context.Context, id string, req *domain.FavoriteProduk) (*domain.FavoriteProduk, error) {
	ctx, cancel := context.WithTimeout(ctx, fu.contextTimeout)
	defer cancel()

	result, err := fu.FavoriteRepo.FindOne(ctx, id)
	if err != nil {
		return nil, errors.New("favorite not found")
	}

	// Jika produk ada di favorite, maka tampilkan message
	for _, v := range result.Produk {
		if v.ID == req.Produk[0].ID {
			return nil, errors.New("produk already in favorite")
		}
	}

	result.Produk = append(result.Produk, req.Produk...)

	_, err = fu.FavoriteRepo.UpdateOne(ctx, result, id)
	if err != nil {
		return nil, err
	}

	return result, nil
}

// RemoveProduct godoc
// @Summary      Remove Product from Favorite
// @Description  Remove Product from Favorite
// @Tags         User - Favorite
// @Accept       json
// @Produce      json
// @Param        request body dtos.DeleteFavoriteRequest true "Payload Body [RAW]"
// @Success      200 {object} dtos.DeleteProductFavoriteResponse
// @Failure      400 {object} dtos.BadRequestResponse
// @Failure      401 {object} dtos.UnauthorizedResponse
// @Failure      403 {object} dtos.ForbiddenResponse
// @Failure      404 {object} dtos.NotFoundResponse
// @Failure      500 {object} dtos.InternalServerErrorResponse
// @Router       /favorite/deleteproduct [post]
// @Security BearerAuth
func (fu *FavoriteUsecase) RemoveProduct(ctx context.Context, favoriteID string, productID string) (*domain.DeleteProductFavoriteResponse, error) {
	ctx, cancel := context.WithTimeout(ctx, fu.contextTimeout)
	defer cancel()

	// Check Favorite apakah sudah ada atau belum
	result, err := fu.FavoriteRepo.FindOneFavorite(ctx, favoriteID)
	if err != nil {
		return nil, errors.New("favorite not found")
	}

	// Check Produk apakah ada atau tidak
	produk, err := fu.ProdukRepo.FindOne(ctx, productID)
	if err != nil {
		return nil, errors.New("produk not found")
	}

	// Check apakah produk sudah ada di favorite atau belum
	var index = -1
	for i, v := range result.Produk {
		if v.ID.Hex() == produk.ID.Hex() {
			index = i
			break
		}
	}

	// Jika produk tidak ada di favorite
	if index == -1 {
		return nil, errors.New("product is not in favorite")
	}

	err = fu.FavoriteRepo.RemoveProduct(ctx, favoriteID, productID)
	if err != nil {
		return nil, errors.New("failed to remove produk")
	}

	check, err := fu.FavoriteRepo.FindOneFavorite(ctx, favoriteID)
	if err != nil {
		return nil, errors.New("favorite id not found")
	}

	if len(check.Produk) == 0 {
		err = fu.FavoriteRepo.DeleteOne(ctx, favoriteID)
		if err != nil {
			return nil, errors.New("cannot delete favorite")
		}
	}

	res := &domain.DeleteProductFavoriteResponse{
		Name: produk.Name,
	}

	return res, nil
}
