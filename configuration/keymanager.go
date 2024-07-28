package configuration

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"log"
	"os"
	"sync"
)

var (
	JWTGenerationKeys KeyPair
	JWEGenerationKeys KeyPair
	once              sync.Once
)

// KeyPair holds a private key and its corresponding public key.
type KeyPair struct {
	PrivateKey *rsa.PrivateKey
	PublicKey  *rsa.PublicKey
}

// Initialize initializes the JWT and JWE keys. It is thread-safe and will only run once.
func Initialize() error {
	var err error
	once.Do(func() {
		if err = loadKeys(); err != nil {
			log.Printf("Failed to load keys: %v", err)
			err = generateAndSaveKeys()
		}
	})
	return err
}

// GetJWTPublicKey returns the public key for JWT.
func GetJWTPublicKey() (*rsa.PublicKey, error) {
	return getPublicKey(JWTGenerationKeys.PublicKey, "JWT")
}

// GetJWTPrivateKey returns the private key for JWT.
func GetJWTPrivateKey() (*rsa.PrivateKey, error) {
	return getPrivateKey(JWTGenerationKeys.PrivateKey, "JWT")
}

// GetJWEPublicKey returns the public key for JWE.
func GetJWEPublicKey() (*rsa.PublicKey, error) {
	return getPublicKey(JWEGenerationKeys.PublicKey, "JWE")
}

// GetJWEPrivateKey returns the private key for JWE.
func GetJWEPrivateKey() (*rsa.PrivateKey, error) {
	return getPrivateKey(JWEGenerationKeys.PrivateKey, "JWE")
}

func getPublicKey(key *rsa.PublicKey, keyType string) (*rsa.PublicKey, error) {
	if key == nil {
		return nil, fmt.Errorf("%s public key is not initialized", keyType)
	}
	return key, nil
}

func getPrivateKey(key *rsa.PrivateKey, keyType string) (*rsa.PrivateKey, error) {
	if key == nil {
		return nil, fmt.Errorf("%s private key is not initialized", keyType)
	}
	return key, nil
}

func generateAndSaveKeys() error {
	var err error

	// Generate JWT keys
	JWTGenerationKeys.PrivateKey, err = rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		return fmt.Errorf("failed to generate JWT private key: %w", err)
	}
	JWTGenerationKeys.PublicKey = &JWTGenerationKeys.PrivateKey.PublicKey

	// Generate JWE keys
	JWEGenerationKeys.PrivateKey, err = rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		return fmt.Errorf("failed to generate JWE private key: %w", err)
	}
	JWEGenerationKeys.PublicKey = &JWEGenerationKeys.PrivateKey.PublicKey

	if err = saveKeys(); err != nil {
		return fmt.Errorf("failed to save keys: %w", err)
	}

	return nil
}

func saveKeys() error {
	if err := saveKeyToFile("jwt_private_key.pem", JWTGenerationKeys.PrivateKey); err != nil {
		return err
	}
	if err := savePublicKeyToFile("jwt_public_key.pem", JWTGenerationKeys.PublicKey); err != nil {
		return err
	}

	if err := saveKeyToFile("jwe_private_key.pem", JWEGenerationKeys.PrivateKey); err != nil {
		return err
	}
	if err := savePublicKeyToFile("jwe_public_key.pem", JWEGenerationKeys.PublicKey); err != nil {
		return err
	}

	return nil
}

func saveKeyToFile(filename string, key *rsa.PrivateKey) error {
	return saveToFile(filename, "RSA PRIVATE KEY", x509.MarshalPKCS1PrivateKey(key))
}

func savePublicKeyToFile(filename string, key *rsa.PublicKey) error {
	pubBytes, err := x509.MarshalPKIXPublicKey(key)
	if err != nil {
		return fmt.Errorf("failed to marshal public key: %w", err)
	}
	return saveToFile(filename, "PUBLIC KEY", pubBytes)
}

func saveToFile(filename, blockType string, bytes []byte) error {
	file, err := os.Create(filename)
	if err != nil {
		return fmt.Errorf("failed to create file %s: %w", filename, err)
	}
	defer file.Close()

	if err := pem.Encode(file, &pem.Block{Type: blockType, Bytes: bytes}); err != nil {
		return fmt.Errorf("failed to encode %s in %s: %w", blockType, filename, err)
	}
	return nil
}

func loadKeys() error {
	if err := loadPrivateKeyFromFile("jwt_private_key.pem", &JWTGenerationKeys.PrivateKey); err != nil {
		return err
	}
	if err := loadPublicKeyFromFile("jwt_public_key.pem", &JWTGenerationKeys.PublicKey); err != nil {
		return err
	}

	if err := loadPrivateKeyFromFile("jwe_private_key.pem", &JWEGenerationKeys.PrivateKey); err != nil {
		return err
	}
	if err := loadPublicKeyFromFile("jwe_public_key.pem", &JWEGenerationKeys.PublicKey); err != nil {
		return err
	}
	return nil
}

func loadPrivateKeyFromFile(filename string, key **rsa.PrivateKey) error {
	bytes, err := readFile(filename)
	if err != nil {
		return err
	}

	block, _ := pem.Decode(bytes)
	if block == nil || block.Type != "RSA PRIVATE KEY" {
		return fmt.Errorf("invalid or missing PEM block in %s", filename)
	}

	privKey, err := x509.ParsePKCS1PrivateKey(block.Bytes)
	if err != nil {
		return fmt.Errorf("failed to parse private key from %s: %w", filename, err)
	}
	*key = privKey
	return nil
}

func loadPublicKeyFromFile(filename string, key **rsa.PublicKey) error {
	bytes, err := readFile(filename)
	if err != nil {
		return err
	}

	block, _ := pem.Decode(bytes)
	if block == nil || block.Type != "PUBLIC KEY" {
		return fmt.Errorf("invalid or missing PEM block in %s", filename)
	}

	pubKey, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		return fmt.Errorf("failed to parse public key from %s: %w", filename, err)
	}
	*key = pubKey.(*rsa.PublicKey)
	return nil
}

func readFile(filename string) ([]byte, error) {
	bytes, err := os.ReadFile(filename)
	if err != nil {
		return nil, fmt.Errorf("failed to read file %s: %w", filename, err)
	}
	return bytes, nil
}
