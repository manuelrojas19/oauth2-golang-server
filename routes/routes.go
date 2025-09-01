package routes

import (
	"net/http"

	"github.com/manuelrojas19/go-oauth2-server/handlers"
)

func Routes(registerHandler handlers.RegisterHandler,
	tokenHandler handlers.TokenHandler,
	authorizeHandler handlers.AuthorizeHandler,
	requestConsentHandler handlers.RequestConsentHandler,
	jwksHandler handlers.JwksHandler,
	authorizeCallbackHandler handlers.AuthorizeCallbackHandler,
	loginHandler handlers.LoginHandler,
	healthHandler handlers.HealthHandler,
	userinfoHandler handlers.UserinfoHandler,
	logoutHandler handlers.LogoutHandler,
	introspectionHandler handlers.IntrospectionHandler,
	revocationHandler handlers.RevocationHandler,
) map[string]http.HandlerFunc {
	return map[string]http.HandlerFunc{
		"/oauth/register":            registerHandler.Register,
		"/oauth/token":               tokenHandler.Token,
		"/oauth/authorize":           authorizeHandler.Authorize,
		"/oauth/consent":             requestConsentHandler.RequestConsent,
		"/oauth/login":               loginHandler.Login,
		"/.well-known/jwks.json":     jwksHandler.Jwks,
		"/oauth/userinfo":            userinfoHandler.Userinfo,
		"/oauth/logout":              logoutHandler.Logout,
		"/oauth/introspect":          introspectionHandler.Introspect,
		"/oauth/revoke":              revocationHandler.Revoke,
		"/google/authorize/callback": authorizeCallbackHandler.ProcessCallback,
		"/health":                    healthHandler.Health,
	}
}
