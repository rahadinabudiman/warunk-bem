package helpers

import (
	"log"

	"github.com/spf13/viper"
)

type Config struct {
	EmailFrom    string `mapstructure:"EMAIL_FROM"`
	SMTPHost     string `mapstructure:"SMTP_HOST"`
	SMTPPass     string `mapstructure:"SMTP_PASS"`
	SMTPPort     int    `mapstructure:"SMTP_PORT"`
	SMTPUser     string `mapstructure:"SMTP_USER"`
	FromName     string `mapstructure:"FROM_NAME"`
	ClientOrigin string `mapstructure:"CLIENT_ORIGIN"`
}

func LoadConfig(path string) (Config, error) {
	var config Config

	viper.AddConfigPath(path)
	viper.SetConfigType("env")
	viper.SetConfigName(".env")

	viper.AutomaticEnv()

	err := viper.ReadInConfig()
	if err != nil {
		log.Fatalf("Failed to read config file: %v", err)
		return config, err
	}

	err = viper.Unmarshal(&config)
	if err != nil {
		log.Fatalf("Failed to unmarshal config file: %v", err)
		return config, err
	}

	return config, nil
}
