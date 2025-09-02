package repositories

import (
	"errors"
	"fmt"

	"github.com/manuelrojas19/go-oauth2-server/store"
	"github.com/manuelrojas19/go-oauth2-server/utils"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type authCodeRepository struct {
	Db     *gorm.DB
	logger *zap.Logger
}

// NewAuthCodeRepository initializes a new AuthCodeRepository
func NewAuthCodeRepository(db *gorm.DB, logger *zap.Logger) AuthorizationRepository {
	return &authCodeRepository{
		Db:     db,
		logger: logger,
	}
}

// Save saves an AuthCode to the database
func (r *authCodeRepository) Save(authCode *store.AuthCode) (*store.AuthCode, error) {
	r.logger.Info("Attempting to save authorization code", zap.String("code", authCode.Code), zap.String("clientId", *authCode.ClientId))
	r.logger.Debug("AuthCode entity to save", zap.Any("authCode", authCode))

	if err := r.Db.Create(authCode).Error; err != nil {
		r.logger.Error("Error saving AuthCode to database",
			zap.String("code", authCode.Code),
			zap.String("clientId", utils.StringDeref(authCode.ClientId)),
			zap.Error(err),
			zap.Stack("stacktrace"),
		)
		return nil, fmt.Errorf("failed to save AuthCode: %w", err)
	}
	r.logger.Info("AuthCode saved successfully", zap.String("code", authCode.Code), zap.String("id", authCode.Id))
	return authCode, nil
}

// FindByCode retrieves an AuthCode from the database using the code string
func (r *authCodeRepository) FindByCode(code string) (*store.AuthCode, error) {
	r.logger.Info("Searching for AuthCode by code",
		zap.String("code", code),
	)
	r.logger.Debug("Executing database query to find AuthCode")

	// Initialize a new AuthCode entity
	authCode := new(store.AuthCode)

	// Query the database for the code
	result := r.Db.Where("code = ?", code).First(authCode)

	// Handle errors during the query
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			r.logger.Debug("AuthCode not found in database",
				zap.String("code", code),
				zap.Error(result.Error),
			)
			return nil, fmt.Errorf("AuthCode not found or invalidated: %w", result.Error)
		}
		r.logger.Error("Error finding AuthCode in database",
			zap.String("code", code),
			zap.Error(result.Error),
			zap.Stack("stacktrace"),
		)
		return nil, fmt.Errorf("error finding AuthCode: %w", result.Error)
	}

	r.logger.Info("Successfully found AuthCode",
		zap.String("code", code),
		zap.String("authCodeId", authCode.Id),
	)
	r.logger.Debug("Found AuthCode details", zap.Any("authCode", authCode))

	return authCode, nil
}

// Delete deletes an AuthCode from the database using the code string
func (r *authCodeRepository) Delete(code string) error {
	r.logger.Info("Attempting to delete AuthCode",
		zap.String("code", code),
	)
	r.logger.Debug("Executing database delete operation for AuthCode")

	result := r.Db.Where("code = ?", code).Delete(&store.AuthCode{})
	if result.Error != nil {
		r.logger.Error("Error deleting AuthCode from database",
			zap.String("code", code),
			zap.Error(result.Error),
			zap.Stack("stacktrace"),
		)
		return fmt.Errorf("failed to delete AuthCode: %w", result.Error)
	}

	if result.RowsAffected == 0 {
		r.logger.Warn("AuthCode not found for deletion",
			zap.String("code", code),
		)
	} else {
		r.logger.Info("AuthCode deleted successfully",
			zap.String("code", code),
			zap.Int64("rowsAffected", result.RowsAffected),
		)
	}

	return nil
}
