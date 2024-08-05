package services

import (
	"go.uber.org/fx"
)

var Module = fx.Options(
	fx.Provide(
		NewUserConsentService,
		NewOauthClientService,
		NewTokenService,
		NewAuthorizationService,
		NewWellKnownService,
		NewScopeService,
	),
)
