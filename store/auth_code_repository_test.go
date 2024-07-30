package store

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func setupTestDB(t *testing.T) *gorm.DB {
	// Use an in-memory SQLite database for testing
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	assert.NoError(t, err)

	// Migrate the schema
	err = db.AutoMigrate(&Scope{}, &User{}, &OauthClient{}, &AuthCode{})
	assert.NoError(t, err)

	return db
}

func TestSave(t *testing.T) {
	db := setupTestDB(t)

	// Initialize the repository with the test database
	repo := NewAuthCodeRepository(db)

	// Create a new AuthCode instance
	authCode := &AuthCode{
		Id:          "1",
		Code:        "test_code",
		RedirectURI: "http://localhost/callback",
		Scope:       "read",
		Used:        false,
		UserId:      "user_1",
		ClientId:    "client_1",
		ExpiresAt:   time.Now().Add(1 * time.Hour),
		CreatedAt:   time.Now(),
	}

	// Call the Save method
	result, err := repo.Save(authCode)
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, authCode.Id, result.Id)

	// Verify that the record was saved
	var savedAuthCode AuthCode
	err = db.First(&savedAuthCode, "id = ?", authCode.Id).Error
	assert.NoError(t, err)
	assert.Equal(t, authCode.Code, savedAuthCode.Code)
}

func TestFindByCode(t *testing.T) {
	db := setupTestDB(t)

	// Initialize the repository with the test database
	repo := NewAuthCodeRepository(db)

	// Create a new AuthCode instance and save it to the database
	authCode := &AuthCode{
		Id:          "1",
		Code:        "test_code",
		RedirectURI: "http://localhost/callback",
		Scope:       "read",
		Used:        false,
		UserId:      "user_1",
		ClientId:    "client_1",
		ExpiresAt:   time.Now().Add(1 * time.Hour),
		CreatedAt:   time.Now(),
	}
	err := db.Create(authCode).Error
	assert.NoError(t, err)

	// Call the FindByCode method
	result, err := repo.FindByCode(authCode.Code)
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, authCode.Id, result.Id)
}
