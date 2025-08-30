package utils

import (
	"crypto/sha256"
	"encoding/base64"
)

// S256Challenge generates a code challenge using S256 method from a code verifier.
func S256Challenge(codeVerifier string) string {
	s := sha256.New()
	s.Write([]byte(codeVerifier))
	return base64.RawURLEncoding.EncodeToString(s.Sum(nil))
}
