package utils

import (
	"bytes"
	"crypto/sha256"
	"encoding/base64"
	"errors"
	"github.com/golang-jwt/jwt/v4"
	"strconv"
	"strings"
	"time"

	"github.com/google/uuid"
)

func GenerateToken(clientId string, userId string, createdAt time.Time) (string, error) {
	// Create a buffer and concatenate the values
	buf := bytes.NewBufferString(clientId)
	buf.WriteString(userId)
	buf.WriteString(strconv.FormatInt(createdAt.UnixNano(), 10))

	// Generate a random UUID for better security
	randomUUID, err := uuid.NewRandom()
	if err != nil {
		return "", err
	}

	// Use SHA-256 to generate a secure hash
	tokenUUID := uuid.NewMD5(randomUUID, buf.Bytes())
	tokenHash := sha256.Sum256([]byte(tokenUUID.String()))

	// Encode the hash in base64 and remove padding characters
	token := base64.URLEncoding.EncodeToString(tokenHash[:])
	token = strings.ToUpper(strings.TrimRight(token, "="))

	return token, nil
}

// GenerateJWT generates a JWT token with the given client ID, user ID, secret key, token type, and expiration.
func GenerateJWT(clientId string, userId string, secretKey []byte, tokenType string) (string, error) {
	// Define the expiration time for different token types
	var expirationTime time.Duration
	if tokenType == "access" {
		expirationTime = time.Hour * 1 // Access token expires in 1 hour
	} else if tokenType == "refresh" {
		expirationTime = time.Hour * 24 * 30 // Refresh token expires in 30 days
	} else {
		return "", errors.New("invalid token type")
	}

	// Define JWT claims
	claims := jwt.MapClaims{
		"clientId": clientId,
		"userId":   userId,
		"iat":      time.Now().Unix(),
		"exp":      time.Now().Add(expirationTime).Unix(), // Token expires based on type
		"type":     tokenType,                             // Add token type to claims
	}

	// Create the token with HS256 algorithm
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Sign the token with the secret key
	tokenString, err := token.SignedString(secretKey)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}
