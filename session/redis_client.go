package session

import (
	"github.com/manuelrojas19/go-oauth2-server/configuration"
	"github.com/redis/go-redis/v9"
)

func NewRedisClient() *redis.Client {
	return redis.NewClient(&redis.Options{
		Addr:     configuration.RedisAddr,
		Password: configuration.RedisPassword,
		DB:       configuration.RedisDB,
	})
}
