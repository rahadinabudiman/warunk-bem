package author

import (
	"warunk-bem/mongo"

	"github.com/spf13/viper"
)

var (
	App *Application
)

type Application struct {
	Config *viper.Viper
	Mongo  mongo.Client
}

func init() {
	AppInit()
}

func AppInit() {
	App = &Application{}
	App.Config = InitConfig()
	App.Mongo = InitMongoDatabase()
}
