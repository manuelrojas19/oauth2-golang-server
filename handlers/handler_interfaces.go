package handlers

import "net/http"

type AcceptConsentHandler interface {
	AcceptConsent(http.ResponseWriter, *http.Request)
}

type AuthorizeHandler interface {
	Authorize(http.ResponseWriter, *http.Request)
}

type AuthorizeCallbackHandler interface {
	ProcessCallback(http.ResponseWriter, *http.Request)
}

type JwksHandler interface {
	Jwks(http.ResponseWriter, *http.Request)
}

type LoginHandler interface {
	Login(http.ResponseWriter, *http.Request)
}

type RegisterHandler interface {
	Register(http.ResponseWriter, *http.Request)
}

type RequestConsentHandler interface {
	RequestConsent(http.ResponseWriter, *http.Request)
}

type TokenHandler interface {
	Token(http.ResponseWriter, *http.Request)
}

type ScopeHandler interface {
	CreateScope(http.ResponseWriter, *http.Request)
}
