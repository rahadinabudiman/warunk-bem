package repository

import (
	"context"
	"fmt"
	"warunk-bem/domain"
	"warunk-bem/mongo"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type transaksiRepository struct {
	DB         mongo.Database
	Collection mongo.Collection
}

const (
	timeFormat     = "2006-01-02T15:04:05.999Z07:00" // reduce precision from RFC3339Nano as date format
	collectionName = "transaksi"
)

func NewTransaksiRepository(DB mongo.Database) domain.TransaksiRepository {
	return &transaksiRepository{DB, DB.Collection(collectionName)}
}

func (tr *transaksiRepository) InsertOne(ctx context.Context, req *domain.Transaksi) (*domain.Transaksi, error) {
	var err error

	_, err = tr.Collection.InsertOne(ctx, req)
	if err != nil {
		return nil, err
	}

	return req, nil
}

func (tr *transaksiRepository) FindOne(ctx context.Context, id string) (*domain.Transaksi, error) {
	var (
		transaksi domain.Transaksi
		err       error
	)

	idHex, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return &transaksi, err
	}

	err = tr.Collection.FindOne(ctx, bson.M{"_id": idHex}).Decode(&transaksi)
	if err != nil {
		return nil, err
	}

	return &transaksi, err
}

func (tr *transaksiRepository) FindAllByUserId(ctx context.Context, id string) ([]*domain.Transaksi, error) {
	var transaksis []*domain.Transaksi

	idHex, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return transaksis, err
	}

	filter := bson.M{"user_id": idHex}
	cursor, err := tr.Collection.Find(ctx, filter)
	if err != nil {
		return transaksis, err
	}
	defer cursor.Close(ctx)

	for cursor.Next(ctx) {
		var transaksi domain.Transaksi
		if err := cursor.Decode(&transaksi); err != nil {
			return transaksis, err
		}
		transaksis = append(transaksis, &transaksi)
	}

	if err != nil {
		return transaksis, err
	}

	return transaksis, nil
}

func (tr *transaksiRepository) GetAllWithPage(ctx context.Context, rp int64, p int64, filter interface{}, setsort interface{}) ([]domain.Transaksi, int64, error) {
	var (
		transaksi []domain.Transaksi
		skip      int64
		opts      *options.FindOptions
	)

	skip = (p * rp) - rp
	if setsort != nil {
		opts = options.MergeFindOptions(
			options.Find().SetLimit(rp),
			options.Find().SetSkip(skip),
			options.Find().SetSort(setsort),
		)
	} else {
		opts = options.MergeFindOptions(
			options.Find().SetLimit(rp),
			options.Find().SetSkip(skip),
		)
	}

	cursor, err := tr.Collection.Find(
		ctx,
		filter,
		opts,
	)

	if err != nil {
		return nil, 0, err
	}
	if cursor == nil {
		return nil, 0, fmt.Errorf("nil cursor value")
	}
	err = cursor.All(ctx, &transaksi)
	if err != nil {
		return nil, 0, err
	}

	count, err := tr.Collection.CountDocuments(ctx, filter)
	if err != nil {
		return transaksi, 0, err
	}

	return transaksi, count, err
}

func (tr *transaksiRepository) UpdateOne(ctx context.Context, transaksi *domain.Transaksi, id string) (*domain.Transaksi, error) {
	var err error

	idHex, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return transaksi, err
	}

	filter := bson.M{"_id": idHex}
	update := bson.M{"$set": transaksi}

	_, err = tr.Collection.UpdateOne(ctx, filter, update)
	if err != nil {
		return transaksi, err
	}

	err = tr.Collection.FindOne(ctx, bson.M{"_id": idHex}).Decode(transaksi)
	if err != nil {
		return transaksi, err
	}
	return transaksi, nil
}

func (tr *transaksiRepository) DeleteOne(ctx context.Context, id string) error {
	var err error

	idHex, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}

	_, err = tr.Collection.DeleteOne(ctx, bson.M{"_id": idHex})
	if err != nil {
		return err
	}

	return nil
}
