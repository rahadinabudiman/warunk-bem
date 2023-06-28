package domain

import (
	"context"
	"time"

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
}
