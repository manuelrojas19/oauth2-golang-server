package repositories

import (
	"errors"
	"fmt"
	"time"

	"github.com/manuelrojas19/go-oauth2-server/store"
	"github.com/manuelrojas19/go-oauth2-server/utils"

	"go.uber.org/zap"
	"gorm.io/gorm"
)

type accessTokenRepository struct {
	Db     *gorm.DB
	logger *zap.Logger
}

func NewAccessTokenRepository(db *gorm.DB, logger *zap.Logger) AccessTokenRepository {
	return &accessTokenRepository{Db: db, logger: logger}
}

func (ot *accessTokenRepository) Save(token *store.AccessToken) (*store.AccessToken, error) {
	ot.logger.Info("Starting transaction to save access token", zap.String("clientId", utils.StringDeref(token.ClientId)))
	ot.logger.Debug("Access token details to be saved", zap.Any("token", token))

	tx := ot.Db.Begin()
	if tx.Error != nil {
		ot.logger.Error("Failed to begin transaction for saving access token", zap.Error(tx.Error))
		return nil, fmt.Errorf("failed to begin transaction: %w", tx.Error)
	}

	ot.logger.Debug("Deleting expired access tokens for client", zap.String("clientId", utils.StringDeref(token.ClientId)))
	expiredTokensQuery := tx.Unscoped().
		Where("client_id = ?", token.ClientId).
		Where("expires_at <= ?", time.Now())

	if err := expiredTokensQuery.Delete(new(store.AccessToken)).Error; err != nil {
		ot.logger.Error("Error deleting expired access tokens", zap.String("clientId", utils.StringDeref(token.ClientId)), zap.Error(err))
		tx.Rollback()
		return nil, fmt.Errorf("failed to delete expired access tokens: %w", err)
	}
	ot.logger.Debug("Expired access tokens deleted if any")

	ot.logger.Debug("Creating new access token record in database")
	if err := tx.Create(token).Error; err != nil {
		ot.logger.Error("Error creating access token", zap.String("clientId", utils.StringDeref(token.ClientId)), zap.Error(err))
		tx.Rollback()
		return nil, fmt.Errorf("failed to create access token: %w", err)
	}
	ot.logger.Debug("Access token record created in database", zap.String("accessTokenId", token.Id))

	if err := tx.Commit().Error; err != nil {
		ot.logger.Error("Error committing transaction for saving access token", zap.String("clientId", utils.StringDeref(token.ClientId)), zap.Error(err))
		tx.Rollback()
		return nil, fmt.Errorf("failed to commit transaction: %w", err)
	}

	ot.logger.Info("Access token saved successfully", zap.String("clientId", utils.StringDeref(token.ClientId)), zap.String("accessTokenId", token.Id))
	return token, nil
}

func (ot *accessTokenRepository) FindByAccessToken(accessToken string) (*store.AccessToken, error) {
	ot.logger.Info("Attempting to find access token", zap.String("accessToken", accessToken))
	var token store.AccessToken

	if err := ot.Db.Where("access_token = ?", accessToken).First(&token).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			ot.logger.Debug("Access token not found", zap.String("accessToken", accessToken))
			return nil, fmt.Errorf("access token not found")
		}
		ot.logger.Error("Failed to find access token in database", zap.String("accessToken", accessToken), zap.Error(err))
		return nil, fmt.Errorf("failed to find access token: %w", err)
	}
	ot.logger.Debug("Access token found", zap.String("accessTokenId", token.Id))

	if token.ExpiresAt.Before(time.Now()) {
		ot.logger.Warn("Access token expired", zap.String("accessTokenId", token.Id), zap.Time("expiresAt", token.ExpiresAt))
		return nil, fmt.Errorf("access token expired")
	}
	ot.logger.Info("Access token is valid and active", zap.String("accessTokenId", token.Id))

	return &token, nil
}

func (ot *accessTokenRepository) DeleteByAccessToken(accessToken string) error {
	ot.logger.Info("Attempting to delete access token", zap.String("accessToken", accessToken))
	result := ot.Db.Where("access_token = ?", accessToken).Delete(&store.AccessToken{})
	if result.Error != nil {
		ot.logger.Error("Failed to delete access token from database", zap.String("accessToken", accessToken), zap.Error(result.Error))
		return fmt.Errorf("failed to delete access token: %w", result.Error)
	}

	if result.RowsAffected == 0 {
		ot.logger.Info("No access token found to delete", zap.String("accessToken", accessToken))
	} else {
		ot.logger.Info("Access token deleted successfully", zap.String("accessToken", accessToken))
	}
	// If no rows were affected, it means the token was not found, but we don't return an error as per RFC 7009
	return nil
}
