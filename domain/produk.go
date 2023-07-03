package domain

import (
	"context"
	"time"
	"warunk-bem/domain/dtos"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Produk struct {
	ID        primitive.ObjectID `bson:"_id" json:"id"`
	CreatedAt time.Time          `bson:"created_at" json:"created_at"`
	UpdatedAt time.Time          `bson:"updated_at" json:"updated_at"`
	Slug      string             `bson:"slug" json:"slug"`
	Name      string             `bson:"name" json:"name"`
	Detail    string             `bson:"detail" json:"detail"`
	Price     int64              `bson:"price" json:"price"`
	Stock     int64              `bson:"stock" json:"stock"`
	Category  string             `bson:"category" json:"category"`
	Image     string             `bson:"image" json:"image" form:"image"`
}

type ProdukRepository interface {
	InsertOne(ctx context.Context, req *Produk) (*Produk, error)
	FindOne(ctx context.Context, id string) (*Produk, error)
	GetAllWithPage(ctx context.Context, rp int64, p int64, filter interface{}, setsort interface{}) ([]Produk, int64, error)
	UpdateOne(ctx context.Context, produk *Produk, id string) (*Produk, error)
	DeleteOne(ctx context.Context, id string) error
}

type ProdukUsecase interface {
	InsertOne(ctx context.Context, req *dtos.InsertProdukRequest, url string) (*dtos.InsertProdukResponse, error)
	FindOne(ctx context.Context, id string) (res *dtos.ProdukDetailResponse, err error)
	GetAllWithPage(ctx context.Context, rp int64, p int64, filter interface{}, setsort interface{}) ([]dtos.ProdukDetailResponse, int64, error)
	UpdateOne(ctx context.Context, req *dtos.ProdukUpdateRequest, id string) (*dtos.ProdukDetailResponse, error)
	DeleteOne(ctx context.Context, id string, idAdmin string, req dtos.DeleteProdukRequest) (res dtos.ResponseMessage, err error)
}
