package repositories

import (
	"errors"
	"github.com/manuelrojas19/go-oauth2-server/store"
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
	err = db.AutoMigrate(&store.Scope{}, &store.User{}, &store.OauthClient{}, &store.AuthCode{})
	assert.NoError(t, err)

	return db
}

func TestAuthCodeRepository_Save(t *testing.T) {
	db := setupTestDB(t)
	repo := NewAuthCodeRepository(db)

	authCode := &store.AuthCode{
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

	savedAuthCode, err := repo.Save(authCode)
	assert.NoError(t, err)
	assert.NotNil(t, savedAuthCode)
	assert.Equal(t, authCode.Id, savedAuthCode.Id)

	var result store.AuthCode
	err = db.First(&result, "id = ?", authCode.Id).Error
	assert.NoError(t, err)
	assert.Equal(t, authCode.Code, result.Code)
}

func TestAuthCodeRepository_FindByCode(t *testing.T) {
	db := setupTestDB(t)
	repo := NewAuthCodeRepository(db)

	authCode := &store.AuthCode{
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

	result, err := repo.FindByCode(authCode.Code)
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, authCode.Code, result.Code)
}

func TestAuthCodeRepository_FindByCode_RecordNotFound(t *testing.T) {
	db := setupTestDB(t)
	repo := NewAuthCodeRepository(db)

	result, err := repo.FindByCode("non_existent_code")
	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Equal(t, "AuthCode not found or invalidated", err.Error())
}

func TestAuthCodeRepository_FindByCode_GenericError(t *testing.T) {
	db := setupTestDB(t)
	repo := NewAuthCodeRepository(db)

	// Simulate a generic database error
	db.Error = errors.New("generic database error") // Simulate a generic database error

	result, err := repo.FindByCode("some_code")
	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Equal(t, "error finding AuthCode: generic database error", err.Error())
}
