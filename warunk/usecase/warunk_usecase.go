package usecase

import (
	"context"
	"errors"
	"time"
	"warunk-bem/domain"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type WarunkUsecase struct {
	WarunkRepo     domain.WarunkRepository
	ProdukRepo     domain.ProdukRepository
	UserRepo       domain.UserRepository
	contextTimeout time.Duration
}

func NewWarunkUsecase(WarunkRepo domain.WarunkRepository, ProdukRepo domain.ProdukRepository, UserRepo domain.UserRepository, contextTimeout time.Duration) domain.WarunkUsecase {
	return &WarunkUsecase{
		WarunkRepo:     WarunkRepo,
		ProdukRepo:     ProdukRepo,
		UserRepo:       UserRepo,
		contextTimeout: contextTimeout,
	}
}

func (fu *WarunkUsecase) InsertOne(ctx context.Context, req *domain.InsertWarunkRequest) (*domain.InsertWarunkResponse, error) {
	var res *domain.InsertWarunkResponse

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

	_, err = fu.WarunkRepo.InsertOne(ctx, &domain.Warunk{
		ID:        req.ID,
		CreatedAt: req.CreatedAt,
		UpdatedAt: req.UpdatedAt,
		UserID:    user.ID,
		Produk: []domain.Produk{
			{
				ID:       produk.ID,
				Slug:     produk.Slug,
				Name:     produk.Name,
				Detail:   produk.Detail,
				Price:    produk.Price,
				Stock:    produk.Stock,
				Image:    produk.Image,
				Category: produk.Category,
			},
		},
	})

	if err != nil {
		return nil, errors.New("cannot add produk to Warunk")
	}

	return res, nil
}

func (fu *WarunkUsecase) FindOne(ctx context.Context, id string) (*domain.InsertWarunkResponse, error) {
	ctx, cancel := context.WithTimeout(ctx, fu.contextTimeout)
	defer cancel()

	Warunk, err := fu.WarunkRepo.FindOne(ctx, id)
	if err != nil {
		return nil, errors.New("warunk not found")
	}

	res := &domain.InsertWarunkResponse{
		ID:     Warunk.ID.Hex(),
		UserID: Warunk.UserID.Hex(),
		Produk: Warunk.Produk,
	}

	return res, nil
}

func (fu *WarunkUsecase) UpdateOne(ctx context.Context, id string, req *domain.Warunk) (*domain.Warunk, error) {
	ctx, cancel := context.WithTimeout(ctx, fu.contextTimeout)
	defer cancel()

	result, err := fu.WarunkRepo.FindOne(ctx, id)
	if err != nil {
		return nil, errors.New("warunk not found")
	}

	// Jika produk ada di Warunk, maka tampilkan message
	for _, v := range result.Produk {
		if v.ID == req.Produk[0].ID {
			return nil, errors.New("produk already in Warunk")
		}
	}

	result.Produk = append(result.Produk, req.Produk...)

	_, err = fu.WarunkRepo.UpdateOne(ctx, result, id)
	if err != nil {
		return nil, err
	}

	return result, nil
}
