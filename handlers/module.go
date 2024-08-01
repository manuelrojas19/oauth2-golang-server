package handlers

import (
	"go.uber.org/fx"
)

var Module = fx.Options(
	fx.Provide(
		NewAcceptConsentHandler,
		NewRegisterHandler,
		NewTokenHandler,
		NewJwksHandler,
		NewAuthorizeHandler,
		NewRequestConsentHandler,
		NewAuthorizeCallbackHandler,
		NewLoginHandler,
		NewHandlerRoutes,
		NewHealthHandler,
	),
)
