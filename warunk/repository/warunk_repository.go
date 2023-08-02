package repository

import (
	"context"
	"warunk-bem/domain"
	"warunk-bem/mongo"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type WarunkRepository struct {
	DB         mongo.Database
	Collection mongo.Collection
}

const (
	timeFormat     = "2006-01-02T15:04:05.999Z07:00" // reduce precision from RFC3339Nano as date format
	collectionName = "warunk"
)

func NewWarunkRepository(DB mongo.Database) domain.WarunkRepository {
	return &WarunkRepository{DB, DB.Collection(collectionName)}
}

func (kr *WarunkRepository) InsertOne(ctx context.Context, req *domain.Warunk) (*domain.Warunk, error) {
	var err error

	_, err = kr.Collection.InsertOne(ctx, req)
	if err != nil {
		return nil, err
	}

	return req, nil
}

func (kr *WarunkRepository) FindOneWarunk(ctx context.Context, id string) (*domain.Warunk, error) {
	var (
		Warunk domain.Warunk
		err    error
	)

	idHex, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return &Warunk, err
	}

	err = kr.Collection.FindOne(ctx, bson.M{"_id": idHex}).Decode(&Warunk)
	if err != nil {
		return &Warunk, err
	}

	return &Warunk, nil
}

func (kr *WarunkRepository) FindOne(ctx context.Context, id string) (*domain.Warunk, error) {
	var (
		Warunk domain.Warunk
		err    error
	)

	idHex, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return &Warunk, err
	}

	err = kr.Collection.FindOne(ctx, bson.M{"user_id": idHex}).Decode(&Warunk)
	if err != nil {
		return &Warunk, err
	}

	return &Warunk, nil
}

func (kr *WarunkRepository) GetAllWithPage(ctx context.Context, rp int64, p int64, filter interface{}, setsort interface{}) ([]domain.Warunk, int64, error) {
	var (
		Warunk []domain.Warunk
		total  int64
		err    error
	)

	total, err = kr.Collection.CountDocuments(ctx, filter)
	if err != nil {
		return nil, 0, err
	}

	opts := options.Find()
	opts.SetSort(setsort)
	opts.SetLimit(rp)
	opts.SetSkip((p - 1) * rp)

	cur, err := kr.Collection.Find(ctx, filter, opts)
	if err != nil {
		return nil, 0, err
	}

	for cur.Next(ctx) {
		var k domain.Warunk
		err := cur.Decode(&k)
		if err != nil {
			return nil, 0, err
		}

		Warunk = append(Warunk, k)
	}

	return Warunk, total, nil
}

func (kr *WarunkRepository) UpdateOne(ctx context.Context, Warunk *domain.Warunk, id string) (*domain.Warunk, error) {
	var err error

	idHex, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return Warunk, err
	}

	filter := bson.M{"user_id": idHex}
	update := bson.M{"$set": Warunk}

	_, err = kr.Collection.UpdateOne(ctx, filter, update)
	if err != nil {
		return Warunk, err
	}

	err = kr.Collection.FindOne(ctx, bson.M{"user_id": idHex}).Decode(Warunk)
	if err != nil {
		return Warunk, err
	}
	return Warunk, nil
}

func (kr *WarunkRepository) UpdateOneWarunk(ctx context.Context, Warunk *domain.Warunk, id string) (*domain.Warunk, error) {
	var err error

	idHex, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return Warunk, err
	}

	filter := bson.M{"_id": idHex}
	update := bson.M{"$set": Warunk}

	_, err = kr.Collection.UpdateOne(ctx, filter, update)
	if err != nil {
		return Warunk, err
	}

	err = kr.Collection.FindOne(ctx, bson.M{"_id": idHex}).Decode(Warunk)
	if err != nil {
		return Warunk, err
	}
	return Warunk, nil
}

func (kr *WarunkRepository) RemoveProduct(ctx context.Context, WarunkID string, productID string) error {
	var (
		Warunk *domain.Warunk
		err    error
	)

	WarunkObjectID, err := primitive.ObjectIDFromHex(WarunkID)
	if err != nil {
		return err
	}

	productObjectID, err := primitive.ObjectIDFromHex(productID)
	if err != nil {
		return err
	}

	// Cari Keranjang berdasarkan ID Keranjang
	err = kr.Collection.FindOne(ctx, bson.M{"_id": WarunkObjectID}).Decode(&Warunk)
	if err != nil {
		return err
	}

	// Mencari data Produk di dalam Array Produk Keranjang
	var updateProduk []domain.Produk
	for _, p := range Warunk.Produk {
		if p.ID != productObjectID {
			updateProduk = append(updateProduk, p)
		}
	}

	// Update Produk Array untuk menghapus data produk di keranjang
	update := bson.M{"$set": bson.M{"produk": updateProduk}}
	_, err = kr.Collection.UpdateOne(ctx, bson.M{"_id": WarunkObjectID}, update)
	if err != nil {
		return err
	}

	return nil
}

func (kr *WarunkRepository) DeleteOne(ctx context.Context, id string) error {
	var err error

	idHex, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}

	_, err = kr.Collection.DeleteOne(ctx, bson.M{"_id": idHex})
	if err != nil {
		return err
	}

	return nil
}
