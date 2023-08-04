package usecase

import (
	"context"
	"errors"
	"time"
	"warunk-bem/domain"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type WishlistUsecase struct {
	WishlistRepo   domain.WishlistRepository
	ProdukRepo     domain.ProdukRepository
	UserRepo       domain.UserRepository
	contextTimeout time.Duration
}

func NewWishlistUsecase(WishlistRepo domain.WishlistRepository, ProdukRepo domain.ProdukRepository, UserRepo domain.UserRepository, contextTimeout time.Duration) domain.WishlistUsecase {
	return &WishlistUsecase{
		WishlistRepo:   WishlistRepo,
		ProdukRepo:     ProdukRepo,
		UserRepo:       UserRepo,
		contextTimeout: contextTimeout,
	}
}

// AddWishlist godoc
// @Summary      Add Wishlist
// @Description  Add Wishlist
// @Tags         User - Wishlist
// @Accept       json
// @Produce      json
// @Param        request body dtos.InsertWishlistRequest true "Payload Body [RAW]"
// @Success      201 {object} dtos.WishlistCreatedResponse
// @Failure      400 {object} dtos.BadRequestResponse
// @Failure      401 {object} dtos.UnauthorizedResponse
// @Failure      403 {object} dtos.ForbiddenResponse
// @Failure      404 {object} dtos.NotFoundResponse
// @Failure      500 {object} dtos.InternalServerErrorResponse
// @Router       /wishlist [post]
// @Security BearerAuth
func (fu *WishlistUsecase) InsertOne(ctx context.Context, req *domain.InsertWishlistRequest) (*domain.InsertWishlistResponse, error) {
	var res *domain.InsertWishlistResponse

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

	_, err = fu.WishlistRepo.InsertOne(ctx, &domain.WishlistProduk{
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
		return nil, errors.New("cannot add produk to Wishlist")
	}

	return res, nil
}

func (fu *WishlistUsecase) FindOne(ctx context.Context, id string) (*domain.InsertWishlistResponse, error) {
	ctx, cancel := context.WithTimeout(ctx, fu.contextTimeout)
	defer cancel()

	Wishlist, err := fu.WishlistRepo.FindOne(ctx, id)
	if err != nil {
		return nil, errors.New("wishlist not found")
	}

	res := &domain.InsertWishlistResponse{
		ID:     Wishlist.ID.Hex(),
		UserID: Wishlist.UserID.Hex(),
		Produk: Wishlist.Produk,
	}

	return res, nil
}

func (fu *WishlistUsecase) UpdateOne(ctx context.Context, id string, req *domain.WishlistProduk) (*domain.WishlistProduk, error) {
	ctx, cancel := context.WithTimeout(ctx, fu.contextTimeout)
	defer cancel()

	result, err := fu.WishlistRepo.FindOne(ctx, id)
	if err != nil {
		return nil, errors.New("wishlist not found")
	}

	// Jika produk ada di Wishlist, maka tampilkan message
	for _, v := range result.Produk {
		if v.ID == req.Produk[0].ID {
			return nil, errors.New("produk already in Wishlist")
		}
	}

	result.Produk = append(result.Produk, req.Produk...)

	_, err = fu.WishlistRepo.UpdateOne(ctx, result, id)
	if err != nil {
		return nil, err
	}

	return result, nil
}

func (fu *WishlistUsecase) RemoveProduct(ctx context.Context, WishlistID string, productID string) (*domain.DeleteProductWishlistResponse, error) {
	ctx, cancel := context.WithTimeout(ctx, fu.contextTimeout)
	defer cancel()

	// Check Wishlist apakah sudah ada atau belum
	result, err := fu.WishlistRepo.FindOneWishlist(ctx, WishlistID)
	if err != nil {
		return nil, errors.New("wishlist not found")
	}

	// Check Produk apakah ada atau tidak
	produk, err := fu.ProdukRepo.FindOne(ctx, productID)
	if err != nil {
		return nil, errors.New("produk not found")
	}

	// Check apakah produk sudah ada di Wishlist atau belum
	var index = -1
	for i, v := range result.Produk {
		if v.ID.Hex() == produk.ID.Hex() {
			index = i
			break
		}
	}

	// Jika produk tidak ada di Wishlist
	if index == -1 {
		return nil, errors.New("product is not in Wishlist")
	}

	err = fu.WishlistRepo.RemoveProduct(ctx, WishlistID, productID)
	if err != nil {
		return nil, errors.New("failed to remove produk")
	}

	check, err := fu.WishlistRepo.FindOneWishlist(ctx, WishlistID)
	if err != nil {
		return nil, errors.New("wishlist id not found")
	}

	if len(check.Produk) == 0 {
		err = fu.WishlistRepo.DeleteOne(ctx, WishlistID)
		if err != nil {
			return nil, errors.New("cannot delete Wishlist")
		}
	}

	res := &domain.DeleteProductWishlistResponse{
		Name: produk.Name,
	}

	return res, nil
}
