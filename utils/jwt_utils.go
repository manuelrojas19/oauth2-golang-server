package utils

import (
	"crypto/rand"
	"encoding/hex"
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"github.com/manuelrojas19/go-oauth2-server/configuration"
)

// GenerateJWT creates a JWT token with the given parameters.
func GenerateJWT(clientId *string, userId *string, secretKey interface{}, tokenType string) (string, error) {
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
			"iat":  time.Now().Unix(),
			"exp":  time.Now().Add(expirationTime).Unix(),
			"type": tokenType,
			"jti":  generateRandomString(),
		})

		if clientId != nil {
			token.Claims.(jwt.MapClaims)["clientId"] = *clientId
		}

		if userId != nil {
			token.Claims.(jwt.MapClaims)["userId"] = *userId
		}

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
			return "", fmt.Errorf("invalid key type for HS256; expected []byte, got %T", secretKey)
		}
		token := jwt.NewWithClaims(signingMethod, jwt.MapClaims{
			"iat":  time.Now().Unix(),
			"exp":  time.Now().Add(expirationTime).Unix(),
			"type": tokenType,
			"jti":  generateRandomString(),
		})

		if clientId != nil {
			token.Claims.(jwt.MapClaims)["clientId"] = *clientId
		}
		if userId != nil {
			token.Claims.(jwt.MapClaims)["userId"] = *userId
		}

		tokenString, err := token.SignedString(key)
		if err != nil {
			return "", fmt.Errorf("failed to sign token: %w", err)
		}
		return tokenString, nil

	default:
		return "", fmt.Errorf("unsupported token type: %s", tokenType)
	}
}

// ValidateRefreshToken validates the JWT token using the provided secret key and returns the claims if valid.
func ValidateRefreshToken(tokenString string, secretKey []byte) (jwt.MapClaims, error) {
	// Parse and validate the token
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// Check token signing method
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unsupported JWT signing method for refresh token: %s", token.Method.Alg())
		}
		return secretKey, nil
	})
	if err != nil {
		return nil, fmt.Errorf("failed to parse or validate refresh token: %w", err)
	}

	// Check if the token is valid
	if !token.Valid {
		return nil, errors.New("refresh token is invalid")
	}

	// Extract claims
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil, errors.New("refresh token contains invalid claims")
	}

	// Validate expiration time
	exp, ok := claims["exp"].(float64)
	if !ok || time.Unix(int64(exp), 0).Before(time.Now()) {
		return nil, errors.New("refresh token has expired")
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
