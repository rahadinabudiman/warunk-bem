package helpers

import (
	"os"
	"strconv"
)

func GetEnvInt(key string) (int, error) {
	strVal := os.Getenv(key)
	if strVal == "" {
		return 0, nil // or return an error if you consider a missing value as an error
	}

	val, err := strconv.Atoi(strVal)
	if err != nil {
		return 0, err
	}

	return val, nil
}
