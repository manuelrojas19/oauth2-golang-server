package configuration

import (
	"os"
)

var (
	RedisAddr     string
	RedisPassword string
	RedisDB       int
)

func LoadRedisSecrets() {
	RedisAddr = os.Getenv("REDIS_URL")
	RedisPassword = os.Getenv("REDIS_PASSWORD")
	RedisPassword = os.Getenv("REDIS_DB")
}
