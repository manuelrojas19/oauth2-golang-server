package handlers

import (
	"fmt"
	"github.com/manuelrojas19/go-oauth2-server/services"
	"github.com/manuelrojas19/go-oauth2-server/services/commands"
	"net/http"
	"net/url"
)

type AuthorizeHandler struct {
	authorizationService services.AuthorizationService
}

func NewAuthorizeHandler(authorizationService services.AuthorizationService) Handler {
	return &AuthorizeHandler{
		authorizationService: authorizationService,
	}
}

func (a AuthorizeHandler) Handler(w http.ResponseWriter, r *http.Request) {
	clientId := r.URL.Query().Get("client_id")
	scope := r.URL.Query().Get("scope")
	redirectUri := r.URL.Query().Get("redirect_uri")
	responseType := r.URL.Query().Get("response_type")

	// Validate response_type
	if responseType != "code" {
		http.Error(w, "unsupported_response_type", http.StatusBadRequest)
		return
	}

	// URL-encode parameters
	encodedClientId := url.QueryEscape(clientId)
	encodedScope := url.QueryEscape(scope)
	encodedRedirectUri := url.QueryEscape(redirectUri)
	encodedResponseType := url.QueryEscape(responseType)

	// Prepare the query parameters string
	queryParams := fmt.Sprintf("client_id=%s&scope=%s&redirect_uri=%s&response_type=%s",
		encodedClientId, encodedScope, encodedRedirectUri, encodedResponseType)

	// Authorization request
	command := &commands.Authorization{
		ClientId:     clientId,
		Scope:        scope,
		RedirectUri:  redirectUri,
		ResponseType: responseType,
	}

	authCode, err := a.authorizationService.Authorize(command)

	if err != nil {
		switch err.Error() {
		case services.ErrUserNotAuthenticated:
			loginURL := fmt.Sprintf("/login?%s", queryParams)
			http.Redirect(w, r, loginURL, http.StatusSeeOther)
			return
		case services.ErrConsentRequired:
			consentURL := fmt.Sprintf("/consent?%s", queryParams)
			http.Redirect(w, r, consentURL, http.StatusSeeOther)
			return
		default:
			http.Error(w, "server_error", http.StatusInternalServerError)
			return
		}
	}

	// Redirect to the redirect_uri with authorization code
	redirectURL := fmt.Sprintf("%s?code=%s", encodedRedirectUri, authCode)
	http.Redirect(w, r, redirectURL, http.StatusSeeOther)
}
