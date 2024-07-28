package services

import (
	"github.com/manuelrojas19/go-oauth2-server/store"
)

type userConsentService struct {
	consentRepo store.UserConsentRepository
}

func NewUserConsentService(consentRepo store.UserConsentRepository) UserConsentService {
	return &userConsentService{consentRepo: consentRepo}
}

func (c userConsentService) Save(userID, clientID, scope string) error {
	//TODO implement me
	panic("implement me")
}

func (c userConsentService) HasUserConsented(userID, clientID, scope string) bool {
	//TODO implement me
	return false
}
