package usecase

import (
	"context"
	"time"
	"warunk-bem/domain"

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

func (ku *KeranjangUsecase) FindOne(ctx context.Context, id string) (*domain.InsertKeranjangResponse, error) {
	ctx, cancel := context.WithTimeout(ctx, ku.contextTimeout)
	defer cancel()

	keranjang, err := ku.KeranjangRepo.FindOne(ctx, id)
	if err != nil {
		return nil, err
	}

	res := &domain.InsertKeranjangResponse{
		UserID: keranjang.UserID.Hex(),
		Produk: keranjang.Produk,
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
