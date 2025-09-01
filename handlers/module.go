package handlers

import (
	"go.uber.org/fx"
)

var Module = fx.Options(
	fx.Provide(
		NewRegisterHandler,
		NewTokenHandler,
		NewJwksHandler,
		NewAuthorizeHandler,
		NewRequestConsentHandler,
		NewAuthorizeCallbackHandler,
		NewLoginHandler,
		NewHealthHandler,
		NewUserinfoHandler,
		NewLogoutHandler,
		NewIntrospectionHandler,
		NewRevocationHandler,
		NewAcceptConsentHandler,
	),
)
