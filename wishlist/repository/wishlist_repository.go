package repository

import (
	"context"
	"warunk-bem/domain"
	"warunk-bem/mongo"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type WishlistRepository struct {
	DB         mongo.Database
	Collection mongo.Collection
}

const (
	timeFormat     = "2006-01-02T15:04:05.999Z07:00" // reduce precision from RFC3339Nano as date format
	collectionName = "wishlist"
)

func NewWishlistRepository(DB mongo.Database) domain.WishlistRepository {
	return &WishlistRepository{DB, DB.Collection(collectionName)}
}

func (kr *WishlistRepository) InsertOne(ctx context.Context, req *domain.WishlistProduk) (*domain.WishlistProduk, error) {
	var err error

	_, err = kr.Collection.InsertOne(ctx, req)
	if err != nil {
		return nil, err
	}

	return req, nil
}

func (kr *WishlistRepository) FindOneWishlist(ctx context.Context, id string) (*domain.WishlistProduk, error) {
	var (
		Wishlist domain.WishlistProduk
		err      error
	)

	idHex, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return &Wishlist, err
	}

	err = kr.Collection.FindOne(ctx, bson.M{"_id": idHex}).Decode(&Wishlist)
	if err != nil {
		return &Wishlist, err
	}

	return &Wishlist, nil
}

func (kr *WishlistRepository) FindOne(ctx context.Context, id string) (*domain.WishlistProduk, error) {
	var (
		Wishlist domain.WishlistProduk
		err      error
	)

	idHex, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return &Wishlist, err
	}

	err = kr.Collection.FindOne(ctx, bson.M{"user_id": idHex}).Decode(&Wishlist)
	if err != nil {
		return &Wishlist, err
	}

	return &Wishlist, nil
}

func (kr *WishlistRepository) GetAllWithPage(ctx context.Context, rp int64, p int64, filter interface{}, setsort interface{}) ([]domain.WishlistProduk, int64, error) {
	var (
		Wishlist []domain.WishlistProduk
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
		var k domain.WishlistProduk
		err := cur.Decode(&k)
		if err != nil {
			return nil, 0, err
		}

		Wishlist = append(Wishlist, k)
	}

	return Wishlist, total, nil
}

func (kr *WishlistRepository) UpdateOne(ctx context.Context, Wishlist *domain.WishlistProduk, id string) (*domain.WishlistProduk, error) {
	var err error

	idHex, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return Wishlist, err
	}

	filter := bson.M{"user_id": idHex}
	update := bson.M{"$set": Wishlist}

	_, err = kr.Collection.UpdateOne(ctx, filter, update)
	if err != nil {
		return Wishlist, err
	}

	err = kr.Collection.FindOne(ctx, bson.M{"user_id": idHex}).Decode(Wishlist)
	if err != nil {
		return Wishlist, err
	}
	return Wishlist, nil
}

func (kr *WishlistRepository) UpdateOneWishlist(ctx context.Context, Wishlist *domain.WishlistProduk, id string) (*domain.WishlistProduk, error) {
	var err error

	idHex, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return Wishlist, err
	}

	filter := bson.M{"_id": idHex}
	update := bson.M{"$set": Wishlist}

	_, err = kr.Collection.UpdateOne(ctx, filter, update)
	if err != nil {
		return Wishlist, err
	}

	err = kr.Collection.FindOne(ctx, bson.M{"_id": idHex}).Decode(Wishlist)
	if err != nil {
		return Wishlist, err
	}
	return Wishlist, nil
}

func (kr *WishlistRepository) RemoveProduct(ctx context.Context, WishlistID string, productID string) error {
	var (
		Wishlist *domain.WishlistProduk
		err      error
	)

	WishlistObjectID, err := primitive.ObjectIDFromHex(WishlistID)
	if err != nil {
		return err
	}

	productObjectID, err := primitive.ObjectIDFromHex(productID)
	if err != nil {
		return err
	}

	// Cari Keranjang berdasarkan ID Keranjang
	err = kr.Collection.FindOne(ctx, bson.M{"_id": WishlistObjectID}).Decode(&Wishlist)
	if err != nil {
		return err
	}

	// Mencari data Produk di dalam Array Produk Keranjang
	var updateProduk []domain.Produk
	for _, p := range Wishlist.Produk {
		if p.ID != productObjectID {
			updateProduk = append(updateProduk, p)
		}
	}

	// Update Produk Array untuk menghapus data produk di keranjang
	update := bson.M{"$set": bson.M{"produk": updateProduk}}
	_, err = kr.Collection.UpdateOne(ctx, bson.M{"_id": WishlistObjectID}, update)
	if err != nil {
		return err
	}

	return nil
}

func (kr *WishlistRepository) DeleteOne(ctx context.Context, id string) error {
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
