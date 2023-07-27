package domain

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type FavoriteProduk struct {
	ID        primitive.ObjectID `bson:"_id" json:"id"`
	CreatedAt time.Time          `bson:"created_at" json:"created_at"`
	UpdatedAt time.Time          `bson:"updated_at" json:"updated_at"`
	UserID    primitive.ObjectID `bson:"user_id" json:"user_id"`
	Produk    []Produk           `bson:"produk" json:"produk"`
}

type InsertFavoriteRequest struct {
	ID        primitive.ObjectID `bson:"_id" json:"id"`
	CreatedAt time.Time          `json:"created_at"`
	UpdatedAt time.Time          `json:"updated_at"`
	UserID    string             `json:"user_id"`
	ProdukID  string             `json:"produk_id"`
	Produk    []Produk           `json:"produk"`
}

type DeleteProductFavoriteRequest struct {
	FavoriteID string `json:"favorite_id"`
	ProdukID   string `json:"produk_id"`
}

type DeleteProductFavoriteResponse struct {
	Name string `json:"name"`
}

type InsertFavoriteResponse struct {
	ID     string   `json:"id"`
	UserID string   `json:"user_id"`
	Produk []Produk `json:"produk"`
	Total  int      `json:"total"`
}

type FavoriteRepository interface {
	InsertOne(ctx context.Context, req *FavoriteProduk) (*FavoriteProduk, error)
	FindOne(ctx context.Context, id string) (*FavoriteProduk, error)
	FindOneFavorite(ctx context.Context, id string) (*FavoriteProduk, error)
	GetAllWithPage(ctx context.Context, rp int64, p int64, filter interface{}, setsort interface{}) ([]FavoriteProduk, int64, error)
	UpdateOne(ctx context.Context, favorite *FavoriteProduk, id string) (*FavoriteProduk, error)
	UpdateOneFavorite(ctx context.Context, favorite *FavoriteProduk, id string) (*FavoriteProduk, error)
	RemoveProduct(ctx context.Context, favoriteID string, productID string) error
	DeleteOne(ctx context.Context, id string) error
}

type FavoriteUsecase interface {
	InsertOne(ctx context.Context, req *InsertFavoriteRequest) (*InsertFavoriteResponse, error)
	FindOne(ctx context.Context, id string) (*InsertFavoriteResponse, error)
	UpdateOne(ctx context.Context, id string, req *FavoriteProduk) (*FavoriteProduk, error)
	RemoveProduct(ctx context.Context, favoriteID string, productID string) (*DeleteProductFavoriteResponse, error)
}
