package repositories

import (
	"go.uber.org/fx"
)

var Module = fx.Options(
	fx.Provide(
		NewOauthClientRepository,
		NewAccessTokenRepository,
		NewRefreshTokenRepository,
		NewUserConsentRepository,
		NewAuthCodeRepository,
		NewUserRepository,
	),
)
