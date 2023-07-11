package domain

import (
	"context"
	"time"
	"warunk-bem/dtos"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Keranjang struct {
	ID        primitive.ObjectID `bson:"_id" json:"id"`
	CreatedAt time.Time          `bson:"created_at" json:"created_at"`
	UpdatedAt time.Time          `bson:"updated_at" json:"updated_at"`
	UserID    primitive.ObjectID `bson:"user_id" json:"user_id"`
	Produk    []Produk           `bson:"produk" json:"produk"`
	Total     int                `bson:"total" json:"total"`
}

type KeranjangRepository interface {
	InsertOne(ctx context.Context, req *Keranjang) (*Keranjang, error)
	FindOne(ctx context.Context, id string) (*Keranjang, error)
	GetAllWithPage(ctx context.Context, rp int64, p int64, filter interface{}, setsort interface{}) ([]Keranjang, int64, error)
	UpdateOne(ctx context.Context, keranjang *Keranjang, id string) (*Keranjang, error)
	DeleteOne(ctx context.Context, id string) error
}

type KeranjangUsecase interface {
	InsertOne(ctx context.Context, req *dtos.InsertKeranjangRequest) (*dtos.InsertKeranjangResponse, error)
	// FindOne(ctx context.Context, id string) (*dtos.InsertKeranjangResponse, error)
	// GetAllWithPage(ctx context.Context, rp int64, p int64, filter interface{}, setsort interface{}) ([]dtos.InsertKeranjangResponse, int64, error)
	// UpdateOne(Ctx context.Context, id string, req *Keranjang) (*dtos.InsertKeranjangResponse, error)
	// DeleteOne(ctx context.Context, id string) error
}
