package author

import (
	"log"

	"github.com/spf13/viper"
)

func InitConfig() *viper.Viper {
	config := viper.New()

	config.SetConfigFile(`.env.dev`)
	err := config.ReadInConfig()
	if err != nil {
		panic("Cant Find File .env.dev")
	}

	if config.GetBool(`debug`) {
		log.Println("Service RUN on DEBUG mode")
	}

	return config
}
