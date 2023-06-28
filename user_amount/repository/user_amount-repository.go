package repository

import (
	"context"
	"warunk-bem/domain"
	"warunk-bem/mongo"
)

type userAmountRepository struct {
	DB         mongo.Database
	Collection mongo.Collection
}

const (
	timeFormat     = "2006-01-02T15:04:05.999Z07:00" // reduce precision from RFC3339Nano as date format
	collectionName = "user_amount"
)

func NewUserAmountRepository(DB mongo.Database) domain.UserAmountRepository {
	return &userAmountRepository{DB, DB.Collection(collectionName)}
}

func (r *userAmountRepository) InsertOne(ctx context.Context, req *domain.UserAmount) (res *domain.UserAmount, err error) {
	_, err = r.Collection.InsertOne(ctx, req)
	if err != nil {
		return nil, err
	}

	return req, nil
}
