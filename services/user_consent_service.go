package services

import "github.com/manuelrojas19/go-oauth2-server/store/repositories"

type userConsentService struct {
	consentRepo repositories.UserConsentRepository
}

func NewUserConsentService(consentRepo repositories.UserConsentRepository) UserConsentService {
	return &userConsentService{consentRepo: consentRepo}
}

func (c userConsentService) Save(userID, clientID, scope string) error {
	//TODO implement me
	panic("implement me")
}

func (c userConsentService) HasUserConsented(userID, clientID, scope string) bool {
	//TODO implement me
	panic("implement me")
}
