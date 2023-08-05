package usecase

import (
	"context"
	"encoding/json"
	"errors"
	"time"
	"warunk-bem/domain"
	"warunk-bem/dtos"
	"warunk-bem/helpers"

	"github.com/go-redis/redis/v8"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"golang.org/x/crypto/bcrypt"
)

type produkUsecase struct {
	ProdukRepo     domain.ProdukRepository
	UserRepo       domain.UserRepository
	RedisClient    *redis.Client
	contextTimeout time.Duration
}

func NewProdukUsecase(ProdukRepo domain.ProdukRepository, UserRepo domain.UserRepository, RedisClient *redis.Client, contextTimeout time.Duration) domain.ProdukUsecase {
	return &produkUsecase{
		ProdukRepo:     ProdukRepo,
		UserRepo:       UserRepo,
		RedisClient:    RedisClient,
		contextTimeout: contextTimeout,
	}
}

// AddProduk godoc
// @Summary      Add Produk
// @Description  Add Produk
// @Tags         Admin - Produk
// @Accept       json
// @Produce      json
// @Param        request body dtos.InsertProdukRequest true "Payload Body [RAW]"
// @Success      201 {object} dtos.ProdukCreatedResponse
// @Failure      400 {object} dtos.BadRequestResponse
// @Failure      401 {object} dtos.UnauthorizedResponse
// @Failure      403 {object} dtos.ForbiddenResponse
// @Failure      404 {object} dtos.NotFoundResponse
// @Failure      500 {object} dtos.InternalServerErrorResponse
// @Router       /produk [post]
// @Security BearerAuth
func (pu *produkUsecase) InsertOne(c context.Context, req *dtos.InsertProdukRequest, url string) (*dtos.InsertProdukResponse, error) {
	var res *dtos.InsertProdukResponse

	ctx, cancel := context.WithTimeout(c, pu.contextTimeout)
	defer cancel()

	req.ID = primitive.NewObjectID()
	req.CreatedAt = time.Now()
	req.UpdatedAt = time.Now()

	slug := helpers.CreateSlug(req.Name)
	imageUrl := url

	// Check Slug apakah sudah ada atau belum
	checkSlug, err := pu.ProdukRepo.FindSlug(ctx, slug)
	if err == nil {
		if checkSlug.Slug == slug {
			slug = slug + "-" + helpers.RandomString(3)
		}
	}

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

// GetProdukByID godoc
// @Summary      Get Produk by ID
// @Description  Get Produk by ID
// @Tags         Produk
// @Accept       json
// @Produce      json
// @Param id path string true "ID Produk"
// @Success      200 {object} dtos.ProdukOKResponse
// @Failure      400 {object} dtos.BadRequestResponse
// @Failure      401 {object} dtos.UnauthorizedResponse
// @Failure      403 {object} dtos.ForbiddenResponse
// @Failure      404 {object} dtos.NotFoundResponse
// @Failure      500 {object} dtos.InternalServerErrorResponse
// @Router       /produk/{id} [get]
func (pu *produkUsecase) FindOne(c context.Context, id string) (res *dtos.ProdukDetailResponse, err error) {
	// Check if the result exists in Redis cache
	cacheKey := "produk:" + id
	val, err := pu.RedisClient.Get(c, cacheKey).Result()
	if err == nil {
		// Cache hit, unmarshal the cached value and return it
		res = &dtos.ProdukDetailResponse{}
		if err := json.Unmarshal([]byte(val), res); err != nil {
			return nil, err
		}
		return res, nil
	}

	ctx, cancel := context.WithTimeout(c, pu.contextTimeout)
	defer cancel()

	req, err := pu.ProdukRepo.FindOne(ctx, id)
	if err != nil {
		return res, err
	}

	res = &dtos.ProdukDetailResponse{
		ID:       req.ID.Hex(),
		Name:     req.Name,
		Slug:     req.Slug,
		Detail:   req.Detail,
		Price:    req.Price,
		Stock:    req.Stock,
		Category: req.Category,
		Image:    req.Image,
	}

	cacheValue, err := json.Marshal(res)
	if err == nil {
		pu.RedisClient.Set(c, cacheKey, cacheValue, 10*time.Minute)
	}

	return res, nil
}

// GetAllProduk godoc
// @Summary      Get All Produk
// @Description  Get All Produk
// @Tags         Produk
// @Accept       json
// @Produce      json
// @Success      200 {object} dtos.ProdukDetailResponse
// @Failure      400 {object} dtos.BadRequestResponse
// @Failure      401 {object} dtos.UnauthorizedResponse
// @Failure      403 {object} dtos.ForbiddenResponse
// @Failure      404 {object} dtos.NotFoundResponse
// @Failure      500 {object} dtos.InternalServerErrorResponse
// @Router       /produk [get]
func (pu *produkUsecase) GetAllWithPage(c context.Context, rp int64, p int64, filter interface{}, setsort interface{}) ([]*dtos.ProdukDetailResponse, int64, error) {
	var (
		res []*dtos.ProdukDetailResponse
		err error
	)

	cacheKey := "produk"
	val, err := pu.RedisClient.Get(c, cacheKey).Result()
	if err == nil {
		if err := json.Unmarshal([]byte(val), &res); err != nil {
			return nil, 0, err
		}
		return res, int64(len(res)), nil
	}

	ctx, cancel := context.WithTimeout(c, pu.contextTimeout)
	defer cancel()

	req, count, err := pu.ProdukRepo.GetAllWithPage(ctx, rp, p, filter, setsort)
	if err != nil {
		return nil, count, err
	}

	for _, v := range req {
		res = append(res, &dtos.ProdukDetailResponse{
			ID:       v.ID.Hex(),
			Name:     v.Name,
			Slug:     v.Slug,
			Detail:   v.Detail,
			Price:    v.Price,
			Stock:    v.Stock,
			Category: v.Category,
			Image:    v.Image,
		})
	}

	cacheValue, err := json.Marshal(res)
	if err == nil {
		pu.RedisClient.Set(c, cacheKey, cacheValue, 10*time.Minute)
	}

	return res, count, nil
}

// ProdukUpdate godoc
// @Summary      Update Produk
// @Description  Update Produk
// @Tags         Admin - Produk
// @Accept       json
// @Produce      json
// @Param        request body dtos.ProdukUpdateRequest true "Payload Body [RAW]"
// @Success      200 {object} dtos.ProdukOKResponse
// @Failure      400 {object} dtos.BadRequestResponse
// @Failure      401 {object} dtos.UnauthorizedResponse
// @Failure      403 {object} dtos.ForbiddenResponse
// @Failure      404 {object} dtos.NotFoundResponse
// @Failure      500 {object} dtos.InternalServerErrorResponse
// @Router       /produk/{id} [put]
// @Security BearerAuth
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

	cacheKey := "produk:" + id
	cacheValue, err := json.Marshal(res)
	if err == nil {
		pu.RedisClient.Set(c, cacheKey, cacheValue, 10*time.Minute)
	}

	return res, nil
}

// DeleteProduk godoc
// @Summary      Delete a Produk
// @Description  Delete a Produk
// @Tags         Admin - Produk
// @Accept       json
// @Produce      json
// @Param id path integer true "ID Produk"
// @Success      200 {object} dtos.StatusOKDeletedResponse
// @Failure      400 {object} dtos.BadRequestResponse
// @Failure      401 {object} dtos.UnauthorizedResponse
// @Failure      403 {object} dtos.ForbiddenResponse
// @Failure      404 {object} dtos.NotFoundResponse
// @Failure      500 {object} dtos.InternalServerErrorResponse
// @Router       /produk/{id} [delete]
// @Security BearerAuth
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
