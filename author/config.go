package author

import (
	"log"

	"github.com/spf13/viper"
)

func InitConfig() *viper.Viper {
	config := viper.New()

	config.SetConfigFile(`.env`)
	err := config.ReadInConfig()
	if err != nil {
		panic("Cannot find .env file")
	}

	if config.GetBool(`debug`) {
		log.Println("Service RUN on DEBUG mode")
	}

	return config
}
