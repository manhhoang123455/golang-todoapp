package redis

import (
	"context"
	"fmt"
	"github.com/joho/godotenv"
	"github.com/redis/go-redis/v9"
	"log"
	"os"
)

var ctx = context.Background()

func InitRedis() *redis.Client {
	// Load .env file
	errEnv := godotenv.Load()
	if errEnv != nil {
		panic("Failed to load env file")
	}

	dbHost := os.Getenv("REDIS_HOST")
	dbPassword := os.Getenv("REDIS_PASSWORD")

	log.Println("Testing Golang Redis")

	rdb := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:6379", dbHost),
		Password: dbPassword,
		DB:       0,
	})

	_, err := rdb.Ping(ctx).Result()
	if err != nil {
		fmt.Println("Failed to connect redis : " + err.Error())
		panic("Failed to connect redis : " + err.Error())
	}

	return rdb
}

func Close(rdb *redis.Client) {
	rdb.Close()
}
