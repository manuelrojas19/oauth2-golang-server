package handlers

import (
	"fmt"
	"net/http"
	"net/url"

	"github.com/manuelrojas19/go-oauth2-server/api"
	"github.com/manuelrojas19/go-oauth2-server/errors"
	"github.com/manuelrojas19/go-oauth2-server/oauth"
	"github.com/manuelrojas19/go-oauth2-server/services"
	"go.uber.org/zap"
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
		handleAuthError(w, r, "", "", api.ErrorResponseBody(api.ErrInvalidRequest), a.log)
		return
	}
	a.log.Info("Authorization request decoded", zap.Any("authRequest", authRequest))

	// Validate the authorization request
	if err := authRequest.Validate(); err != nil {
		a.log.Error("Invalid authorization request", zap.Error(err))
		handleAuthError(w, r, authRequest.RedirectUri, authRequest.State, api.ErrorResponseBody(api.ErrInvalidRequest), a.log)
		return
	}
	a.log.Info("Authorization request validated", zap.Any("clientId", authRequest))

	// Create AuthorizeCommand from the request
	command := &services.AuthorizeCommand{
		ClientId:            authRequest.ClientId,
		Scope:               authRequest.Scope,
		RedirectUri:         authRequest.RedirectUri,
		ResponseType:        authRequest.ResponseType,
		State:               authRequest.State,
		CodeChallenge:       authRequest.CodeChallenge,
		CodeChallengeMethod: authRequest.CodeChallengeMethod,
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
		handleAuthError(w, r, authRequest.RedirectUri, authRequest.State, api.ErrorResponseBody(api.ErrUnsupportedResponseType), log)
	default:
		handleAuthError(w, r, authRequest.RedirectUri, authRequest.State, api.ErrorResponseBody(api.ErrServerError), log)
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

func handleAuthError(w http.ResponseWriter, r *http.Request, redirectURI, state string, errorResponse api.ErrorResponse, log *zap.Logger) {
	// Default redirect URI if not provided
	if redirectURI == "" {
		redirectURI = "default/error/page" // Replace with your default error page
	}

	// Construct the error response URL
	errorResponseURL := fmt.Sprintf("%s?error=%s", redirectURI, errorResponse.Error)
	if errorResponse.ErrorDescription != "" {
		errorResponseURL += fmt.Sprintf("&error_description=%s", url.QueryEscape(errorResponse.ErrorDescription))
	}
	if state != "" {
		errorResponseURL += fmt.Sprintf("&state=%s", url.QueryEscape(state))
	}

	// Log the error for debugging purposes
	log.Error("Redirecting with error", zap.String("error_response", errorResponseURL))

	// Redirect the client to the redirect URI with the error
	http.Redirect(w, r, errorResponseURL, http.StatusFound)
}
