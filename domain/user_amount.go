package domain

import (
	"context"
	"time"
	"warunk-bem/dtos"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type UserAmount struct {
	ID        primitive.ObjectID `bson:"_id" json:"id"`
	CreatedAt time.Time          `bson:"created_at" json:"created_at"`
	UpdatedAt time.Time          `bson:"updated_at" json:"updated_at"`
	UserID    primitive.ObjectID `bson:"user_id" json:"user_id"`
	Amount    float64            `bson:"amount" json:"amount"`
}

type UserAmountRepository interface {
	InsertOne(ctx context.Context, req *UserAmount) (res *UserAmount, err error)
	FindOne(ctx context.Context, id string) (res *UserAmount, err error)
	UpdateOne(ctx context.Context, amount *UserAmount, id string) (res *UserAmount, err error)
}

type UserAmountUsecase interface {
	TopUpSaldo(ctx context.Context, req *dtos.TopUpSaldoRequest) (res *dtos.TopUpSaldoResponse, err error)
}
