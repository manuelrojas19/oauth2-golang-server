package repositories

import (
	"errors"
	"fmt"
	"time"

	"github.com/manuelrojas19/go-oauth2-server/store"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type refreshTokenRepository struct {
	Db     *gorm.DB
	logger *zap.Logger
}

// NewRefreshTokenRepository initializes a new instance of RefreshTokenRepository.
func NewRefreshTokenRepository(db *gorm.DB, logger *zap.Logger) RefreshTokenRepository {
	return &refreshTokenRepository{Db: db, logger: logger}
}

func (ot *refreshTokenRepository) InvalidateRefreshTokensByAccessTokenId(accessTokenId string) error {
	ot.logger.Info("Starting transaction to invalidate refresh tokens by access token ID", zap.String("accessTokenId", accessTokenId))

	tx := ot.Db.Begin()
	if tx.Error != nil {
		ot.logger.Error("Failed to begin transaction for invalidating refresh tokens", zap.String("accessTokenId", accessTokenId), zap.Error(tx.Error))
		return fmt.Errorf("failed to start transaction: %w", tx.Error)
	}

	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
			ot.logger.Error("PANIC: Rolled back transaction for access token ID", zap.String("accessTokenId", accessTokenId), zap.Any("panicReason", r), zap.Stack("stacktrace"))
		}
	}()

	ot.logger.Debug("Deleting refresh tokens associated with access token ID", zap.String("accessTokenId", accessTokenId))
	result := tx.Unscoped().Where("access_token_id = ?", accessTokenId).Delete(new(store.RefreshToken))
	if result.Error != nil {
		ot.logger.Error("ERROR: Failed to delete refresh tokens by access token ID", zap.String("accessTokenId", accessTokenId), zap.Error(result.Error))
		tx.Rollback()
		return fmt.Errorf("failed to delete refresh token: %w", result.Error)
	}
	ot.logger.Info("Successfully deleted refresh tokens", zap.String("accessTokenId", accessTokenId), zap.Int64("rowsAffected", result.RowsAffected))

	if err := tx.Commit().Error; err != nil {
		ot.logger.Error("Error committing transaction for invalidating refresh tokens", zap.String("accessTokenId", accessTokenId), zap.Error(err))
		tx.Rollback()
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	ot.logger.Info("Transaction committed successfully for invalidating refresh tokens", zap.String("accessTokenId", accessTokenId))
	return nil
}

func (ot *refreshTokenRepository) Save(token *store.RefreshToken) (*store.RefreshToken, error) {
	ot.logger.Info("Starting transaction to save refresh token", zap.String("accessTokenId", token.AccessTokenId))
	ot.logger.Debug("Refresh token details to be saved", zap.Any("token", token))

	tx := ot.Db.Begin()
	if tx.Error != nil {
		ot.logger.Error("Failed to begin transaction for saving refresh token", zap.String("accessTokenId", token.AccessTokenId), zap.Error(tx.Error))
		return nil, fmt.Errorf("failed to start transaction: %w", tx.Error)
	}

	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
			ot.logger.Error("PANIC: Rolled back transaction for access token ID", zap.String("accessTokenId", token.AccessTokenId), zap.Any("panicReason", r), zap.Stack("stacktrace"))
		}
	}()

	ot.logger.Debug("Creating new refresh token record in database")
	if err := tx.Create(token).Error; err != nil {
		ot.logger.Error("Error creating new refresh token", zap.String("accessTokenId", token.AccessTokenId), zap.Error(err))
		tx.Rollback()
		return nil, fmt.Errorf("failed to create refresh token: %w", err)
	}
	ot.logger.Debug("Refresh token record created in database", zap.String("refreshTokenId", token.Id))

	if err := tx.Commit().Error; err != nil {
		ot.logger.Error("Error committing transaction for new refresh token", zap.String("accessTokenId", token.AccessTokenId), zap.Error(err))
		tx.Rollback()
		return nil, fmt.Errorf("failed to commit transaction: %w", err)
	}

	ot.logger.Info("Successfully saved new refresh token", zap.String("accessTokenId", token.AccessTokenId), zap.String("refreshTokenId", token.Id))
	return token, nil
}

// FindByToken retrieves a refresh token from the database using the token string.
func (ot *refreshTokenRepository) FindByToken(token string) (*store.RefreshToken, error) {
	ot.logger.Info("Searching for refresh token", zap.String("refreshToken", token))
	ot.logger.Debug("Executing database query to find refresh token by token string")

	// Initialize a new RefreshToken entity
	refreshToken := new(store.RefreshToken)

	// Query the database for the token
	result := ot.Db.Where("token = ?", token).First(refreshToken)

	// Handle errors during the query
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			ot.logger.Debug("Refresh token not found in database", zap.String("refreshToken", token))
			return nil, fmt.Errorf("RefreshToken not found or invalidated: %w", result.Error)
		}
		ot.logger.Error("Error finding Refresh Token in database", zap.String("refreshToken", token), zap.Error(result.Error))
		return nil, fmt.Errorf("error finding Refresh Token: %w", result.Error)
	}

	ot.logger.Info("Successfully found refresh token", zap.String("refreshToken", token), zap.String("refreshTokenId", refreshToken.Id))
	ot.logger.Debug("Found Refresh Token details", zap.Any("refreshTokenEntity", refreshToken))

	return refreshToken, nil
}

// FindByRefreshToken retrieves a refresh token from the database using the token string.
func (ot *refreshTokenRepository) FindByRefreshToken(refreshToken string) (*store.RefreshToken, error) {
	ot.logger.Info("Searching for refresh token by refresh token string", zap.String("refreshToken", refreshToken))
	var token store.RefreshToken

	if err := ot.Db.Where("token = ?", refreshToken).First(&token).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			ot.logger.Debug("Refresh token not found", zap.String("refreshToken", refreshToken))
			return nil, fmt.Errorf("refresh token not found")
		}
		ot.logger.Error("Failed to find refresh token in database", zap.String("refreshToken", refreshToken), zap.Error(err))
		return nil, fmt.Errorf("failed to find refresh token: %w", err)
	}
	ot.logger.Debug("Refresh token found", zap.String("refreshTokenId", token.Id))

	if token.ExpiresAt.Before(time.Now()) {
		ot.logger.Warn("Refresh token expired", zap.String("refreshTokenId", token.Id), zap.Time("expiresAt", token.ExpiresAt))
		return nil, fmt.Errorf("refresh token expired")
	}
	ot.logger.Info("Refresh token is valid and active", zap.String("refreshTokenId", token.Id))

	return &token, nil
}

func (ot *refreshTokenRepository) DeleteByRefreshToken(refreshToken string) error {
	ot.logger.Info("Attempting to delete refresh token by refresh token string", zap.String("refreshToken", refreshToken))
	result := ot.Db.Where("refresh_token = ?", refreshToken).Delete(&store.RefreshToken{})
	if result.Error != nil {
		ot.logger.Error("Failed to delete refresh token from database", zap.String("refreshToken", refreshToken), zap.Error(result.Error))
		return fmt.Errorf("failed to delete refresh token: %w", result.Error)
	}

	if result.RowsAffected == 0 {
		ot.logger.Info("No refresh token found to delete", zap.String("refreshToken", refreshToken))
	} else {
		ot.logger.Info("Refresh token deleted successfully", zap.String("refreshToken", refreshToken), zap.Int64("rowsAffected", result.RowsAffected))
	}
	// If no rows were affected, it means the token was not found, but we don't return an error as per RFC 7009
	return nil
}
