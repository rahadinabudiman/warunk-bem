package dtos

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type InsertProdukRequest struct {
	ID        primitive.ObjectID `bson:"_id" json:"id"`
	CreatedAt time.Time          `json:"created_at"`
	UpdatedAt time.Time          `json:"updated_at"`
	Name      string             `bson:"name" json:"name"`
	Detail    string             `bson:"detail" json:"detail"`
	Price     int64              `bson:"price" json:"price"`
	Stock     int64              `bson:"stock" json:"stock"`
	Category  string             `bson:"category" json:"category"`
}

type ProdukUpdateRequest struct {
	Name     string `bson:"name" json:"name"`
	Detail   string `bson:"detail" json:"detail"`
	Price    int64  `bson:"price" json:"price"`
	Stock    int64  `bson:"stock" json:"stock"`
	Category string `bson:"category" json:"category"`
}

type DeleteProdukRequest struct {
	Password string `json:"password" form:"password" validate:"required" example:"rahadinabudimansundara"`
}
