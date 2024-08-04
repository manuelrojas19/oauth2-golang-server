package repositories

import (
	"go.uber.org/fx"
)

var Module = fx.Options(
	fx.Provide(
		NewOauthClientRepository,
		NewAccessTokenRepository,
		NewRefreshTokenRepository,
		NewAccessConsentRepository,
		NewScopeRepository,
		NewAuthCodeRepository,
		NewUserRepository,
	),
)
