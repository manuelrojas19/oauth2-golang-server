package services

import (
	"crypto/rsa"
	"crypto/sha1"
	"encoding/hex"
	"fmt"
	"sync"
	"time"

	"github.com/lestrrat-go/jwx/jwk"
	"github.com/manuelrojas19/go-oauth2-server/configuration"
)

type wellKnownService struct {
	jwkSetCache *jwkCache
	once        sync.Once
}

type jwkCache struct {
	set        *jwk.Set
	expiration time.Time
}

func NewWellKnownService() WellKnownService {
	return &wellKnownService{}
}

// GetJwk retrieves the JWK set containing the public key for JWT.
func (w *wellKnownService) GetJwk() (*jwk.Set, error) {
	w.once.Do(func() {
		w.jwkSetCache = &jwkCache{}
	})

	if w.jwkSetCache.set != nil && time.Now().Before(w.jwkSetCache.expiration) {
		return w.jwkSetCache.set, nil
	}

	set := jwk.NewSet()

	publicKey, err := configuration.GetJWTPublicKey()
	if err != nil {
		return nil, fmt.Errorf("failed to get JWT public key: %w", err)
	}

	jwtKey, err := jwk.New(publicKey)
	if err != nil {
		return nil, fmt.Errorf("failed to create JWK from public key: %w", err)
	}

	params := map[string]interface{}{
		jwk.KeyIDKey:     calculateKid(publicKey),
		jwk.KeyUsageKey:  "sig",
		jwk.AlgorithmKey: "RS256",
	}

	if err := setJWKParameters(jwtKey, params); err != nil {
		return nil, fmt.Errorf("failed to set JWK parameters: %w", err)
	}

	set.Add(jwtKey)

	// Cache for 5 minutes (or a configurable duration)
	w.jwkSetCache.set = &set
	w.jwkSetCache.expiration = time.Now().Add(5 * time.Minute)

	return &set, nil
}

// calculateKid generates a key ID based on the public key.
func calculateKid(publicKey *rsa.PublicKey) string {
	keyData := publicKey.N.Bytes()
	hasher := sha1.New()
	hasher.Write(keyData)
	return hex.EncodeToString(hasher.Sum(nil))
}

// setJWKParameters sets multiple parameters on a JWK key and returns any errors encountered.
func setJWKParameters(key jwk.Key, params map[string]interface{}) error {
	for k, v := range params {
		if err := key.Set(k, v); err != nil {
			return err
		}
	}
	return nil
}
