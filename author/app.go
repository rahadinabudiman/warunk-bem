package author

import (
	"warunk-bem/mongo"

	"github.com/go-redis/redis/v8"
)

var App *Application

type Application struct {
	Mongo mongo.Client
	Redis redis.Client
}

func init() {
	AppInit()
}

func AppInit() {
	App = &Application{}
	App.Mongo = InitMongoDatabase()
	App.Redis = *InitRedisClient()
}
