package keymanager

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"os"
	"sync"
)

var (
	PrivateKey *rsa.PrivateKey
	PublicKey  *rsa.PublicKey
	once       sync.Once
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

// GetPrivateKey returns the public key
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
	pem.Encode(privFile, &pem.Block{Type: "RSA PRIVATE KEY", Bytes: privBytes})

	pubFile, err := os.Create("public_key.pem")
	if err != nil {
		return err
	}
	defer pubFile.Close()

	pubBytes, err := x509.MarshalPKIXPublicKey(publicKey)
	if err != nil {
		return err
	}
	pem.Encode(pubFile, &pem.Block{Type: "PUBLIC KEY", Bytes: pubBytes})

	return nil
}

// LoadKeys loads the private and public keys from disk
func LoadKeys() error {
	privFile, err := os.ReadFile("private_key.pem")
	if err != nil {
		return err
	}
	privBlock, _ := pem.Decode(privFile)
	if privBlock == nil || privBlock.Type != "RSA PRIVATE KEY" {
		return fmt.Errorf("failed to decode PEM block containing private key")
	}
	var privateKey *rsa.PrivateKey
	privateKey, err = x509.ParsePKCS1PrivateKey(privBlock.Bytes)
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
