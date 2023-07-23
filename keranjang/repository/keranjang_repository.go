package repository

import (
	"context"
	"warunk-bem/domain"
	"warunk-bem/mongo"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type keranjangRepository struct {
	DB         mongo.Database
	Collection mongo.Collection
}

const (
	timeFormat     = "2006-01-02T15:04:05.999Z07:00" // reduce precision from RFC3339Nano as date format
	collectionName = "keranjang"
)

func NewKeranjangRepository(DB mongo.Database) domain.KeranjangRepository {
	return &keranjangRepository{DB, DB.Collection(collectionName)}
}

func (kr *keranjangRepository) InsertOne(ctx context.Context, req *domain.Keranjang) (*domain.Keranjang, error) {
	var err error

	_, err = kr.Collection.InsertOne(ctx, req)
	if err != nil {
		return nil, err
	}

	return req, nil
}

func (kr *keranjangRepository) FindOneKeranjang(ctx context.Context, id string) (*domain.Keranjang, error) {
	var (
		keranjang domain.Keranjang
		err       error
	)

	idHex, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return &keranjang, err
	}

	err = kr.Collection.FindOne(ctx, bson.M{"_id": idHex}).Decode(&keranjang)
	if err != nil {
		return &keranjang, err
	}

	return &keranjang, nil
}

func (kr *keranjangRepository) FindOne(ctx context.Context, id string) (*domain.Keranjang, error) {
	var (
		keranjang domain.Keranjang
		err       error
	)

	idHex, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return &keranjang, err
	}

	err = kr.Collection.FindOne(ctx, bson.M{"user_id": idHex}).Decode(&keranjang)
	if err != nil {
		return &keranjang, err
	}

	return &keranjang, nil
}

func (kr *keranjangRepository) GetAllWithPage(ctx context.Context, rp int64, p int64, filter interface{}, setsort interface{}) ([]domain.Keranjang, int64, error) {
	var (
		keranjang []domain.Keranjang
		total     int64
		err       error
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
		var k domain.Keranjang
		err := cur.Decode(&k)
		if err != nil {
			return nil, 0, err
		}

		keranjang = append(keranjang, k)
	}

	return keranjang, total, nil
}

func (kr *keranjangRepository) UpdateOne(ctx context.Context, keranjang *domain.Keranjang, id string) (*domain.Keranjang, error) {
	var err error

	idHex, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return keranjang, err
	}

	filter := bson.M{"user_id": idHex}
	update := bson.M{"$set": keranjang}

	_, err = kr.Collection.UpdateOne(ctx, filter, update)
	if err != nil {
		return keranjang, err
	}

	err = kr.Collection.FindOne(ctx, bson.M{"user_id": idHex}).Decode(keranjang)
	if err != nil {
		return keranjang, err
	}
	return keranjang, nil
}

func (kr *keranjangRepository) UpdateOneKeranjang(ctx context.Context, keranjang *domain.Keranjang, id string) (*domain.Keranjang, error) {
	var err error

	idHex, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return keranjang, err
	}

	filter := bson.M{"_id": idHex}
	update := bson.M{"$set": keranjang}

	_, err = kr.Collection.UpdateOne(ctx, filter, update)
	if err != nil {
		return keranjang, err
	}

	err = kr.Collection.FindOne(ctx, bson.M{"_id": idHex}).Decode(keranjang)
	if err != nil {
		return keranjang, err
	}
	return keranjang, nil
}

func (kr *keranjangRepository) RemoveProduct(ctx context.Context, keranjangID string, productID string) error {
	var (
		keranjang *domain.Keranjang
		err       error
	)

	keranjangObjectID, err := primitive.ObjectIDFromHex(keranjangID)
	if err != nil {
		return err
	}

	productObjectID, err := primitive.ObjectIDFromHex(productID)
	if err != nil {
		return err
	}

	// Cari Keranjang berdasarkan ID Keranjang
	err = kr.Collection.FindOne(ctx, bson.M{"_id": keranjangObjectID}).Decode(&keranjang)
	if err != nil {
		return err
	}

	// Mencari data Produk di dalam Array Produk Keranjang
	var updateProduk []domain.Produk
	for _, p := range keranjang.Produk {
		if p.ID != productObjectID {
			updateProduk = append(updateProduk, p)
		}
	}

	// Update Produk Array untuk menghapus data produk di keranjang
	update := bson.M{"$set": bson.M{"produk": updateProduk}}
	_, err = kr.Collection.UpdateOne(ctx, bson.M{"_id": keranjangObjectID}, update)
	if err != nil {
		return err
	}

	return nil
}

func (kr *keranjangRepository) DeleteOne(ctx context.Context, id string) error {
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
