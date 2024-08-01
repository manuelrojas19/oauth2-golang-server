package handlers

import (
	"fmt"
	"net/http"
	"net/url"
)

type acceptConsentHandler struct{}

func NewAcceptConsentHandler() AcceptConsentHandler {
	return &acceptConsentHandler{}
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
		// Handle consent denial (e.g., show an error message)
		http.Error(w, "Consent denied", http.StatusForbidden)
	}
}
