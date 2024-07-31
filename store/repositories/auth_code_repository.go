package repositories

import (
	"errors"
	"fmt"
	"github.com/manuelrojas19/go-oauth2-server/store"
	"gorm.io/gorm"
	"log"
)

type authCodeRepository struct {
	Db *gorm.DB
}

func NewAuthCodeRepository(db *gorm.DB) AuthorizationRepository {
	return &authCodeRepository{
		Db: db,
	}
}

func (r *authCodeRepository) Save(authCode *store.AuthCode) (*store.AuthCode, error) {
	if err := r.Db.Create(authCode).Error; err != nil {
		return nil, err
	}
	return authCode, nil
}

// FindByCode retrieves a refresh token from the database using the token string.
func (ot *authCodeRepository) FindByCode(code string) (*store.AuthCode, error) {
	log.Printf("Searching  AuthCode %s", code)

	// Initialize a new RefreshToken entity
	authCode := new(store.AuthCode)

	// Query the database for the token
	result := ot.Db.Where("code = ?", code).First(authCode)

	// Handle errors during the query
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("AuthCode not found or invalidated")
		}
		return nil, fmt.Errorf("error finding AuthCode: %w", result.Error)
	}

	log.Printf("Successfully found AuthCode %s", code)
	return authCode, nil
}