package repository

import (
	"context"
	"warunk-bem/domain"
	"warunk-bem/mongo"

	"go.mongodb.org/mongo-driver/bson"
)

type dashboardRepository struct {
	DB               mongo.Database
	ProdukCollection mongo.Collection
	UserCollection   mongo.Collection
}

const (
	produkCollectionName = "produk" // Ganti dengan nama koleksi produk yang sesuai
	userCollectionName   = "user"   // Ganti dengan nama koleksi user yang sesuai
)

func NewDashboardRepository(DB mongo.Database) domain.DashboardRepository {
	return &dashboardRepository{
		DB,
		DB.Collection(produkCollectionName),
		DB.Collection(userCollectionName),
	}
}

func (r *dashboardRepository) DashboardGetAll() ([]domain.Produk, *domain.User, error) {
	ctx := context.TODO()

	// Aggregation pipeline
	pipeline := []bson.M{
		{
			"$lookup": bson.M{
				"from":         userCollectionName,
				"localField":   "user_id",
				"foreignField": "_id",
				"as":           "user",
			},
		},
		{
			"$unwind": "$user",
		},
	}

	// Execute aggregation
	cursor, err := r.ProdukCollection.Aggregate(ctx, pipeline)
	if err != nil {
		return nil, nil, err
	}
	defer cursor.Close(ctx)

	// Process aggregation results
	var produk []domain.Produk
	var user *domain.User
	for cursor.Next(ctx) {
		var result struct {
			Produk domain.Produk `bson:"_id"`
			User   *domain.User  `bson:"user"`
		}
		if err := cursor.Decode(&result); err != nil {
			return nil, nil, err
		}

		produk = append(produk, result.Produk)
		user = result.User
	}

	return produk, user, nil
}
