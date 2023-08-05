package usecase

import (
	"context"
	"errors"
	"time"
	"warunk-bem/domain"

	"github.com/go-redis/redis/v8"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type KeranjangUsecase struct {
	KeranjangRepo  domain.KeranjangRepository
	ProdukRepo     domain.ProdukRepository
	UserRepo       domain.UserRepository
	RedisClient    *redis.Client
	contextTimeout time.Duration
}

func NewKeranjangUsecase(KeranjangRepo domain.KeranjangRepository, ProdukRepo domain.ProdukRepository, UserRepo domain.UserRepository, RedisClient *redis.Client, contextTimeout time.Duration) domain.KeranjangUsecase {
	return &KeranjangUsecase{
		KeranjangRepo:  KeranjangRepo,
		ProdukRepo:     ProdukRepo,
		UserRepo:       UserRepo,
		RedisClient:    RedisClient,
		contextTimeout: contextTimeout,
	}
}

// AddKeranjang godoc
// @Summary      Add Keranjang
// @Description  Add Keranjang
// @Tags         User - Keranjang
// @Accept       json
// @Produce      json
// @Param        request body dtos.InsertKeranjangRequest true "Payload Body [RAW]"
// @Success      201 {object} dtos.DetailKeranjang
// @Failure      400 {object} dtos.BadRequestResponse
// @Failure      401 {object} dtos.UnauthorizedResponse
// @Failure      403 {object} dtos.ForbiddenResponse
// @Failure      404 {object} dtos.NotFoundResponse
// @Failure      500 {object} dtos.InternalServerErrorResponse
// @Router       /keranjang [post]
// @Security BearerAuth
func (ku *KeranjangUsecase) InsertOne(ctx context.Context, req *domain.InsertKeranjangRequest) (*domain.InsertKeranjangResponse, error) {
	var res *domain.InsertKeranjangResponse

	ctx, cancel := context.WithTimeout(ctx, ku.contextTimeout)
	defer cancel()

	produk, err := ku.ProdukRepo.FindOne(ctx, req.ProdukID)
	if err != nil {
		return res, err
	}

	user, err := ku.UserRepo.FindOne(ctx, req.UserID)
	if err != nil {
		return res, err
	}

	req.ID = primitive.NewObjectID()
	req.CreatedAt = time.Now()
	req.UpdatedAt = time.Now()

	_, err = ku.KeranjangRepo.InsertOne(ctx, &domain.Keranjang{
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
				Stock:    int64(req.Total),
				Image:    produk.Image,
				Category: produk.Category,
			},
		},
		Total: req.Total,
	})
	if err != nil {
		return res, err
	}

	return res, nil
}

// GetKeranjang godoc
// @Summary      Get Keranjang by UserID
// @Description  Get Keranjang by UserID
// @Tags         User - Keranjang
// @Accept       json
// @Produce      json
// @Success      200 {object} dtos.KeranjangOKResponse
// @Failure      400 {object} dtos.BadRequestResponse
// @Failure      401 {object} dtos.UnauthorizedResponse
// @Failure      403 {object} dtos.ForbiddenResponse
// @Failure      404 {object} dtos.NotFoundResponse
// @Failure      500 {object} dtos.InternalServerErrorResponse
// @Router       /keranjang [get]
// @Security BearerAuth
func (ku *KeranjangUsecase) FindOne(ctx context.Context, id string) (*domain.InsertKeranjangResponse, error) {
	ctx, cancel := context.WithTimeout(ctx, ku.contextTimeout)
	defer cancel()

	keranjang, err := ku.KeranjangRepo.FindOne(ctx, id)
	if err != nil {
		return nil, err
	}

	res := &domain.InsertKeranjangResponse{
		ID:     keranjang.ID.Hex(),
		UserID: keranjang.UserID.Hex(),
		Produk: keranjang.Produk,
	}

	return res, nil
}

// RemoveProduct godoc
// @Summary      Remove Product from Keranjang
// @Description  Remove Product from Keranjang
// @Tags         User - Keranjang
// @Accept       json
// @Produce      json
// @Param        request body dtos.DeleteProductKeranjangRequest true "Payload Body [RAW]"
// @Success      200 {object} dtos.DeleteProductKeranjangResponse
// @Failure      400 {object} dtos.BadRequestResponse
// @Failure      401 {object} dtos.UnauthorizedResponse
// @Failure      403 {object} dtos.ForbiddenResponse
// @Failure      404 {object} dtos.NotFoundResponse
// @Failure      500 {object} dtos.InternalServerErrorResponse
// @Router       /keranjang/deleteproduct [post]
// @Security BearerAuth
func (ku *KeranjangUsecase) RemoveProduct(ctx context.Context, keranjangID string, productID string) (*domain.DeleteProductKeranjangResponse, error) {
	ctx, cancel := context.WithTimeout(ctx, ku.contextTimeout)
	defer cancel()

	// Check Keranjang apakah ada atau tidak
	result, err := ku.KeranjangRepo.FindOneKeranjang(ctx, keranjangID)
	if err != nil {
		return nil, err
	}

	// Check Produk apakah ada atau tidak
	produk, err := ku.ProdukRepo.FindOne(ctx, productID)
	if err != nil {
		return nil, errors.New("produk tidak ditemukan")
	}

	// Check apakah produk ada di keranjang atau tidak
	var index = -1
	for i, v := range result.Produk {
		if v.ID.Hex() == produk.ID.Hex() {
			index = i
			break
		}
	}

	// Jika produk tidak ada di keranjang
	if index == -1 {
		return nil, errors.New("produk tidak ada di keranjang")
	}

	result.Total -= int(result.Produk[index].Stock)
	_, err = ku.KeranjangRepo.UpdateOneKeranjang(ctx, result, keranjangID)
	if err != nil {
		return nil, errors.New("tidak dapat mengupdate total produk")
	}

	err = ku.KeranjangRepo.RemoveProduct(ctx, keranjangID, productID)
	if err != nil {
		return nil, errors.New("produk tidak ditemukan")
	}

	check, err := ku.KeranjangRepo.FindOneKeranjang(ctx, keranjangID)
	if err != nil {
		return nil, errors.New("tidak ada keranjang dengan id tersebut")
	}

	if len(check.Produk) == 0 {
		err = ku.KeranjangRepo.DeleteOne(ctx, keranjangID)
		if err != nil {
			return nil, errors.New("tidak dapat menghapus keranjang")
		}
	}

	res := &domain.DeleteProductKeranjangResponse{
		Name: produk.Name,
	}

	return res, nil
}

func (ku *KeranjangUsecase) UpdateOne(ctx context.Context, id string, req *domain.Keranjang) (*domain.Keranjang, error) {
	ctx, cancel := context.WithTimeout(ctx, ku.contextTimeout)
	defer cancel()

	result, err := ku.KeranjangRepo.FindOne(ctx, id)
	if err != nil {
		return nil, err
	}

	// Jika produk sudah ada di keranjang, tambahkan stoknya saja jangan tambahkan array produknya
	for i, v := range result.Produk {
		if v.ID == req.Produk[0].ID {
			result.Produk[i].Stock += req.Produk[0].Stock
			result.Total += int(req.Produk[0].Stock)
			_, err = ku.KeranjangRepo.UpdateOne(ctx, result, id)
			if err != nil {
				return nil, err
			}
			return result, nil
		}
	}

	result.Produk = append(result.Produk, req.Produk...)
	result.Total += req.Total

	_, err = ku.KeranjangRepo.UpdateOne(ctx, result, id)
	if err != nil {
		return nil, err
	}

	return result, nil
}
