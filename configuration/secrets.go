package configuration

import (
	"fmt"
	"github.com/joho/godotenv"
	"os"
)

var (
	RedisAddr          string
	RedisPassword      string
	RedisDB            int
	DatabaseUrl        string
	GoogleClientID     string
	GoogleClientSecret string
	GoogleRedirectURL  string
	GoogleAuthURL      string
	GoogleTokenURL     string
	GoogleUserInfoURL  string
	Scopes             string
)

func LoadSecrets() error {
	if err := godotenv.Load(); err != nil {
		fmt.Println("Error loading .env file")
		return err
	}
	loadGoogleSecrets()
	loadDbSecrets()
	loadRedisSecrets()
	return nil
}

func loadRedisSecrets() {
	RedisAddr = os.Getenv("REDIS_URL")
	RedisPassword = os.Getenv("REDIS_PASSWORD")
	RedisPassword = os.Getenv("REDIS_DB")
}

func loadDbSecrets() {
	DatabaseUrl = os.Getenv("DATABASE_URL")
}

func loadGoogleSecrets() {
	// Load configuration from environment variables
	GoogleClientID = os.Getenv("GOOGLE_CLIENT_ID")
	GoogleClientSecret = os.Getenv("GOOGLE_CLIENT_SECRET")
	GoogleRedirectURL = os.Getenv("GOOGLE_REDIRECT_URL")
	GoogleAuthURL = os.Getenv("GOOGLE_AUTH_URL")
	GoogleTokenURL = os.Getenv("GOOGLE_TOKEN_URL")
	GoogleUserInfoURL = os.Getenv("GOOGLE_USER_INFO_URL")
	Scopes = os.Getenv("SCOPES")

	// Optionally, you can set default values if not defined
	if GoogleClientID == "" {
		GoogleClientID = "default-client-id"
	}
	if GoogleClientSecret == "" {
		GoogleClientSecret = "default-client-secret"
	}
}
