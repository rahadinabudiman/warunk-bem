package usecase

import (
	"context"
	"errors"
	"time"
	"warunk-bem/domain"
	"warunk-bem/domain/dtos"
	"warunk-bem/helpers"

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

func (pu *produkUsecase) InsertOne(c context.Context, req *dtos.InsertProdukRequest, url string) (*dtos.InsertProdukResponse, error) {
	var res *dtos.InsertProdukResponse

	ctx, cancel := context.WithTimeout(c, pu.contextTimeout)
	defer cancel()

	req.ID = primitive.NewObjectID()
	req.CreatedAt = time.Now()
	req.UpdatedAt = time.Now()

	slug := helpers.CreateSlug(req.Name)
	imageUrl := url

	CreateProduk := &domain.Produk{
		ID:        req.ID,
		CreatedAt: req.CreatedAt,
		UpdatedAt: req.UpdatedAt,
		Slug:      slug,
		Name:      req.Name,
		Detail:    req.Detail,
		Price:     req.Price,
		Stock:     req.Stock,
		Category:  req.Category,
		Image:     imageUrl,
	}

	createdProduk, err := pu.ProdukRepo.InsertOne(ctx, CreateProduk)
	if err != nil {
		return res, errors.New("failed to create Produk")
	}

	res = &dtos.InsertProdukResponse{
		Name:     createdProduk.Name,
		Slug:     createdProduk.Slug,
		Detail:   createdProduk.Detail,
		Price:    createdProduk.Price,
		Stock:    createdProduk.Stock,
		Category: createdProduk.Category,
		Image:    createdProduk.Image,
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
		Slug:     req.Slug,
		Detail:   req.Detail,
		Price:    req.Price,
		Stock:    req.Stock,
		Category: req.Category,
		Image:    req.Image,
	}

	return res, nil
}

func (pu *produkUsecase) GetAllWithPage(c context.Context, rp int64, p int64, filter interface{}, setsort interface{}) ([]dtos.ProdukDetailResponse, int64, error) {
	var res []dtos.ProdukDetailResponse

	ctx, cancel := context.WithTimeout(c, pu.contextTimeout)
	defer cancel()

	req, count, err := pu.ProdukRepo.GetAllWithPage(ctx, rp, p, filter, setsort)
	if err != nil {
		return res, count, err
	}

	for _, v := range req {
		res = append(res, dtos.ProdukDetailResponse{
			Name:     v.Name,
			Slug:     v.Slug,
			Detail:   v.Detail,
			Price:    v.Price,
			Stock:    v.Stock,
			Category: v.Category,
			Image:    v.Image,
		})
	}

	return res, count, nil

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
	slug := helpers.CreateSlug(result.Name)
	result.Slug = slug
	result.Detail = req.Detail
	result.Price = req.Price
	result.Stock = req.Stock
	result.Category = req.Category
	result.Image = req.Image

	resp, err := pu.ProdukRepo.UpdateOne(ctx, result, id)
	if err != nil {
		return res, err
	}

	res = &dtos.ProdukDetailResponse{
		Name:     resp.Name,
		Slug:     resp.Slug,
		Detail:   resp.Detail,
		Price:    resp.Price,
		Stock:    resp.Stock,
		Category: resp.Category,
		Image:    resp.Image,
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
