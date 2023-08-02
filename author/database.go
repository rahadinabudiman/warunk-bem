package author

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"
	"warunk-bem/mongo"

	"github.com/joho/godotenv"
)

func InitMongoDatabase() mongo.Client {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	dbHost := os.Getenv(`MONGODB_HOST`)
	dbPort := os.Getenv(`MONGODB_PORT`)
	dbUser := os.Getenv(`MONGODB_USER`)
	dbPass := os.Getenv(`MONGODB_PASS`)
	dbName := os.Getenv(`MONGODB_NAME`)
	MongoDBStatus := os.Getenv(`MONGODB_STATUS`)
	MongoDBKey := os.Getenv(`MONGODB_KEY`)
	mongodbURI := os.Getenv(`MONGODB_URI`)
	if MongoDBStatus == MongoDBKey {
		mongodbURI = os.Getenv("MONGODB_URI")
	} else {
		mongodbURI = fmt.Sprintf("mongodb://%s:%s@%s:%s/?authMechanism=SCRAM-SHA-1&authSource=%s", dbUser, dbPass, dbHost, dbPort, dbName)
		if dbUser == "" || dbPass == "" {
			mongodbURI = fmt.Sprintf("mongodb://%s:%s/%s/", dbHost, dbPort, dbName)
		}
	}

	client, err := mongo.NewClient(mongodbURI)
	if err != nil {
		log.Fatal(err)
	}

	err = client.Connect(ctx)
	if err != nil {
		log.Fatal(err)
	}

	err = client.Ping(ctx)
	if err != nil {
		log.Fatal(err)
	}

	return client
}
