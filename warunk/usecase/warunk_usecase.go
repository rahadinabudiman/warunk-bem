package usecase

import (
	"context"
	"errors"
	"time"
	"warunk-bem/domain"

	"github.com/go-redis/redis/v8"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type WarunkUsecase struct {
	WarunkRepo     domain.WarunkRepository
	ProdukRepo     domain.ProdukRepository
	UserRepo       domain.UserRepository
	RedisClient    *redis.Client
	contextTimeout time.Duration
}

func NewWarunkUsecase(WarunkRepo domain.WarunkRepository, ProdukRepo domain.ProdukRepository, UserRepo domain.UserRepository, RedisClient *redis.Client, contextTimeout time.Duration) domain.WarunkUsecase {
	return &WarunkUsecase{
		WarunkRepo:     WarunkRepo,
		ProdukRepo:     ProdukRepo,
		UserRepo:       UserRepo,
		RedisClient:    RedisClient,
		contextTimeout: contextTimeout,
	}
}

func (fu *WarunkUsecase) InsertOne(ctx context.Context, req *domain.InsertWarunkRequest) (*domain.InsertWarunkResponse, error) {
	var res *domain.InsertWarunkResponse

	ctx, cancel := context.WithTimeout(ctx, fu.contextTimeout)
	defer cancel()

	findLatest, err := fu.WarunkRepo.FindLatestByStatus(ctx, req.Status)
	if err != nil {
		return nil, err
	}

	if findLatest != nil {
		tanggalPembuatan := findLatest.CreatedAt
		StatusPembuatan := findLatest.Status
		formattedTanggal := tanggalPembuatan.Format("2006-01-02")

		user, err := fu.UserRepo.FindOne(ctx, req.UserID)
		if err != nil {
			return nil, errors.New("user not found")
		}

		// Create an array of domain.Produk from the req.Produk
		produks := make([]domain.Produk, 0, len(req.Produk))
		for _, v := range req.Produk {
			produk, err := fu.ProdukRepo.FindOne(ctx, v.ID.Hex())
			if err != nil {
				return nil, errors.New("produk not found")
			}

			// Create a new instance of domain.Produk and set its properties
			newProduk := domain.Produk{
				ID:       produk.ID,
				Slug:     produk.Slug,
				Name:     produk.Name,
				Detail:   produk.Detail,
				Price:    produk.Price,
				Stock:    v.Stock,
				Image:    produk.Image,
				Category: produk.Category,
			}

			produks = append(produks, newProduk)

			if req.Status == "Buka" {
				if v.ID == produk.ID {
					produk.Stock = v.Stock
					_, err := fu.ProdukRepo.UpdateOne(ctx, produk, produk.ID.Hex())
					if err != nil {
						return nil, errors.New("cannot update stock produk")
					}
				}
			}
		}

		req.ID = primitive.NewObjectID()
		req.CreatedAt = time.Now()
		req.UpdatedAt = time.Now()

		tanggalBuat := req.CreatedAt
		FormatTanggalBuat := tanggalBuat.Format("2006-01-02")
		if FormatTanggalBuat == formattedTanggal && req.Status == StatusPembuatan {
			return nil, errors.New("warunk already open")
		}

		_, err = fu.WarunkRepo.InsertOne(ctx, &domain.Warunk{
			ID:        req.ID,
			CreatedAt: req.CreatedAt,
			UpdatedAt: req.UpdatedAt,
			UserID:    user.ID,
			Produk:    produks,
			Status:    req.Status,
		})

		if err != nil {
			return nil, errors.New("cannot add produk to Warunk")
		}

		res = &domain.InsertWarunkResponse{
			ID:     req.ID.Hex(),
			UserID: req.UserID,
			Produk: produks,
			Status: req.Status,
		}

		return res, nil
	} else {
		user, err := fu.UserRepo.FindOne(ctx, req.UserID)
		if err != nil {
			return nil, errors.New("user not found")
		}

		// Create an array of domain.Produk from the req.Produk
		produks := make([]domain.Produk, 0, len(req.Produk))
		for _, v := range req.Produk {
			produk, err := fu.ProdukRepo.FindOne(ctx, v.ID.Hex())
			if err != nil {
				return nil, errors.New("produk not found")
			}

			// Create a new instance of domain.Produk and set its properties
			newProduk := domain.Produk{
				ID:       produk.ID,
				Slug:     produk.Slug,
				Name:     produk.Name,
				Detail:   produk.Detail,
				Price:    produk.Price,
				Stock:    v.Stock,
				Image:    produk.Image,
				Category: produk.Category,
			}

			produks = append(produks, newProduk)

			if req.Status == "Buka" {
				if v.ID == produk.ID {
					produk.Stock = v.Stock
					_, err := fu.ProdukRepo.UpdateOne(ctx, produk, produk.ID.Hex())
					if err != nil {
						return nil, errors.New("cannot update stock produk")
					}
				}
			}
		}

		req.ID = primitive.NewObjectID()
		req.CreatedAt = time.Now()
		req.UpdatedAt = time.Now()
		_, err = fu.WarunkRepo.InsertOne(ctx, &domain.Warunk{
			ID:        req.ID,
			CreatedAt: req.CreatedAt,
			UpdatedAt: req.UpdatedAt,
			UserID:    user.ID,
			Produk:    produks,
			Status:    req.Status,
		})

		if err != nil {
			return nil, errors.New("cannot add produk to Warunk")
		}

		res = &domain.InsertWarunkResponse{
			ID:     req.ID.Hex(),
			UserID: req.UserID,
			Produk: produks,
			Status: req.Status,
		}

		return res, nil
	}
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
