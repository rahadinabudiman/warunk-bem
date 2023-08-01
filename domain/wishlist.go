package domain

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type WishlistProduk struct {
	ID        primitive.ObjectID `bson:"_id" json:"id"`
	CreatedAt time.Time          `bson:"created_at" json:"created_at"`
	UpdatedAt time.Time          `bson:"updated_at" json:"updated_at"`
	UserID    primitive.ObjectID `bson:"user_id" json:"user_id"`
	Produk    []Produk           `bson:"produk" json:"produk"`
}

type InsertWishlistRequest struct {
	ID        primitive.ObjectID `bson:"_id" json:"id"`
	CreatedAt time.Time          `json:"created_at"`
	UpdatedAt time.Time          `json:"updated_at"`
	UserID    string             `json:"user_id"`
	ProdukID  string             `json:"produk_id"`
	Produk    []Produk           `json:"produk"`
}

type DeleteProductWishlistRequest struct {
	WishlistID string `json:"Wishlist_id"`
	ProdukID   string `json:"produk_id"`
}

type DeleteProductWishlistResponse struct {
	Name string `json:"name"`
}

type InsertWishlistResponse struct {
	ID     string   `json:"id"`
	UserID string   `json:"user_id"`
	Produk []Produk `json:"produk"`
	Total  int      `json:"total"`
}

type WishlistRepository interface {
	InsertOne(ctx context.Context, req *WishlistProduk) (*WishlistProduk, error)
	FindOne(ctx context.Context, id string) (*WishlistProduk, error)
	FindOneWishlist(ctx context.Context, id string) (*WishlistProduk, error)
	GetAllWithPage(ctx context.Context, rp int64, p int64, filter interface{}, setsort interface{}) ([]WishlistProduk, int64, error)
	UpdateOne(ctx context.Context, Wishlist *WishlistProduk, id string) (*WishlistProduk, error)
	UpdateOneWishlist(ctx context.Context, Wishlist *WishlistProduk, id string) (*WishlistProduk, error)
	RemoveProduct(ctx context.Context, WishlistID string, productID string) error
	DeleteOne(ctx context.Context, id string) error
}

type WishlistUsecase interface {
	InsertOne(ctx context.Context, req *InsertWishlistRequest) (*InsertWishlistResponse, error)
	FindOne(ctx context.Context, id string) (*InsertWishlistResponse, error)
	UpdateOne(ctx context.Context, id string, req *WishlistProduk) (*WishlistProduk, error)
	RemoveProduct(ctx context.Context, WishlistID string, productID string) (*DeleteProductWishlistResponse, error)
}
