package utils

import (
	"strings"

	"github.com/manuelrojas19/go-oauth2-server/store"
)

func StringDeref(s *string) string {
	if s != nil {
		return *s
	}
	return ""
}

// JoinStringSlice joins a slice of strings with a separator.
func JoinStringSlice(s []string, sep string) string {
	return strings.Join(s, sep)
}

// ScopesToStringSlice converts a slice of store.Scope to a slice of scope names (strings).
func ScopesToStringSlice(scopes []store.Scope) []string {
	names := make([]string, len(scopes))
	for i, scope := range scopes {
		names[i] = scope.Name
	}
	return names
}
