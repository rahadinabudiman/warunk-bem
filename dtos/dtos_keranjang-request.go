package dtos

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type InsertKeranjangRequest struct {
	ID        primitive.ObjectID `bson:"_id" json:"id"`
	CreatedAt time.Time          `json:"created_at"`
	UpdatedAt time.Time          `json:"updated_at"`
	UserID    string             `json:"user_id"`
	ProdukID  string             `json:"produk_id"`
	Total     int                `json:"total"`
}
