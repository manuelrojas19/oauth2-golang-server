package repositories

import (
	"errors"
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
	err := a.Db.Where("user_id = ? AND client_id = ? AND scope_id = ?", userID, clientID, scopeID).First(&consent).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return false, nil
	} else if err != nil {
		return false, err
	}
	return consent.Consented, nil
}

func (a *accessConsentRepository) Save(userId, clientId, scopeId string) (*store.AccessConsent, error) {
	var consent *store.AccessConsent

	err := a.Db.Transaction(func(tx *gorm.DB) error {
		// Check if the user exists
		_, err := a.userRepository.FindByUserId(userId)
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				a.logger.Error("User not found", zap.String("userId", userId))
				return errors.New("user not found")
			}
			a.logger.Error("Error finding user", zap.Error(err))
			return err
		}
		a.logger.Info("User found", zap.String("userId", userId))

		// Check if the client exists
		_, err = a.oauthClientRepository.FindByClientId(clientId)
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				a.logger.Error("Client not found", zap.String("clientId", clientId))
				return errors.New("client not found")
			}
			a.logger.Error("Error finding client", zap.Error(err))
			return err
		}
		a.logger.Info("Client found", zap.String("clientId", clientId))

		// Create the AccessConsent record using builder
		consent = store.NewUserConsentBuilder().
			WithUserId(userId).
			WithClientId(clientId).
			WithScopeId(scopeId).
			WithConsented(true).
			Build()

		// Save the consent to the database
		if err := tx.Create(consent).Error; err != nil {
			a.logger.Error("Error saving consent", zap.Error(err))
			return err
		}
		a.logger.Info("Consent saved successfully", zap.String("userId", userId), zap.String("clientId", clientId), zap.String("scopeId", scopeId))

		return nil
	})

	if err != nil {
		return nil, err
	}

	return consent, nil
}
