package keymanager

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
			log.Printf("Failed to load keys: %v", err)
			// If loading keys fails, generate new keys
			PrivateKey, err = rsa.GenerateKey(rand.Reader, 2048)
			if err != nil {
				log.Printf("Failed to generate RSA key: %v", err)
				return
			}
			PublicKey = &PrivateKey.PublicKey
			err = SaveKeys(PrivateKey, PublicKey)
			if err != nil {
				log.Printf("Failed to save keys: %v", err)
			}
		}
	})
	return err
}

// GetPublicKey returns the public key
func GetPublicKey() (*rsa.PublicKey, error) {
	if PublicKey == nil {
		err := fmt.Errorf("public key is not initialized")
		log.Println(err)
		return nil, err
	}

	return PublicKey, nil
}

// GetPrivateKey returns the private key
func GetPrivateKey() (*rsa.PrivateKey, error) {
	if PrivateKey == nil {
		err := fmt.Errorf("private key is not initialized")
		log.Println(err)
		return nil, err
	}
	return PrivateKey, nil
}

// SaveKeys saves the private and public keys to disk
func SaveKeys(privateKey *rsa.PrivateKey, publicKey *rsa.PublicKey) error {
	privFile, err := os.Create("private_key.pem")
	if err != nil {
		log.Printf("Failed to create private key file: %v", err)
		return err
	}
	defer privFile.Close()

	privBytes := x509.MarshalPKCS1PrivateKey(privateKey)
	err = pem.Encode(privFile, &pem.Block{Type: "RSA PRIVATE KEY", Bytes: privBytes})
	if err != nil {
		log.Printf("Failed to encode private key: %v", err)
		return err
	}

	pubFile, err := os.Create("public_key.pem")
	if err != nil {
		log.Printf("Failed to create public key file: %v", err)
		return err
	}
	defer pubFile.Close()

	pubBytes, err := x509.MarshalPKIXPublicKey(publicKey)
	if err != nil {
		log.Printf("Failed to marshal public key: %v", err)
		return err
	}
	err = pem.Encode(pubFile, &pem.Block{Type: "PUBLIC KEY", Bytes: pubBytes})
	if err != nil {
		log.Printf("Failed to encode public key: %v", err)
		return err
	}

	return nil
}

// LoadKeys loads the private and public keys from disk
func LoadKeys() error {
	privFile, err := os.ReadFile("private_key.pem")
	if err != nil {
		log.Printf("Failed to read private key file: %v", err)
		return err
	}
	privBlock, _ := pem.Decode(privFile)
	if privBlock == nil || privBlock.Type != "RSA PRIVATE KEY" {
		err := fmt.Errorf("failed to decode PEM block containing private key")
		log.Println(err)
		return err
	}
	privateKey, err := x509.ParsePKCS1PrivateKey(privBlock.Bytes)
	if err != nil {
		log.Printf("Failed to parse private key: %v", err)
		return err
	}
	PrivateKey = privateKey

	pubFile, err := os.ReadFile("public_key.pem")
	if err != nil {
		log.Printf("Failed to read public key file: %v", err)
		return err
	}
	pubBlock, _ := pem.Decode(pubFile)
	if pubBlock == nil || pubBlock.Type != "PUBLIC KEY" {
		err := fmt.Errorf("failed to decode PEM block containing public key")
		log.Println(err)
		return err
	}
	pubKey, err := x509.ParsePKIXPublicKey(pubBlock.Bytes)
	if err != nil {
		log.Printf("Failed to parse public key: %v", err)
		return err
	}
	PublicKey = pubKey.(*rsa.PublicKey)

	return nil
}
