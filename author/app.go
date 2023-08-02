package author

import (
	"warunk-bem/mongo"
)

var (
	App *Application
)

type Application struct {
	Mongo mongo.Client
}

func init() {
	AppInit()
}

func AppInit() {
	App = &Application{}
	App.Mongo = InitMongoDatabase()
}
