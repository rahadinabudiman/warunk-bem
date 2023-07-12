package domain

import (
	"context"
	"time"
	"warunk-bem/dtos"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Transaksi struct {
	ID        primitive.ObjectID `bson:"_id" json:"id"`
	CreatedAt time.Time          `bson:"created_at" json:"created_at"`
	UpdatedAt time.Time          `bson:"updated_at" json:"updated_at"`
	UserID    primitive.ObjectID `bson:"user_id" json:"user_id"`
	ProdukID  primitive.ObjectID `bson:"produk_id" json:"produk_id"`
	Total     int64              `bson:"total" json:"total"`
	Status    string             `bson:"status" json:"status"`
}

type TransaksiRepository interface {
	InsertOne(ctx context.Context, req *Transaksi) (*Transaksi, error)
	FindOne(ctx context.Context, id string) (*Transaksi, error)
	GetAllWithPage(ctx context.Context, rp int64, p int64, filter interface{}, setsort interface{}) ([]Transaksi, int64, error)
	UpdateOne(ctx context.Context, transaksi *Transaksi, id string) (*Transaksi, error)
	DeleteOne(ctx context.Context, id string) error
}

type TransaksiUsecase interface {
	InsertOne(ctx context.Context, req *dtos.InsertTransaksiRequest) (*dtos.InsertTransaksiResponse, error)
	InsertByKeranjang(ctx context.Context, req *dtos.InsertTransaksiKeranjangRequest) (*dtos.InsertTransaksiResponse, error)
	// FindOne(ctx context.Context, id string) (res *dtos.UserProfileResponse, err error)
	// GetAllWithPage(ctx context.Context, rp int64, p int64, filter interface{}, setsort interface{}) ([]dtos.UserProfileResponse, int64, error)
	// UpdateOne(ctx context.Context, user *dtos.UpdateUserRequest, id string) (*dtos.UpdateUserResponse, error)
	// DeleteOne(c context.Context, id string, req dtos.DeleteUserRequest) (res dtos.ResponseMessage, err error)
}
