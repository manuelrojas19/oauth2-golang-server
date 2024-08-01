package configuration

import (
	"fmt"
	"github.com/joho/godotenv"
)

func LoadSecrets() error {
	if err := godotenv.Load(); err != nil {
		fmt.Println("Error loading .env file")
		return err
	}
	LoadGoogleSecrets()
	LoadDbSecrets()
	LoadRedisSecrets()
	return nil
}
