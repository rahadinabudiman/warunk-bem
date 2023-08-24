package repository

import (
	"context"
	"time"
	"warunk-bem/domain"
	"warunk-bem/mongo"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type produkRepository struct {
	DB         mongo.Database
	Collection mongo.Collection
}

const (
	timeFormat     = "2006-01-02T15:04:05.999Z07:00" // reduce precision from RFC3339Nano as date format
	collectionName = "produk"
)

func NewProdukRepository(DB mongo.Database) domain.ProdukRepository {
	return &produkRepository{DB, DB.Collection(collectionName)}
}

func (r *produkRepository) InsertOne(ctx context.Context, req *domain.Produk) (*domain.Produk, error) {
	var (
		err error
	)

	_, err = r.Collection.InsertOne(ctx, req)
	if err != nil {
		return nil, err
	}

	return req, nil
}

func (r *produkRepository) FindSlug(ctx context.Context, slug string) (*domain.Produk, error) {
	var (
		produk domain.Produk
		err    error
	)

	err = r.Collection.FindOne(ctx, bson.M{"slug": slug}).Decode(&produk)
	if err != nil {
		return nil, err
	}

	return &produk, nil
}

func (r *produkRepository) FindOne(ctx context.Context, id string) (*domain.Produk, error) {
	var (
		produk domain.Produk
		err    error
	)

	idHex, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return &produk, err
	}

	err = r.Collection.FindOne(ctx, bson.M{"_id": idHex}).Decode(&produk)
	if err != nil {
		return &produk, err
	}

	return &produk, nil
}

func (r *produkRepository) GetAllWithPage(ctx context.Context, rp int64, p int64, filter interface{}, setsort interface{}) ([]domain.Produk, int64, error) {
	var (
		produk []domain.Produk
		err    error
	)

	// Add the condition for deleted_at being null or nil
	nullOrNilCondition := bson.M{"$in": []interface{}{nil, primitive.Null{}}}
	filterWithDeletedAt := bson.M{
		"$and": []interface{}{
			filter,
			bson.M{"deleted_at": nullOrNilCondition},
		},
	}

	findOptions := options.Find()
	findOptions.SetLimit(rp)
	findOptions.SetSkip((p - 1) * rp)
	findOptions.SetSort(setsort)

	cursor, err := r.Collection.Find(ctx, filterWithDeletedAt, findOptions)
	if err != nil {
		return produk, 0, err
	}

	for cursor.Next(ctx) {
		var produkTemp domain.Produk
		err = cursor.Decode(&produkTemp)
		if err != nil {
			return produk, 0, err
		}
		produk = append(produk, produkTemp)
	}

	total, err := r.Collection.CountDocuments(ctx, filterWithDeletedAt)
	if err != nil {
		return produk, 0, err
	}

	return produk, total, nil
}

func (r *produkRepository) UpdateOne(ctx context.Context, produk *domain.Produk, id string) (*domain.Produk, error) {
	var (
		err error
	)

	idHex, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return produk, err
	}

	produk.UpdatedAt = time.Now()

	_, err = r.Collection.UpdateOne(ctx, bson.M{"_id": idHex}, bson.M{"$set": produk})
	if err != nil {
		return produk, err
	}

	return produk, nil
}
func (r *produkRepository) DeleteOne(ctx context.Context, id string) error {
	idHex, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}

	filter := bson.M{"_id": idHex}
	update := bson.M{
		"$set": bson.M{
			"deleted_at": time.Now(),
		},
	}

	_, err = r.Collection.UpdateOne(ctx, filter, update)
	if err != nil {
		return err
	}

	return nil
}
