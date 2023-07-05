package usecase

import (
	"context"
	"errors"
	"time"
	"warunk-bem/domain"
	"warunk-bem/dtos"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type TransaksiUsecase struct {
	TransaksiRepo  domain.TransaksiRepository
	ProdukRepo     domain.ProdukRepository
	UserRepo       domain.UserRepository
	contextTimeout time.Duration
}

func NewTransaksiUsecase(TransaksiRepo domain.TransaksiRepository, ProdukRepo domain.ProdukRepository, UserRepo domain.UserRepository, contextTimeout time.Duration) domain.TransaksiUsecase {
	return &TransaksiUsecase{
		TransaksiRepo:  TransaksiRepo,
		ProdukRepo:     ProdukRepo,
		UserRepo:       UserRepo,
		contextTimeout: contextTimeout,
	}
}

func (tu *TransaksiUsecase) InsertOne(ctx context.Context, req *dtos.InsertTransaksiRequest) (*dtos.InsertTransaksiResponse, error) {
	var res *dtos.InsertTransaksiResponse

	ctx, cancel := context.WithTimeout(ctx, tu.contextTimeout)
	defer cancel()

	produk, err := tu.ProdukRepo.FindOne(ctx, req.ProdukID.Hex())
	if err != nil {
		return res, err
	}

	user, err := tu.UserRepo.FindOne(ctx, req.UserID.Hex())
	if err != nil {
		return res, err
	}

	if produk.Stock == 0 {
		return nil, errors.New("produk telah habis terjual")
	}

	if produk.Stock < int64(req.Total) {
		return nil, errors.New("stok produk tidak mencukupi")
	}

	produk.Stock = produk.Stock - int64(req.Total)
	produk.UpdatedAt = time.Now()

	_, err = tu.ProdukRepo.UpdateOne(ctx, produk, produk.ID.Hex())
	if err != nil {
		return res, err
	}

	req.ID = primitive.NewObjectID()
	req.CreatedAt = time.Now()
	req.UpdatedAt = time.Now()

	transaksireq := &domain.Transaksi{
		ID:        req.ID,
		CreatedAt: req.CreatedAt,
		UpdatedAt: req.UpdatedAt,
		UserID:    req.UserID,
		ProdukID:  req.ProdukID,
		Total:     int64(req.Total),
		Status:    "Berhasil",
	}

	resp, err := tu.TransaksiRepo.InsertOne(ctx, transaksireq)
	if err != nil {
		return res, err
	}

	res = &dtos.InsertTransaksiResponse{
		Name:       user.Name,
		ProdukName: produk.Name,
		Total:      resp.Total,
	}

	return res, nil
}
