package repositories

import (
	"errors"
	"fmt"

	"github.com/manuelrojas19/go-oauth2-server/store"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type accessConsentRepository struct {
	Db                    *gorm.DB
	userRepository        UserRepository
	oauthClientRepository OauthClientRepository
	logger                *zap.Logger
}

func NewAccessConsentRepository(db *gorm.DB, userRepo UserRepository, clientRepo OauthClientRepository, logger *zap.Logger) AccessConsentRepository {
	return &accessConsentRepository{
		Db:                    db,
		userRepository:        userRepo,
		oauthClientRepository: clientRepo,
		logger:                logger,
	}
}

func (a *accessConsentRepository) HasUserConsented(userID, clientID, scopeID string) (bool, error) {
	var consent store.AccessConsent
	a.logger.Debug("Querying for user consent", zap.String("userId", userID), zap.String("clientId", clientID), zap.String("scopeId", scopeID))
	err := a.Db.Where("user_id = ? AND client_id = ? AND scope_id = ?", userID, clientID, scopeID).First(&consent).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		a.logger.Debug("No consent record found", zap.String("userId", userID), zap.String("clientId", clientID), zap.String("scopeId", scopeID))
		return false, nil
	} else if err != nil {
		a.logger.Error("Error checking user consent from database", zap.String("userId", userID), zap.String("clientId", clientID), zap.String("scopeId", scopeID), zap.Error(err))
		return false, fmt.Errorf("failed to check user consent: %w", err)
	}
	a.logger.Debug("Consent record found", zap.String("userId", userID), zap.String("clientId", clientID), zap.String("scopeId", scopeID), zap.Bool("consented", consent.Consented))
	return consent.Consented, nil
}

func (a *accessConsentRepository) Save(userId, clientId, scopeId string) (*store.AccessConsent, error) {
	var consent *store.AccessConsent

	a.logger.Info("Attempting to save access consent in transaction", zap.String("userId", userId), zap.String("clientId", clientId), zap.String("scopeId", scopeId))
	err := a.Db.Transaction(func(tx *gorm.DB) error {
		// Check if the user exists
		a.logger.Debug("Checking if user exists", zap.String("userId", userId))
		_, err := a.userRepository.FindByUserId(userId)
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				a.logger.Error("User not found for consent", zap.String("userId", userId))
				return fmt.Errorf("user not found: %w", err)
			}
			a.logger.Error("Error finding user for consent", zap.String("userId", userId), zap.Error(err))
			return fmt.Errorf("failed to find user for consent: %w", err)
		}
		a.logger.Debug("User exists", zap.String("userId", userId))

		// Check if the client exists
		a.logger.Debug("Checking if client exists", zap.String("clientId", clientId))
		_, err = a.oauthClientRepository.FindByClientId(clientId)
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				a.logger.Error("Client not found for consent", zap.String("clientId", clientId))
				return fmt.Errorf("client not found: %w", err)
			}
			a.logger.Error("Error finding client for consent", zap.String("clientId", clientId), zap.Error(err))
			return fmt.Errorf("failed to find client for consent: %w", err)
		}
		a.logger.Debug("Client exists", zap.String("clientId", clientId))

		// Create the AccessConsent record using builder
		consent = store.NewUserConsentBuilder().
			WithUserId(userId).
			WithClientId(clientId).
			WithScopeId(scopeId).
			WithConsented(true).
			Build()
		a.logger.Debug("AccessConsent entity built", zap.Any("consentEntity", consent))

		// Save the consent to the database
		if err := tx.Create(consent).Error; err != nil {
			a.logger.Error("Error saving consent to database", zap.String("userId", userId), zap.String("clientId", clientId), zap.String("scopeId", scopeId), zap.Error(err))
			return fmt.Errorf("failed to save access consent: %w", err)
		}
		a.logger.Info("Consent saved successfully in transaction", zap.String("consentId", consent.Id))

		return nil
	})

	if err != nil {
		a.logger.Error("Transaction failed for saving access consent", zap.Error(err))
		return nil, fmt.Errorf("access consent transaction failed: %w", err)
	}

	a.logger.Info("Access consent saved successfully", zap.String("consentId", consent.Id))
	return consent, nil
}
