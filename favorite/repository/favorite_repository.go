package repository

import (
	"context"
	"warunk-bem/domain"
	"warunk-bem/mongo"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type FavoriteRepository struct {
	DB         mongo.Database
	Collection mongo.Collection
}

const (
	timeFormat     = "2006-01-02T15:04:05.999Z07:00" // reduce precision from RFC3339Nano as date format
	collectionName = "Favorite"
)

func NewFavoriteRepository(DB mongo.Database) domain.FavoriteRepository {
	return &FavoriteRepository{DB, DB.Collection(collectionName)}
}

func (kr *FavoriteRepository) InsertOne(ctx context.Context, req *domain.FavoriteProduk) (*domain.FavoriteProduk, error) {
	var err error

	_, err = kr.Collection.InsertOne(ctx, req)
	if err != nil {
		return nil, err
	}

	return req, nil
}

func (kr *FavoriteRepository) FindOneFavorite(ctx context.Context, id string) (*domain.FavoriteProduk, error) {
	var (
		favorite domain.FavoriteProduk
		err      error
	)

	idHex, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return &favorite, err
	}

	err = kr.Collection.FindOne(ctx, bson.M{"_id": idHex}).Decode(&favorite)
	if err != nil {
		return &favorite, err
	}

	return &favorite, nil
}

func (kr *FavoriteRepository) FindOne(ctx context.Context, id string) (*domain.FavoriteProduk, error) {
	var (
		favorite domain.FavoriteProduk
		err      error
	)

	idHex, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return &favorite, err
	}

	err = kr.Collection.FindOne(ctx, bson.M{"user_id": idHex}).Decode(&favorite)
	if err != nil {
		return &favorite, err
	}

	return &favorite, nil
}

func (kr *FavoriteRepository) GetAllWithPage(ctx context.Context, rp int64, p int64, filter interface{}, setsort interface{}) ([]domain.FavoriteProduk, int64, error) {
	var (
		favorite []domain.FavoriteProduk
		total    int64
		err      error
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
		var k domain.FavoriteProduk
		err := cur.Decode(&k)
		if err != nil {
			return nil, 0, err
		}

		favorite = append(favorite, k)
	}

	return favorite, total, nil
}

func (kr *FavoriteRepository) UpdateOne(ctx context.Context, favorite *domain.FavoriteProduk, id string) (*domain.FavoriteProduk, error) {
	var err error

	idHex, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return favorite, err
	}

	filter := bson.M{"user_id": idHex}
	update := bson.M{"$set": favorite}

	_, err = kr.Collection.UpdateOne(ctx, filter, update)
	if err != nil {
		return favorite, err
	}

	err = kr.Collection.FindOne(ctx, bson.M{"user_id": idHex}).Decode(favorite)
	if err != nil {
		return favorite, err
	}
	return favorite, nil
}

func (kr *FavoriteRepository) UpdateOneFavorite(ctx context.Context, favorite *domain.FavoriteProduk, id string) (*domain.FavoriteProduk, error) {
	var err error

	idHex, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return favorite, err
	}

	filter := bson.M{"_id": idHex}
	update := bson.M{"$set": favorite}

	_, err = kr.Collection.UpdateOne(ctx, filter, update)
	if err != nil {
		return favorite, err
	}

	err = kr.Collection.FindOne(ctx, bson.M{"_id": idHex}).Decode(favorite)
	if err != nil {
		return favorite, err
	}
	return favorite, nil
}

func (kr *FavoriteRepository) RemoveProduct(ctx context.Context, favoriteID string, productID string) error {
	var (
		favorite *domain.FavoriteProduk
		err      error
	)

	favoriteObjectID, err := primitive.ObjectIDFromHex(favoriteID)
	if err != nil {
		return err
	}

	productObjectID, err := primitive.ObjectIDFromHex(productID)
	if err != nil {
		return err
	}

	// Cari Keranjang berdasarkan ID Keranjang
	err = kr.Collection.FindOne(ctx, bson.M{"_id": favoriteObjectID}).Decode(&favorite)
	if err != nil {
		return err
	}

	// Mencari data Produk di dalam Array Produk Keranjang
	var updateProduk []domain.Produk
	for _, p := range favorite.Produk {
		if p.ID != productObjectID {
			updateProduk = append(updateProduk, p)
		}
	}

	// Update Produk Array untuk menghapus data produk di keranjang
	update := bson.M{"$set": bson.M{"produk": updateProduk}}
	_, err = kr.Collection.UpdateOne(ctx, bson.M{"_id": favoriteObjectID}, update)
	if err != nil {
		return err
	}

	return nil
}

func (kr *FavoriteRepository) DeleteOne(ctx context.Context, id string) error {
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
