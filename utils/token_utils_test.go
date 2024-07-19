package utils_test

import (
	"testing"
	"time"

	"github.com/manuelrojas19/go-oauth2-server/utils"
)

func TestToken(t *testing.T) {
	// Test inputs
	clientId := "testClientId"
	userId := "testUserId"
	createdAt := time.Date(2024, time.July, 14, 12, 0, 0, 0, time.UTC)

	// Expected output format (we cannot predict the exact token, but we can check its length)
	token, err := utils.GenerateToken(clientId, userId, createdAt)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	// Check if the token is not empty
	if token == "" {
		t.Fatal("Expected a token, got an empty string")
	}

	// Check if the token has the expected length of a SHA256 hash encoded in Base64 (44 chars when trimmed)
	expectedTokenLength := 43 // 44 chars in base64, with one '=' padding removed
	if len(token) != expectedTokenLength {
		t.Fatalf("Expected token length to be %d, got %d", expectedTokenLength, len(token))
	}

	// Edge case: empty clientId and userId
	emptyClientId := ""
	emptyUserId := ""
	token, err = utils.GenerateToken(emptyClientId, emptyUserId, createdAt)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	if token == "" {
		t.Fatal("Expected a token, got an empty string")
	}

	// Edge case: extreme date (Unix epoch)
	epochTime := time.Unix(0, 0)
	token, err = utils.GenerateToken(clientId, userId, epochTime)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	if token == "" {
		t.Fatal("Expected a token, got an empty string")
	}
}
