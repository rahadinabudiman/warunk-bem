package author

import "github.com/spf13/viper"

type Config struct {
	RedisUri string `mapstructure:"REDIS_URL"`
	Port     string `mapstructure:"PORT"`
}

func LoadConfig(path string) (config Config, err error) {
	viper.AddConfigPath(path)
	viper.SetConfigType("env")

	viper.AutomaticEnv()

	err = viper.ReadInConfig()
	if err != nil {
		return
	}

	err = viper.Unmarshal(&config)
	return
}
