package store

type OauthClientRepository interface {
	Save(client *OauthClient) (*OauthClient, error)
	FindByClientId(clientKey string) (*OauthClient, error)
}

type AccessTokenRepository interface {
	Save(token *AccessToken) (*AccessToken, error)
}

type RefreshTokenRepository interface {
	Save(token *RefreshToken) (*RefreshToken, error)
	FindByToken(token string) (*RefreshToken, error)
	InvalidateRefreshTokensByAccessTokenId(tokenId string) error
}

type UserConsentRepository interface {
	HasUserConsented(userID, clientID, scope string) (bool, error)
	Save(userID, clientID, scope string) (bool, error)
}

type AuthorizationRepository interface {
	Save(authCode *AuthCode) (*AuthCode, error)
	FindByCode(code string) (*AuthCode, error)
}

type UserRepository interface {
	Save(authCode *User) (*User, error)
	FindByUserId(id string) (*User, error)
}
