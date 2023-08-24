package dtos

import (
	"mime/multipart"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type InsertProdukRequest struct {
	ID        primitive.ObjectID    `bson:"_id" json:"id"`
	CreatedAt time.Time             `json:"created_at"`
	UpdatedAt time.Time             `json:"updated_at"`
	DeletedAt primitive.ObjectID    `bson:"deleted_at" json:"deleted_at"`
	Slug      string                `bson:"slug" json:"slug" form:"slug"`
	Name      string                `bson:"name" json:"name" form:"name" validate:"required"`
	Detail    string                `bson:"detail" json:"detail" form:"detail" validate:"required"`
	Price     int64                 `bson:"price" json:"price" form:"price" validate:"required"`
	Stock     int64                 `bson:"stock" json:"stock" form:"stock" validate:"required"`
	Category  string                `bson:"category" json:"category" form:"category" validate:"required"`
	Image     *multipart.FileHeader `bson:"image" json:"image" form:"image" validate:"required"`
}

type ImageProdukRequest struct {
	Image string `json:"image" form:"image" validate:"required"`
}

type ProdukUpdateRequest struct {
	Name     string `bson:"name" json:"name"`
	Detail   string `bson:"detail" json:"detail"`
	Price    int64  `bson:"price" json:"price"`
	Stock    int64  `bson:"stock" json:"stock"`
	Category string `bson:"category" json:"category"`
	Image    string `bson:"image" json:"image"`
}

type DeleteProdukRequest struct {
	Password string `json:"password" form:"password" validate:"required" example:"rahadinabudimansundara"`
}
