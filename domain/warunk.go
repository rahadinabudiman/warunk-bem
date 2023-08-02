package domain

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Warunk struct {
	ID        primitive.ObjectID `bson:"_id" json:"id"`
	CreatedAt time.Time          `bson:"created_at" json:"created_at"`
	UpdatedAt time.Time          `bson:"updated_at" json:"updated_at"`
	DeletedAt time.Time          `bson:"deleted_at" json:"deleted_at"`
	UserID    primitive.ObjectID `bson:"user_id" json:"user_id"`
	Produk    []Produk           `bson:"produk" json:"produk"`
	Status    string             `bson:"status" json:"status"`
}

type CatalogWarunk struct {
	ID    primitive.ObjectID `bson:"_id" json:"id"`
	Stock int64              `bson:"stock" json:"stock"`
}

type InsertWarunkRequest struct {
	ID        primitive.ObjectID `bson:"_id" json:"id"`
	CreatedAt time.Time          `json:"created_at"`
	UpdatedAt time.Time          `json:"updated_at"`
	DeletedAt time.Time          `json:"deleted_at"`
	UserID    string             `json:"user_id"`
	Produk    []CatalogWarunk    `json:"produk"`
	Status    string             `json:"status"`
}

type InsertWarunkResponse struct {
	ID     string   `json:"id"`
	UserID string   `json:"user_id"`
	Produk []Produk `json:"produk"`
	Status string   `json:"status"`
}

type WarunkRepository interface {
	InsertOne(ctx context.Context, req *Warunk) (*Warunk, error)
	FindOne(ctx context.Context, id string) (*Warunk, error)
	FindOneByStatusAndDate(ctx context.Context, status string, date string) (*Warunk, error)
	FindOneByStatus(ctx context.Context, status string) (*Warunk, error)
	FindOneWarunk(ctx context.Context, id string) (*Warunk, error)
	GetAllWithPage(ctx context.Context, rp int64, p int64, filter interface{}, setsort interface{}) ([]Warunk, int64, error)
	UpdateOne(ctx context.Context, Warunk *Warunk, id string) (*Warunk, error)
	UpdateOneWarunk(ctx context.Context, Warunk *Warunk, id string) (*Warunk, error)
	RemoveProduct(ctx context.Context, WarunkID string, productID string) error
	DeleteOne(ctx context.Context, id string) error
}

type WarunkUsecase interface {
	InsertOne(ctx context.Context, req *InsertWarunkRequest) (*InsertWarunkResponse, error)
	FindOne(ctx context.Context, id string) (*InsertWarunkResponse, error)
	UpdateOne(ctx context.Context, id string, req *Warunk) (*Warunk, error)
}
