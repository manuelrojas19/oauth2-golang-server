package repositories

import (
	"errors"
	"fmt"
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

func (ot *refreshTokenRepository) InvalidateRefreshTokensByAccessTokenId(tokenId string) error {
	tx := ot.Db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
			log.Printf("PANIC: Rolled back transaction for access token ID '%s' due to: %v", tokenId, r)
		}
	}()

	if err := tx.Error; err != nil {
		log.Printf("ERROR: Failed to start transaction for access token ID '%s': %v", tokenId, err)
		return fmt.Errorf("failed to start transaction: %w", err)
	}

	usedDeleteTokenQuery := tx.Unscoped().Where("access_token_id = ?", tokenId)
	if err := usedDeleteTokenQuery.Delete(new(entities.RefreshToken)).Error; err != nil {
		log.Printf("ERROR: Failed to execute delete query for access token ID '%s': %v", tokenId, err)
		tx.Rollback()
		return fmt.Errorf("failed to delete refresh token: %w", err)
	}

	if err := tx.Commit().Error; err != nil {
		log.Printf("ERROR: Error committing transaction for access token ID '%s': %v", tokenId, err)
		tx.Rollback()
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	log.Printf("INFO: Successfully deleted refresh tokens for access token ID '%s'", tokenId)
	return nil
}

func (ot *refreshTokenRepository) Save(token *entities.RefreshToken) (*entities.RefreshToken, error) {
	log.Printf("Starting transaction to save refresh token for access_token_id %s", token.AccessTokenId)

	tx := ot.Db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
			log.Printf("PANIC: Rolled back transaction for access token ID '%s' due to: %v", token.AccessTokenId, r)
		}
	}()

	if err := tx.Error; err != nil {
		log.Printf("ERROR: Failed to start transaction for access token ID '%s': %v", token.AccessTokenId, err)
		return nil, fmt.Errorf("failed to start transaction: %w", err)
	}

	if err := tx.Create(token).Error; err != nil {
		log.Printf("ERROR: Error creating new refresh token for access_token_id %s: %v", token.AccessTokenId, err)
		tx.Rollback()
		return nil, fmt.Errorf("failed to create refresh token: %w", err)
	}

	if err := tx.Commit().Error; err != nil {
		log.Printf("ERROR: Error committing transaction for new refresh token with access_token_id %s: %v", token.AccessTokenId, err)
		tx.Rollback()
		return nil, fmt.Errorf("failed to commit transaction: %w", err)
	}

	log.Printf("INFO: Successfully saved new refresh token for access_token_id %s", token.AccessTokenId)
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
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("RefreshToken not found or invalidated")
		}
		return nil, fmt.Errorf("error finding Refresh Token: %w", result.Error)
	}

	log.Printf("Successfully found refresh token with token string %s", token)
	return refreshToken, nil
}
