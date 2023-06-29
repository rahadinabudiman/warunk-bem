package usecase

import (
	"context"
	"errors"
	"time"
	"warunk-bem/domain"
	"warunk-bem/domain/dtos"
	"warunk-bem/user/usecase/helpers"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"golang.org/x/crypto/bcrypt"
)

type produkUsecase struct {
	ProdukRepo     domain.ProdukRepository
	UserRepo       domain.UserRepository
	contextTimeout time.Duration
}

func NewProdukUsecase(ProdukRepo domain.ProdukRepository, UserRepo domain.UserRepository, contextTimeout time.Duration) domain.ProdukUsecase {
	return &produkUsecase{
		ProdukRepo:     ProdukRepo,
		UserRepo:       UserRepo,
		contextTimeout: contextTimeout,
	}
}

func (pu *produkUsecase) InsertOne(c context.Context, req *dtos.InsertProdukRequest) (*dtos.InsertProdukResponse, error) {
	var res *dtos.InsertProdukResponse

	ctx, cancel := context.WithTimeout(c, pu.contextTimeout)
	defer cancel()

	req.ID = primitive.NewObjectID()
	req.CreatedAt = time.Now()
	req.UpdatedAt = time.Now()

	CreateProduk := &domain.Produk{
		ID:        req.ID,
		CreatedAt: req.CreatedAt,
		UpdatedAt: req.UpdatedAt,
		Name:      req.Name,
		Detail:    req.Detail,
		Price:     req.Price,
		Stock:     req.Stock,
		Category:  req.Category,
	}

	createdProduk, err := pu.ProdukRepo.InsertOne(ctx, CreateProduk)
	if err != nil {
		return res, err
	}

	res = &dtos.InsertProdukResponse{
		Name:     createdProduk.Name,
		Detail:   createdProduk.Detail,
		Price:    createdProduk.Price,
		Stock:    createdProduk.Stock,
		Category: createdProduk.Category,
	}

	return res, nil
}
func (pu *produkUsecase) FindOne(c context.Context, id string) (res *dtos.ProdukDetailResponse, err error) {
	ctx, cancel := context.WithTimeout(c, pu.contextTimeout)
	defer cancel()

	req, err := pu.ProdukRepo.FindOne(ctx, id)
	if err != nil {
		return res, err
	}

	res = &dtos.ProdukDetailResponse{
		Name:     req.Name,
		Detail:   req.Detail,
		Price:    req.Price,
		Stock:    req.Stock,
		Category: req.Category,
	}

	return res, nil
}

func (pu *produkUsecase) GetAllWithPage(c context.Context, rp int64, p int64, filter interface{}, setsort interface{}) ([]dtos.ProdukDetailResponse, int64, error) {
	return nil, 0, nil
}

func (pu *produkUsecase) UpdateOne(c context.Context, req *dtos.ProdukUpdateRequest, id string) (*dtos.ProdukDetailResponse, error) {
	var res *dtos.ProdukDetailResponse

	ctx, cancel := context.WithTimeout(c, pu.contextTimeout)
	defer cancel()

	result, err := pu.ProdukRepo.FindOne(ctx, id)
	if err != nil {
		return res, err
	}

	result.Name = req.Name
	result.Detail = req.Detail
	result.Price = req.Price
	result.Stock = req.Stock
	result.Category = req.Category

	resp, err := pu.ProdukRepo.UpdateOne(ctx, result, id)
	if err != nil {
		return res, err
	}

	res = &dtos.ProdukDetailResponse{
		Name:     resp.Name,
		Detail:   resp.Detail,
		Price:    resp.Price,
		Stock:    resp.Stock,
		Category: resp.Category,
	}

	return res, nil
}

func (pu *produkUsecase) DeleteOne(c context.Context, id string, idAdmin string, req dtos.DeleteProdukRequest) (res dtos.ResponseMessage, err error) {
	ctx, cancel := context.WithTimeout(c, pu.contextTimeout)
	defer cancel()

	userAdmin, err := pu.UserRepo.FindOne(ctx, idAdmin)
	if err != nil {
		return res, err
	}

	err = helpers.ComparePassword(req.Password, userAdmin.Password)
	if err != nil && err == bcrypt.ErrMismatchedHashAndPassword {
		return res, errors.New("password is incorrect")
	}

	err = pu.ProdukRepo.DeleteOne(ctx, id)
	if err != nil {
		return res, err
	}

	return res, nil
}
