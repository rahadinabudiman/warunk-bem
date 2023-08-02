package author

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

func InitConfig() {
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("err loading: %v", err)
	}

	// Load individual environment variables using os.Getenv
	debugMode := os.Getenv("debug")

	if debugMode != "" && debugMode == "true" {
		log.Println("Service RUN on DEBUG mode")
	}
}
