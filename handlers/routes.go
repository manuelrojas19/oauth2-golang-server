package handlers

import "net/http"

type Routes interface {
	Routes() map[string]http.HandlerFunc
}

type handlerRoutes struct {
	RegisterHandler          RegisterHandler
	TokenHandler             TokenHandler
	AuthorizeHandler         AuthorizeHandler
	RequestConsentHandler    RequestConsentHandler
	JwksHandler              JwksHandler
	AuthorizeCallbackHandler AuthorizeCallbackHandler
	LoginHandler             LoginHandler
	HealthHandler            HealthHandler
}

func NewHandlerRoutes(
	registerHandler RegisterHandler,
	tokenHandler TokenHandler,
	authorizeHandler AuthorizeHandler,
	requestConsentHandler RequestConsentHandler,
	jwksHandler JwksHandler,
	authorizeCallbackHandler AuthorizeCallbackHandler,
	loginHandler LoginHandler,
	healthHandler HealthHandler,
) Routes {
	return &handlerRoutes{RegisterHandler: registerHandler,
		TokenHandler:             tokenHandler,
		AuthorizeHandler:         authorizeHandler,
		RequestConsentHandler:    requestConsentHandler,
		JwksHandler:              jwksHandler,
		AuthorizeCallbackHandler: authorizeCallbackHandler,
		LoginHandler:             loginHandler,
		HealthHandler:            healthHandler,
	}
}

func (h *handlerRoutes) Routes() map[string]http.HandlerFunc {
	return map[string]http.HandlerFunc{
		"/oauth/register":            h.RegisterHandler.Register,
		"/oauth/token":               h.TokenHandler.Token,
		"/oauth/authorize":           h.AuthorizeHandler.Authorize,
		"/oauth/consent":             h.RequestConsentHandler.RequestConsent,
		"/oauth/login":               h.LoginHandler.Login,
		"/.well-known/jwks.json":     h.JwksHandler.Jwks,
		"/google/authorize/callback": h.AuthorizeCallbackHandler.ProcessCallback,
		"/health":                    h.HealthHandler.Health,
	}
}
