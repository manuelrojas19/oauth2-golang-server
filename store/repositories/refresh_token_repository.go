package repositories

import (
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/manuelrojas19/go-oauth2-server/store"
	"gorm.io/gorm"
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
			log.Printf("PANIC: Rolled back transaction for access token ScopeId '%s' due to: %v", tokenId, r)
		}
	}()

	if err := tx.Error; err != nil {
		log.Printf("ERROR: Failed to start transaction for access token ScopeId '%s': %v", tokenId, err)
		return fmt.Errorf("failed to start transaction: %w", err)
	}

	usedDeleteTokenQuery := tx.Unscoped().Where("access_token_id = ?", tokenId)
	if err := usedDeleteTokenQuery.Delete(new(store.RefreshToken)).Error; err != nil {
		log.Printf("ERROR: Failed to execute delete query for access token ScopeId '%s': %v", tokenId, err)
		tx.Rollback()
		return fmt.Errorf("failed to delete refresh token: %w", err)
	}

	if err := tx.Commit().Error; err != nil {
		log.Printf("ERROR: Error committing transaction for access token ScopeId '%s': %v", tokenId, err)
		tx.Rollback()
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	log.Printf("INFO: Successfully deleted refresh tokens for access token ScopeId '%s'", tokenId)
	return nil
}

func (ot *refreshTokenRepository) Save(token *store.RefreshToken) (*store.RefreshToken, error) {
	log.Printf("Starting transaction to save refresh token for access_token_id %s", token.AccessTokenId)

	tx := ot.Db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
			log.Printf("PANIC: Rolled back transaction for access token ScopeId '%s' due to: %v", token.AccessTokenId, r)
		}
	}()

	if err := tx.Error; err != nil {
		log.Printf("ERROR: Failed to start transaction for access token ScopeId '%s': %v", token.AccessTokenId, err)
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
func (ot *refreshTokenRepository) FindByToken(token string) (*store.RefreshToken, error) {
	log.Printf("Searching for refresh token with token string %s", token)

	// Initialize a new RefreshToken entity
	refreshToken := new(store.RefreshToken)

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

// FindByRefreshToken retrieves a refresh token from the database using the token string.
func (ot *refreshTokenRepository) FindByRefreshToken(refreshToken string) (*store.RefreshToken, error) {
	var token store.RefreshToken

	if err := ot.Db.Where("refresh_token = ?", refreshToken).First(&token).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("refresh token not found")
		}
		return nil, fmt.Errorf("failed to find refresh token: %w", err)
	}

	if token.ExpiresAt.Before(time.Now()) {
		return nil, fmt.Errorf("refresh token expired")
	}

	return &token, nil
}

func (ot *refreshTokenRepository) DeleteByRefreshToken(refreshToken string) error {
	result := ot.Db.Where("refresh_token = ?", refreshToken).Delete(&store.RefreshToken{})
	if result.Error != nil {
		return fmt.Errorf("failed to delete refresh token: %w", result.Error)
	}
	// If no rows were affected, it means the token was not found, but we don't return an error as per RFC 7009
	return nil
}
