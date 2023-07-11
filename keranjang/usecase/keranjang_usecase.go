package usecase

import (
	"context"
	"time"
	"warunk-bem/domain"
	"warunk-bem/dtos"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type KeranjangUsecase struct {
	KeranjangRepo  domain.KeranjangRepository
	ProdukRepo     domain.ProdukRepository
	UserRepo       domain.UserRepository
	contextTimeout time.Duration
}

func NewKeranjangUsecase(KeranjangRepo domain.KeranjangRepository, ProdukRepo domain.ProdukRepository, UserRepo domain.UserRepository, contextTimeout time.Duration) domain.KeranjangUsecase {
	return &KeranjangUsecase{
		KeranjangRepo:  KeranjangRepo,
		ProdukRepo:     ProdukRepo,
		UserRepo:       UserRepo,
		contextTimeout: contextTimeout,
	}
}

func (ku *KeranjangUsecase) InsertOne(ctx context.Context, req *dtos.InsertKeranjangRequest) (*dtos.InsertKeranjangResponse, error) {
	var res *dtos.InsertKeranjangResponse

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
				Stock:    produk.Stock,
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
