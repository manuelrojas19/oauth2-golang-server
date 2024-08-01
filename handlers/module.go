package handlers

import (
	"go.uber.org/fx"
	"net/http"
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
		func() *http.ServeMux { return http.NewServeMux() },
	),
)
