package services

import (
	"crypto/rsa"
	"crypto/sha1"
	"encoding/hex"
	"github.com/lestrrat-go/jwx/jwk"
	"github.com/manuelrojas19/go-oauth2-server/configuration"
)

type wellKnownService struct {
}

func NewWellKnownService() WellKnownService {
	return &wellKnownService{}
}

// GetJwk retrieves the JWK set containing the public key for JWT.
func (w wellKnownService) GetJwk() (*jwk.Set, error) {
	set := jwk.NewSet()

	publicKey, err := configuration.GetJWTPublicKey()
	if err != nil {
		return nil, err
	}

	jwtKey, err := jwk.New(publicKey)
	if err != nil {
		return nil, err
	}

	params := map[string]interface{}{
		jwk.KeyIDKey:     calculateKid(publicKey),
		jwk.KeyUsageKey:  "sig",
		jwk.AlgorithmKey: "RS256",
	}

	if err := setJWKParameters(jwtKey, params); err != nil {
		return nil, err
	}

	set.Add(jwtKey)
	return &set, nil
}

// calculateKid generates a key ScopeId based on the public key.
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
