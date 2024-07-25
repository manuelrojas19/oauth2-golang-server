package handlers

import (
	"fmt"
	"github.com/manuelrojas19/go-oauth2-server/services"
	"github.com/manuelrojas19/go-oauth2-server/services/commands"
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
	if responseType != "code" {
		http.Error(w, "Unsupported response type", http.StatusBadRequest)
		return
	}

	// Call the authorization function
	command := &commands.Authorization{
		ClientId:     clientId,
		Scope:        scope,
		RedirectUri:  redirectUri,
		ResponseType: responseType,
	}

	authCode, err := a.authorizationService.Authorize(command)

	if err != nil {
		// Handle redirections
		switch err.Error() {
		case services.ErrNotAuthenticated:
			http.Redirect(w, r, "/login", http.StatusSeeOther)
			return
		case services.ErrConsentRequired:
			consentURL := fmt.Sprintf("/consent?client_id=%s&scope=%s", clientId, scope)
			http.Redirect(w, r, consentURL, http.StatusSeeOther)
			return
		default:
			http.Error(w, "Authorization failed", http.StatusInternalServerError)
			return
		}
	}

	// If no errors, proceed with the authorization code flow
	redirectURI := r.URL.Query().Get("redirect_uri")
	http.Redirect(w, r, fmt.Sprintf("%s?code=%s", redirectURI, authCode), http.StatusSeeOther)
}
