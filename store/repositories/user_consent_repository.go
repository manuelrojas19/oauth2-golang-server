package repositories

import "gorm.io/gorm"

type userConsentRepository struct {
	Db *gorm.DB
}

func NewUserConsentRepository(db *gorm.DB) UserConsentRepository {
	return &userConsentRepository{
		Db: db,
	}
}

func (userConsentRepository) HasUserConsented(userID, clientID, scope string) (bool, error) {
	//TODO implement me
	panic("implement me")
}

func (userConsentRepository) Save(userID, clientID, scope string) (bool, error) {
	//TODO implement me
	panic("implement me")
}
