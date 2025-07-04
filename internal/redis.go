package internal

import (
	"github.com/go-redis/redis/v8"
	"os"
)

var RedisClient *redis.Client

func RedisConnection() {
	RedisClient = redis.NewClient(&redis.Options{
		Addr:       os.Getenv("REDIS_CONN"),
		MaxRetries: 5,
	})
}
