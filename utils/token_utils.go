package utils

import (
	"bytes"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"errors"
	"fmt"
	"github.com/golang-jwt/jwt/v4"
	"github.com/manuelrojas19/go-oauth2-server/configuration/keymanager"
	"gopkg.in/square/go-jose.v2"
	"strconv"
	"time"

	"github.com/google/uuid"
)

// GenerateToken creates a unique, secure token.
func GenerateToken(clientId string, userId string, createdAt time.Time) (string, error) {
	// Create a buffer and concatenate the values
	buf := bytes.NewBufferString(clientId)
	buf.WriteString(userId)
	buf.WriteString(strconv.FormatInt(createdAt.UnixNano(), 10))

	// Generate a random UUID for added security
	randomUUID, err := uuid.NewRandom()
	if err != nil {
		return "", err
	}

	// Use SHA-256 to generate a secure hash
	tokenHash := sha256.Sum256([]byte(randomUUID.String() + buf.String()))

	// Encode the hash in Base64
	token := base64.URLEncoding.EncodeToString(tokenHash[:])
	return token, nil
}

// GenerateJWT creates a JWT token with the given parameters.
func GenerateJWT(clientId string, userId string, secretKey []byte, tokenType string) (string, error) {
	var expirationTime time.Duration
	switch tokenType {
	case "access":
		expirationTime = time.Hour
	case "refresh":
		expirationTime = 24 * 30 * time.Hour // 30 days
	default:
		return "", errors.New("invalid token type")
	}

	// Generate a random unique identifier for the token
	tokenID, err := generateRandomString(32) // 32 bytes = 64 characters
	if err != nil {
		return "", err
	}

	// Define JWT claims
	claims := jwt.MapClaims{
		"clientId": clientId,
		"userId":   userId,
		"iat":      time.Now().Unix(),
		"exp":      time.Now().Add(expirationTime).Unix(),
		"type":     tokenType,
		"jti":      tokenID,
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

// VerifyToken parses and verifies a JWT token.
func VerifyToken(tokenString string, secretKey []byte) (*jwt.Token, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// Ensure the token's signing method matches the expected method
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return secretKey, nil
	})
	if err != nil {
		return nil, err
	}
	return token, nil
}

// GenerateJWE creates a JWE token from the JWT payload using global keys
func GenerateJWE(jwtToken string) (string, error) {
	publicKey, err := keymanager.GetPublicKey()
	if err != nil {
		return "", err
	}
	cipher, err := jose.NewEncrypter(jose.A256GCM, jose.Recipient{Algorithm: jose.RSA_OAEP, Key: publicKey}, nil)

	if err != nil {
		return "", err
	}

	object, err := cipher.Encrypt([]byte(jwtToken))
	if err != nil {
		return "", err
	}

	serialized, err := object.CompactSerialize()
	if err != nil {
		return "", err
	}

	return serialized, nil
}

// DecryptJWE decrypts a JWE token using the provided RSA private key.
func DecryptJWE(privateKey *rsa.PrivateKey, jweToken string) (string, error) {
	// Parse the JWE token
	object, err := jose.ParseEncrypted(jweToken)
	if err != nil {
		return "", fmt.Errorf("error parsing JWE token: %w", err)
	}

	// Decrypt the token using the RSA private key
	decrypted, err := object.Decrypt(privateKey)
	if err != nil {
		return "", fmt.Errorf("error decrypting JWE token: %w", err)
	}

	return string(decrypted), nil
}

// GenerateRandomString generates a random string of the specified length.
func generateRandomString(length int) (string, error) {
	if length <= 0 {
		return "", errors.New("length must be positive")
	}

	// Generate random bytes
	b := make([]byte, length)
	_, err := rand.Read(b)
	if err != nil {
		return "", err
	}

	// Encode bytes to hex
	return hex.EncodeToString(b), nil
}
