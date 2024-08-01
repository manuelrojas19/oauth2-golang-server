package repositories_test

import (
	"errors"
	"github.com/manuelrojas19/go-oauth2-server/store"
	"github.com/manuelrojas19/go-oauth2-server/store/repositories"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func TestSave(t *testing.T) {
	// Create an in-memory SQLite database
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Fatalf("Failed to connect to database: %v", err)
	}

	// Automigrate the schema
	err = db.AutoMigrate(&store.AccessToken{})
	if err != nil {
		t.Fatalf("Failed to migrate database: %v", err)
	}

	repo := repositories.NewAccessTokenRepository(db)

	// Define a valid access token
	token := &store.AccessToken{
		ClientId:  "tests-client",
		ExpiresAt: time.Now().Add(1 * time.Hour),
	}

	t.Run("Successful token save", func(t *testing.T) {
		// Save the token
		savedToken, err := repo.Save(token)
		assert.NoError(t, err)
		assert.NotNil(t, savedToken)
	})

	t.Run("Failed to delete expired tokens", func(t *testing.T) {
		// Insert an expired token
		expiredToken := &store.AccessToken{
			ClientId:  "tests-client",
			ExpiresAt: time.Now().Add(-1 * time.Hour),
		}
		db.Create(expiredToken)

		// Mock the delete operation to fail
		db.Exec("PRAGMA foreign_keys = OFF")
		err := db.Delete(&store.AccessToken{}, "client_id = ? AND expires_at <= ?", token.ClientId, time.Now()).Error
		assert.NoError(t, err)
		db.Exec("PRAGMA foreign_keys = ON")

		// Try to save the token
		_, err = repo.Save(token)
		assert.Error(t, err)
	})

	t.Run("Failed to create token", func(t *testing.T) {
		// Close the database connection to force an error on create
		sqlDB, _ := db.DB()
		sqlDB.Close()

		// Try to save the token
		_, err := repo.Save(token)
		assert.Error(t, err)

		// Reopen the database connection
		db, _ = gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
		db.AutoMigrate(&store.AccessToken{})
	})

	t.Run("Failed to commit transaction", func(t *testing.T) {
		// Use a real database but mock commit to fail
		db.Exec("PRAGMA foreign_keys = OFF")
		err := db.Transaction(func(tx *gorm.DB) error {
			tx.Create(token)
			return errors.New("commit error")
		})
		assert.Error(t, err)
		db.Exec("PRAGMA foreign_keys = ON")
	})
}
