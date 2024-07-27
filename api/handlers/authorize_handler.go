package handlers

import (
	"fmt"
	"github.com/manuelrojas19/go-oauth2-server/services"
	"github.com/manuelrojas19/go-oauth2-server/services/commands"
	"log"
	"net/http"
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

	// Authorize request
	command := &commands.Authorize{
		ClientId:     clientId,
		Scope:        scope,
		RedirectUri:  redirectUri,
		ResponseType: responseType,
	}

	// SessionId
	cookie, err := r.Cookie("session_id")
	if err == nil && cookie != nil {
		log.Println("Session cookie found")
		command.SessionId = cookie.Value
	}

	authCode, err := a.authorizationService.Authorize(command)

	// Prepare the query parameters string
	queryParams := fmt.Sprintf("client_id=%s&scope=%s&redirect_uri=%s&response_type=%s",
		clientId, scope, redirectUri, responseType)

	if err != nil {
		switch err.Error() {
		case services.ErrUserNotAuthenticated:
			loginURL := fmt.Sprintf("/google/authorize?%s", queryParams)
			http.Redirect(w, r, loginURL, http.StatusSeeOther)
			return
		case services.ErrConsentRequired:
			consentURL := fmt.Sprintf("/oauth/consent?%s", queryParams)
			http.Redirect(w, r, consentURL, http.StatusSeeOther)
			return
		default:
			http.Error(w, "server_error", http.StatusInternalServerError)
			return
		}
	}

	// Redirect to the redirect_uri with authorization code
	redirectURL := fmt.Sprintf("%s?code=%s", redirectUri, authCode.Code)
	log.Printf("Redirecting to: %s", redirectURL)
	http.Redirect(w, r, redirectURL, http.StatusSeeOther)
}
