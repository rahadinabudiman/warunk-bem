package dtos

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type InsertTransaksiRequest struct {
	ID        primitive.ObjectID `bson:"_id" json:"id"`
	CreatedAt time.Time          `json:"created_at"`
	UpdatedAt time.Time          `json:"updated_at"`
	UserID    primitive.ObjectID `bson:"user_id" json:"user_id"`
	ProdukID  primitive.ObjectID `bson:"produk_id" json:"produk_id"`
	Total     int                `bson:"total" json:"total"`
}

type InsertTransaksiKeranjangRequest struct {
	ID     primitive.ObjectID `bson:"_id" json:"id"`
	UserID primitive.ObjectID `bson:"user_id" json:"user_id"`
}
