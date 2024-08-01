package repositories

import (
	"errors"
	"fmt"
	"github.com/manuelrojas19/go-oauth2-server/store"
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
	if err := r.Db.Create(authCode).Error; err != nil {
		r.logger.Error("Error saving AuthCode",
			zap.String("code", authCode.Code),
			zap.Error(err),
			zap.Stack("stacktrace"),
		)
		return nil, fmt.Errorf("failed to save AuthCode: %w", err)
	}
	return authCode, nil
}

// FindByCode retrieves an AuthCode from the database using the code string
func (r *authCodeRepository) FindByCode(code string) (*store.AuthCode, error) {
	r.logger.Info("Searching for AuthCode",
		zap.String("code", code),
	)

	// Initialize a new AuthCode entity
	authCode := new(store.AuthCode)

	// Query the database for the code
	result := r.Db.Where("code = ?", code).First(authCode)

	// Handle errors during the query
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			r.logger.Error("AuthCode not found or invalidated",
				zap.String("code", code),
				zap.Error(result.Error),
				zap.Stack("stacktrace"),
			)
			return nil, fmt.Errorf("AuthCode not found or invalidated: %w", result.Error)
		}
		r.logger.Error("Error finding AuthCode",
			zap.String("code", code),
			zap.Error(result.Error),
			zap.Stack("stacktrace"),
		)
		return nil, fmt.Errorf("error finding AuthCode: %w", result.Error)
	}

	r.logger.Info("Successfully found AuthCode",
		zap.String("code", code),
	)

	return authCode, nil
}
