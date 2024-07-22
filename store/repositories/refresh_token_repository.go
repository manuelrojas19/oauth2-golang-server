package repositories

import (
	"errors"
	"github.com/manuelrojas19/go-oauth2-server/store/entities"
	"gorm.io/gorm"
	"log"
)

type refreshTokenRepository struct {
	Db *gorm.DB
}

// NewRefreshTokenRepository initializes a new instance of RefreshTokenRepository.
func NewRefreshTokenRepository(db *gorm.DB) RefreshTokenRepository {
	return &refreshTokenRepository{Db: db}
}

// Save stores a new refresh token in the database after invalidating previous tokens associated with the same access token ID.
func (ot *refreshTokenRepository) Save(token *entities.RefreshToken) (*entities.RefreshToken, error) {
	log.Printf("Starting transaction to save refresh token for access_token_id %s", token.AccessTokenId)

	// Begin a new transaction
	tx := ot.Db.Begin()

	// Invalidate previous refresh tokens associated with the same access token ID
	previousRefreshTokens := tx.Unscoped().
		Where("access_token_id = ?", token.AccessTokenId)

	if err := previousRefreshTokens.Delete(new(entities.RefreshToken)).Error; err != nil {
		log.Printf("Error deleting previous refresh tokens for access_token_id %s: %v", token.AccessTokenId, err)
		tx.Rollback()
		return nil, err
	}

	// Create the new refresh token
	if err := tx.Create(token).Error; err != nil {
		log.Printf("Error creating new refresh token for access_token_id %s: %v", token.AccessTokenId, err)
		tx.Rollback()
		return nil, err
	}

	// Commit the transaction
	if err := tx.Commit().Error; err != nil {
		log.Printf("Error committing transaction for new refresh token with access_token_id %s: %v", token.AccessTokenId, err)
		tx.Rollback()
		return nil, err
	}

	log.Printf("Successfully saved new refresh token for access_token_id %s", token.AccessTokenId)
	return token, nil
}

// FindByToken retrieves a refresh token from the database using the token string.
func (ot *refreshTokenRepository) FindByToken(token string) (*entities.RefreshToken, error) {
	log.Printf("Searching for refresh token with token string %s", token)

	// Initialize a new RefreshToken entity
	refreshToken := new(entities.RefreshToken)

	// Query the database for the token
	result := ot.Db.Where("token = ?", token).First(refreshToken)

	// Handle errors during the query
	if result.Error != nil {
		log.Printf("Error finding refresh token with token string %s: %v", token, result.Error)
		return nil, result.Error
	}

	// Handle case where no token was found
	if result.RowsAffected == 0 {
		err := errors.New("refresh token not found")
		log.Printf("No refresh token found with token string %s: %v", token, err)
		return nil, err
	}

	log.Printf("Successfully found refresh token with token string %s", token)
	return refreshToken, nil
}
