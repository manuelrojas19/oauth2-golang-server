package store

import (
	"errors"
	"fmt"
	"gorm.io/gorm"
	"log"
)

type authorizationCodeRepository struct {
	Db *gorm.DB
}

func NewAuthorizationRepository(db *gorm.DB) AuthorizationRepository {
	return &authorizationCodeRepository{
		Db: db,
	}
}

func (r *authorizationCodeRepository) Save(authCode *AuthCode) (*AuthCode, error) {
	if err := r.Db.Create(authCode).Error; err != nil {
		return nil, err
	}
	return authCode, nil
}

// FindByCode retrieves a refresh token from the database using the token string.
func (ot *authorizationCodeRepository) FindByCode(code string) (*AuthCode, error) {
	log.Printf("Searching  AuthCode %s", code)

	// Initialize a new RefreshToken entity
	authCode := new(AuthCode)

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
