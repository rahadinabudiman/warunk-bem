package author

import (
	"context"
	"fmt"
	"log"
	"time"
	"warunk-bem/mongo"
)

func InitMongoDatabase() mongo.Client {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	dbHost := App.Config.GetString(`MONGODB_HOST`)
	dbPort := App.Config.GetString(`MONGODB_PORT`)
	dbUser := App.Config.GetString(`MONGODB_USER`)
	dbPass := App.Config.GetString(`MONGODB_PASS`)
	dbName := App.Config.GetString(`MONGODB_NAME`)
	mongodbURI := fmt.Sprintf("mongodb://%s:%s@%s:%s/?authMechanism=SCRAM-SHA-1&authSource=%s", dbUser, dbPass, dbHost, dbPort, dbName)
	if dbUser == "" || dbPass == "" {
		mongodbURI = fmt.Sprintf("mongodb://%s:%s/%s/", dbHost, dbPort, dbName)
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
