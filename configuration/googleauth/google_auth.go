package googleauth

import (
	"fmt"
	"github.com/joho/godotenv"
	"os"
)

// Load the configuration from environment variables
var (
	GoogleClientID     string
	GoogleClientSecret string
	GoogleRedirectURL  string
	GoogleAuthURL      string
	GoogleTokenURL     string
	GoogleUserInfoURL  string
	Scopes             string
)

func LoadSecrets() {
	// Load .env file
	err := godotenv.Load()
	if err != nil {
		fmt.Println("Error loading .env file")
		return
	}

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
