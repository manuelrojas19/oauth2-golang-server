package handlers

import (
	"fmt"
	"net/http"
	"net/url"

	"github.com/manuelrojas19/go-oauth2-server/api"
	"go.uber.org/zap"
)

type acceptConsentHandler struct {
	log *zap.Logger
}

func NewAcceptConsentHandler(log *zap.Logger) AcceptConsentHandler {
	return &acceptConsentHandler{
		log: log,
	}
}

func (h *acceptConsentHandler) AcceptConsent(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	clientId := r.FormValue("client_id")
	scope := r.FormValue("scope")
	redirectUri := r.FormValue("redirect_uri")
	responseType := r.FormValue("response_type")
	consent := r.FormValue("consent")

	// URL-encode parameters for redirect
	encodedClientId := url.QueryEscape(clientId)
	encodedScope := url.QueryEscape(scope)
	encodedRedirectUri := url.QueryEscape(redirectUri)
	encodedResponseType := url.QueryEscape(responseType)

	if consent == "approve" {
		// Redirect back to the original authorization endpoint with original parameters
		redirectURL := fmt.Sprintf("/oauth/authorize?client_id=%s&scope=%s&redirect_uri=%s&response_type=%s",
			encodedClientId, encodedScope, encodedRedirectUri, encodedResponseType)
		http.Redirect(w, r, redirectURL, http.StatusSeeOther)
	} else {
		// Handle consent denial (e.g., redirect with an access_denied error)
		errorResponse := api.ErrorResponse{
			Error:            api.ErrAccessDenied.Error(),
			ErrorDescription: "Resource owner denied the request",
		}
		h.log.Error("Consent denied, redirecting with error", zap.Any("error_response", errorResponse))
		handleAuthError(w, r, redirectUri, "", api.ErrorResponseBody(api.ErrAccessDenied), h.log)
	}
}
