package handlers

import (
	"fmt"
	"github.com/manuelrojas19/go-oauth2-server/api"
	"github.com/manuelrojas19/go-oauth2-server/services"
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
	// Decode the authorization request
	authRequest, err := api.DecodeAuthorizeRequest(r)

	if err != nil {
		log.Printf("Failed to decode authorization request: %v", err)
		http.Error(w, "invalid_request", http.StatusBadRequest)
		return
	}

	// Validate the authorization request
	if err := authRequest.Validate(); err != nil {
		log.Printf("Invalid authorization request: %v", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Create AuthorizeCommand from the request
	command := &services.AuthorizeCommand{
		ClientId:     authRequest.ClientId,
		Scope:        authRequest.Scope,
		RedirectUri:  authRequest.RedirectUri,
		ResponseType: authRequest.ResponseType,
		State:        authRequest.State,
	}

	// Retrieve SessionId from the cookie
	cookie, err := r.Cookie("session_id")
	if err == nil && cookie != nil {
		log.Println("Session cookie found")
		command.SessionId = cookie.Value
	}

	authCode, err := a.authorizationService.Authorize(command)

	// Prepare the query parameters string
	queryParams := fmt.Sprintf("client_id=%s&scope=%s&redirect_uri=%s&response_type=%s",
		authRequest.ClientId, authRequest.Scope, authRequest.RedirectUri, string(authRequest.ResponseType))

	if err != nil {
		switch err.Error() {
		case services.ErrUserNotAuthenticated:
			loginURL := fmt.Sprintf("/oauth/login?%s", queryParams)
			http.Redirect(w, r, loginURL, http.StatusSeeOther)
			return
		case services.ErrConsentRequired:
			consentURL := fmt.Sprintf("/oauth/consent?%s", queryParams)
			http.Redirect(w, r, consentURL, http.StatusSeeOther)
			return
		default:
			http.Error(w, fmt.Sprintf("server error: %s", err), http.StatusInternalServerError)
			return
		}
	}

	// Redirect to the redirect_uri with the authorization code
	redirectURL := fmt.Sprintf("%s?code=%s&state=%s", authRequest.RedirectUri, authCode.Code, authRequest.State)
	log.Printf("Redirecting to: %s", redirectURL)
	http.Redirect(w, r, redirectURL, http.StatusSeeOther)
}
