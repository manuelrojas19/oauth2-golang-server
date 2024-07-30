package handlers

import (
	"fmt"
	"github.com/manuelrojas19/go-oauth2-server/api"
	"github.com/manuelrojas19/go-oauth2-server/services"
	"log"
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
	// Decode the authorization request
	authRequest, err := api.DecodeAuthorizeRequest(r)

	if err != nil {
		log.Printf("Failed to decode authorization request: %v", err)
		handleAuthError(w, r, "", "", "invalid_request", err.Error())
		return
	}

	// Validate the authorization request
	if err := authRequest.Validate(); err != nil {
		log.Printf("Invalid authorization request: %v", err)
		handleAuthError(w, r, authRequest.RedirectUri, authRequest.State, "invalid_request", err.Error())
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
		case services.ErrUnsupportedResponseType:
			handleAuthError(w, r, authRequest.RedirectUri, authRequest.State, "unsupported_response_type", err.Error())
			return
		default:
			handleAuthError(w, r, authRequest.RedirectUri, authRequest.State, "server_error",
				"The authorization server encountered an unexpected condition that prevented it from fulfilling the request.")
			return
		}
	}

	// Build the redirect URL
	redirectURL := authRequest.RedirectUri

	// Check if state is not empty
	if authRequest.State != "" {
		redirectURL = fmt.Sprintf("%s?code=%s&state=%s", redirectURL, authCode.Code, authRequest.State)
	} else {
		redirectURL = fmt.Sprintf("%s?code=%s", redirectURL, authCode.Code)
	}

	// Redirect to the redirect_uri with the authorization code
	log.Printf("Redirecting to: %s", redirectURL)
	http.Redirect(w, r, redirectURL, http.StatusSeeOther)
}

func handleAuthError(w http.ResponseWriter, r *http.Request, redirectURI, state, errorCode, errorDescription string) {
	// Default redirect URI if not provided
	if redirectURI == "" {
		redirectURI = "default/error/page" // Replace with your default error page
	}

	// Construct the error response URL
	errorResponse := fmt.Sprintf("%s?error=%s", redirectURI, errorCode)

	if errorDescription != "" {
		errorResponse += fmt.Sprintf("&error_description=%s", url.QueryEscape(errorDescription))
	}

	if state != "" {
		errorResponse += fmt.Sprintf("&state=%s", url.QueryEscape(state))
	}

	// Log the error for debugging purposes
	log.Printf("Redirecting with error: %s", errorResponse)

	// Redirect the client to the redirect URI with the error
	http.Redirect(w, r, errorResponse, http.StatusFound)
}
