package keymanager

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"golang.org/x/crypto/pbkdf2"
	"io"
	"os"
	"sync"
)

var (
	PrivateKey *rsa.PrivateKey
	PublicKey  *rsa.PublicKey
	once       sync.Once
	passphrase = getPassphrase()            // Retrieve the passphrase from an environment variable
	salt       = []byte("some-random-salt") // Use a random salt in practice, stored securely
)

// Initialize initializes the keys. It is thread-safe and will only run once.
func Initialize() error {
	var err error
	once.Do(func() {
		err = LoadKeys()
		if err != nil {
			// If loading keys fails, generate new keys
			PrivateKey, err = rsa.GenerateKey(rand.Reader, 2048)
			if err != nil {
				return
			}
			PublicKey = &PrivateKey.PublicKey
			err = SaveKeys(PrivateKey, PublicKey)
		}
	})
	return err
}

// GetPublicKey returns the public key
func GetPublicKey() (*rsa.PublicKey, error) {
	if PublicKey == nil {
		return nil, fmt.Errorf("public key is not initialized")
	}
	return PublicKey, nil
}

// GetPrivateKey returns the private key
func GetPrivateKey() (*rsa.PrivateKey, error) {
	if PrivateKey == nil {
		return nil, fmt.Errorf("private key is not initialized")
	}
	return PrivateKey, nil
}

// SaveKeys saves the private and public keys to disk
func SaveKeys(privateKey *rsa.PrivateKey, publicKey *rsa.PublicKey) error {
	privFile, err := os.Create("private_key.pem")
	if err != nil {
		return err
	}
	defer privFile.Close()

	privBytes := x509.MarshalPKCS1PrivateKey(privateKey)
	privBlock := &pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: privBytes,
	}
	// Encrypt the private key with the passphrase
	encPrivBytes, err := encryptWithPassphrase(privBlock.Bytes, passphrase)
	if err != nil {
		return err
	}
	err = pem.Encode(privFile, &pem.Block{Type: "ENCRYPTED RSA PRIVATE KEY", Bytes: encPrivBytes})
	if err != nil {
		return err
	}

	pubFile, err := os.Create("public_key.pem")
	if err != nil {
		return err
	}
	defer pubFile.Close()

	pubBytes, err := x509.MarshalPKIXPublicKey(publicKey)
	if err != nil {
		return err
	}
	err = pem.Encode(pubFile, &pem.Block{Type: "PUBLIC KEY", Bytes: pubBytes})
	if err != nil {
		return err
	}

	return nil
}

// LoadKeys loads the private and public keys from disk
func LoadKeys() error {
	privFile, err := os.ReadFile("private_key.pem")
	if err != nil {
		return err
	}
	privBlock, _ := pem.Decode(privFile)
	if privBlock == nil || privBlock.Type != "ENCRYPTED RSA PRIVATE KEY" {
		return fmt.Errorf("failed to decode PEM block containing encrypted private key")
	}
	privBytes, err := decryptWithPassphrase(privBlock.Bytes, passphrase)
	if err != nil {
		return err
	}
	privBlock, _ = pem.Decode(privBytes)
	if privBlock == nil || privBlock.Type != "RSA PRIVATE KEY" {
		return fmt.Errorf("failed to decode PEM block containing private key")
	}
	privateKey, err := x509.ParsePKCS1PrivateKey(privBlock.Bytes)
	if err != nil {
		return err
	}
	PrivateKey = privateKey

	pubFile, err := os.ReadFile("public_key.pem")
	if err != nil {
		return err
	}
	pubBlock, _ := pem.Decode(pubFile)
	if pubBlock == nil || pubBlock.Type != "PUBLIC KEY" {
		return fmt.Errorf("failed to decode PEM block containing public key")
	}
	pubKey, err := x509.ParsePKIXPublicKey(pubBlock.Bytes)
	if err != nil {
		return err
	}
	PublicKey = pubKey.(*rsa.PublicKey)

	return nil
}

// encryptWithPassphrase encrypts data with a passphrase
func encryptWithPassphrase(data []byte, passphrase string) ([]byte, error) {
	key := pbkdf2.Key([]byte(passphrase), salt, 10000, 32, sha256.New) // 32 bytes key for AES-256
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	// Create a random initialization vector (IV) for AES
	ciphertext := make([]byte, aes.BlockSize+len(data))
	iv := ciphertext[:aes.BlockSize]
	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		return nil, err
	}

	stream := cipher.NewCFBEncrypter(block, iv)
	stream.XORKeyStream(ciphertext[aes.BlockSize:], data)

	return ciphertext, nil
}

// decryptWithPassphrase decrypts data with a passphrase
func decryptWithPassphrase(data []byte, passphrase string) ([]byte, error) {
	key := pbkdf2.Key([]byte(passphrase), salt, 10000, 32, sha256.New) // 32 bytes key for AES-256
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	if len(data) < aes.BlockSize {
		return nil, fmt.Errorf("ciphertext too short")
	}
	iv := data[:aes.BlockSize]
	ciphertext := data[aes.BlockSize:]

	stream := cipher.NewCFBDecrypter(block, iv)
	stream.XORKeyStream(ciphertext, ciphertext)

	return ciphertext, nil
}

// getPassphrase retrieves the passphrase from an environment variable
func getPassphrase() string {
	pass := os.Getenv("KEY_PASSPHRASE")
	if pass == "" {
		return "default-passphrase" // Provide a default value or handle accordingly
	}
	return pass
}
