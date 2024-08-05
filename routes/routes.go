package routes

import (
	"github.com/manuelrojas19/go-oauth2-server/handlers"
	"net/http"
)

func Routes(registerHandler handlers.RegisterHandler,
	tokenHandler handlers.TokenHandler,
	authorizeHandler handlers.AuthorizeHandler,
	requestConsentHandler handlers.RequestConsentHandler,
	jwksHandler handlers.JwksHandler,
	authorizeCallbackHandler handlers.AuthorizeCallbackHandler,
	loginHandler handlers.LoginHandler,
	healthHandler handlers.HealthHandler,
	scopesHandler handlers.ScopeHandler,
) map[string]http.HandlerFunc {
	return map[string]http.HandlerFunc{
		"/oauth/scope":               scopesHandler.CreateScope,
		"/oauth/register":            registerHandler.Register,
		"/oauth/token":               tokenHandler.Token,
		"/oauth/authorize":           authorizeHandler.Authorize,
		"/oauth/consent":             requestConsentHandler.RequestConsent,
		"/oauth/login":               loginHandler.Login,
		"/.well-known/jwks.json":     jwksHandler.Jwks,
		"/google/authorize/callback": authorizeCallbackHandler.ProcessCallback,
		"/health":                    healthHandler.Health,
	}
}
