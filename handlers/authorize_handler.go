package handlers

import (
	"fmt"
	"github.com/manuelrojas19/go-oauth2-server/api"
	"github.com/manuelrojas19/go-oauth2-server/errors"
	"github.com/manuelrojas19/go-oauth2-server/oauth"
	"github.com/manuelrojas19/go-oauth2-server/services"
	"go.uber.org/zap"
	"net/http"
	"net/url"
)

type authorizeHandler struct {
	authorizationService services.AuthorizationService
	log                  *zap.Logger
}

func NewAuthorizeHandler(authorizationService services.AuthorizationService, logger *zap.Logger) AuthorizeHandler {
	return &authorizeHandler{
		authorizationService: authorizationService,
		log:                  logger,
	}
}

func (a authorizeHandler) Authorize(w http.ResponseWriter, r *http.Request) {
	a.log.Info("Entered Authorize handler", zap.String("method", r.Method), zap.String("url", r.URL.String()))

	// Decode the authorization request
	authRequest, err := api.DecodeAuthorizeRequest(r)
	if err != nil {
		a.log.Error("Failed to decode authorization request", zap.Error(err))
		handleAuthError(w, r, "", "", "invalid_request", err.Error(), a.log)
		return
	}
	a.log.Info("Authorization request decoded", zap.Any("authRequest", authRequest))

	// Validate the authorization request
	if err := authRequest.Validate(); err != nil {
		a.log.Error("Invalid authorization request", zap.Error(err))
		handleAuthError(w, r, authRequest.RedirectUri, authRequest.State, "invalid_request", err.Error(), a.log)
		return
	}
	a.log.Info("Authorization request validated", zap.Any("clientId", authRequest))

	// Create AuthorizeCommand from the request
	command := &services.AuthorizeCommand{
		ClientId:     authRequest.ClientId,
		Scope:        authRequest.Scope,
		RedirectUri:  authRequest.RedirectUri,
		ResponseType: authRequest.ResponseType,
		State:        authRequest.State,
	}
	a.log.Info("AuthorizeCommand created", zap.Any("command", command))

	// Retrieve SessionId from the cookie
	cookie, err := r.Cookie("session_id")
	if err == nil && cookie != nil {
		a.log.Info("Session cookie found", zap.String("session_id", cookie.Value))
		command.SessionId = cookie.Value
	} else {
		a.log.Warn("Session cookie not found or error retrieving cookie", zap.Error(err))
	}

	// Authorize the request
	authCode, err := a.authorizationService.Authorize(command)
	if err != nil {
		a.log.Error("Authorization service error", zap.Error(err), zap.Stack("stacktrace"))
		handleAuthorizationError(err, w, r, authRequest, command, a.log)
		return
	}
	a.log.Info("Authorization successful", zap.String("authCode", authCode.Code))

	// Build the redirect URL
	redirectURL := getRedirectURL(authRequest, authCode)
	a.log.Info("Redirect URL built", zap.String("getRedirectURL", redirectURL))

	// Redirect to the redirect_uri with the authorization code
	http.Redirect(w, r, redirectURL, http.StatusSeeOther)
}

func handleAuthorizationError(err error, w http.ResponseWriter, r *http.Request, authRequest *api.AuthorizeRequest, command *services.AuthorizeCommand, log *zap.Logger) {
	queryParams := fmt.Sprintf("client_id=%s&scope=%s&redirect_uri=%s&response_type=%s",
		authRequest.ClientId, authRequest.Scope, authRequest.RedirectUri, string(authRequest.ResponseType))

	switch err.Error() {
	case errors.ErrUserNotAuthenticated:
		loginURL := fmt.Sprintf("/oauth/login?%s", queryParams)
		log.Warn("User not authenticated, redirecting to login", zap.String("loginURL", loginURL))
		http.Redirect(w, r, loginURL, http.StatusSeeOther)
	case errors.ErrConsentRequired:
		consentURL := fmt.Sprintf("/oauth/consent?%s", queryParams)
		log.Warn("User consent required, redirecting to consent", zap.String("consentURL", consentURL))
		http.Redirect(w, r, consentURL, http.StatusSeeOther)
	case errors.ErrUnsupportedResponseType:
		handleAuthError(w, r, authRequest.RedirectUri, authRequest.State, "unsupported_response_type", err.Error(), log)
	default:
		handleAuthError(w, r, authRequest.RedirectUri, authRequest.State, "server_error",
			"the authorization server encountered an unexpected condition that prevented it from fulfilling the request", log)
	}
}

func getRedirectURL(authRequest *api.AuthorizeRequest, authCode *oauth.AuthCode) string {
	redirectURL := authRequest.RedirectUri
	if authRequest.State != "" {
		redirectURL = fmt.Sprintf("%s?code=%s&state=%s", redirectURL, authCode.Code, authRequest.State)
	} else {
		redirectURL = fmt.Sprintf("%s?code=%s", redirectURL, authCode.Code)
	}
	return redirectURL
}

func handleAuthError(w http.ResponseWriter, r *http.Request, redirectURI, state, errorCode, errorDescription string, log *zap.Logger) {
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
	log.Error("Redirecting with error", zap.String("error_response", errorResponse))

	// Redirect the client to the redirect URI with the error
	http.Redirect(w, r, errorResponse, http.StatusFound)
}
