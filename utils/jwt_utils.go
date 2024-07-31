package utils

import (
	"crypto/rand"
	"encoding/hex"
	"errors"
	"fmt"
	"github.com/golang-jwt/jwt/v4"
	"github.com/manuelrojas19/go-oauth2-server/configuration"
	"time"
)

// GenerateJWT creates a JWT token with the given parameters.
func GenerateJWT(clientId string, userId string, secretKey interface{}, tokenType string) (string, error) {
	var expirationTime time.Duration
	var signingMethod jwt.SigningMethod

	switch tokenType {
	case "access":
		expirationTime = time.Hour
		signingMethod = jwt.SigningMethodRS256
		privateKey, err := configuration.GetJWTPrivateKey()
		if err != nil {
			return "", fmt.Errorf("private key is not initialized: %w", err)
		}
		token := jwt.NewWithClaims(signingMethod, jwt.MapClaims{
			"clientId": clientId,
			"userId":   userId,
			"iat":      time.Now().Unix(),
			"exp":      time.Now().Add(expirationTime).Unix(),
			"type":     tokenType,
			"jti":      generateRandomString(),
		})
		tokenString, err := token.SignedString(privateKey)
		if err != nil {
			return "", fmt.Errorf("failed to sign token: %w", err)
		}
		return tokenString, nil

	case "refresh":
		expirationTime = 24 * 30 * time.Hour
		signingMethod = jwt.SigningMethodHS256
		key, ok := secretKey.([]byte)
		if !ok {
			return "", errors.New("invalid key type for HS256")
		}
		token := jwt.NewWithClaims(signingMethod, jwt.MapClaims{
			"iat":  time.Now().Unix(),
			"exp":  time.Now().Add(expirationTime).Unix(),
			"type": tokenType,
			"jti":  generateRandomString(),
		})
		tokenString, err := token.SignedString(key)
		if err != nil {
			return "", fmt.Errorf("failed to sign token: %w", err)
		}
		return tokenString, nil

	default:
		return "", errors.New("invalid token type")
	}
}

// ValidateRefreshToken validates the JWT token using the provided secret key and returns the claims if valid.
func ValidateRefreshToken(tokenString string, secretKey []byte) (jwt.MapClaims, error) {
	// Parse and validate the token
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// Check token signing method
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return secretKey, nil
	})
	if err != nil {
		return nil, err
	}

	// Check if the token is valid
	if !token.Valid {
		return nil, errors.New("invalid token")
	}

	// Extract claims
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil, errors.New("invalid claims")
	}

	// Validate expiration time
	exp, ok := claims["exp"].(float64)
	if !ok || time.Unix(int64(exp), 0).Before(time.Now()) {
		return nil, errors.New("token has expired")
	}

	return claims, nil
}

// GenerateRandomString generates a random string of the specified length.
func generateRandomString() string {
	// Generate random bytes
	b := make([]byte, 32)
	_, err := rand.Read(b)
	if err != nil {
		return ""
	}

	// Encode bytes to hex
	return hex.EncodeToString(b)
}
