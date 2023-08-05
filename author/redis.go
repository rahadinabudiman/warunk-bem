package author

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/go-redis/redis/v8"
	"github.com/joho/godotenv"
)

func InitRedisClient() *redis.Client {
	var (
		ctx         context.Context
		redisclient *redis.Client
	)

	ctx = context.TODO()

	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	if os.Getenv("REDIS_STATUS") == os.Getenv("REDIS_KEY") {
		redisURL := os.Getenv("REDIS_URL_CLOUD")
		redisOptions, err := redis.ParseURL(redisURL)
		if err != nil {
			panic(err)
		}

		redisclient = redis.NewClient(redisOptions)

		err = redisclient.Ping(ctx).Err()
		if err != nil {
			panic(err)
		}

		fmt.Println("Redis client connected successfully...")

		return redisclient
	} else {
		redisclient = redis.NewClient(&redis.Options{
			Addr: os.Getenv("REDIS_URL"),
		})

		if _, err := redisclient.Ping(ctx).Result(); err != nil {
			panic(err)
		}

		err = redisclient.Set(ctx, "test", "Welcome to Golang with Redis and MongoDB",
			0).Err()
		if err != nil {
			panic(err)
		}

		fmt.Println("Redis client connected successfully...")

		return redisclient
	}
}
