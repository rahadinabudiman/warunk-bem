package repository

import (
	"context"
	"warunk-bem/domain"
	"warunk-bem/mongo"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
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

func (r *userAmountRepository) FindOne(ctx context.Context, id string) (res *domain.UserAmount, err error) {
	var (
		amount domain.UserAmount
	)

	idHex, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return &amount, err
	}

	err = r.Collection.FindOne(ctx, bson.M{"user_id": idHex}).Decode(&amount)
	if err != nil {
		return &amount, err
	}

	return &amount, nil
}

func (r *userAmountRepository) UpdateOne(ctx context.Context, amount *domain.UserAmount, id string) (res *domain.UserAmount, err error) {
	idHex, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return amount, err
	}

	filter := bson.M{"_id": idHex}
	update := bson.M{"$set": bson.M{
		"amount": amount.Amount,
	}}

	_, err = r.Collection.UpdateOne(ctx, filter, update)
	if err != nil {
		return amount, err
	}

	err = r.Collection.FindOne(ctx, bson.M{"_id": idHex}).Decode(amount)
	if err != nil {
		return amount, err
	}
	return amount, nil
}
