package session

import (
	"github.com/redis/go-redis/v9"
	"os"
)

var (
	RedisAddr     = os.Getenv("localhost:6379")
	RedisPassword = os.Getenv("my-password")
	RedisDB       = 0
)

func NewRedisClient() *redis.Client {
	return redis.NewClient(&redis.Options{
		Addr:     RedisAddr,
		Password: RedisPassword,
		DB:       RedisDB,
	})
}
