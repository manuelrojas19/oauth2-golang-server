package utils

import (
	"bytes"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/base64"
	"encoding/hex"
	"encoding/pem"
	"errors"
	"fmt"
	"github.com/golang-jwt/jwt/v4"
	"github.com/manuelrojas19/go-oauth2-server/configuration/keymanager"
	"gopkg.in/square/go-jose.v2"
	"io/ioutil"
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

// LoadRSAPrivateKey loads an RSA private key from a PEM file.
func LoadRSAPrivateKey() (*rsa.PrivateKey, error) {
	keyData, err := ioutil.ReadFile("private_key.pem")
	if err != nil {
		return nil, err
	}

	block, _ := pem.Decode(keyData)
	if block == nil || block.Type != "RSA PRIVATE KEY" {
		return nil, errors.New("invalid key type")
	}

	privateKey, err := x509.ParsePKCS1PrivateKey(block.Bytes)
	if err != nil {
		return nil, err
	}

	return privateKey, nil
}

// GenerateJWT creates a JWT token with the given parameters.
func GenerateJWT(clientId string, userId string, secretKey interface{}, tokenType string) (string, error) {
	var expirationTime time.Duration
	var signingMethod jwt.SigningMethod

	switch tokenType {
	case "access":
		expirationTime = time.Hour
		signingMethod = jwt.SigningMethodRS256
		privateKey, err := keymanager.GetJWTPrivateKey()
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

// GenerateJWE creates a JWE token from the JWT payload using global keys
func GenerateJWE(jwtToken string) (string, error) {
	publicKey, err := keymanager.GetJWTPublicKey()
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
