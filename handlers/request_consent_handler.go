package handlers

import (
	"fmt"
	"html/template"
	"net/http"
	"net/url"

	"go.uber.org/zap"
)

var requestTmpl *template.Template

func init() {
	var err error
	requestTmpl, err = template.ParseFiles("templates/authorize.html")
	if err != nil {
		panic(fmt.Sprintf("Error parsing template: %v", err))
	}
}

type requestConsentHandler struct {
	logger *zap.Logger
}

func NewRequestConsentHandler(logger *zap.Logger) RequestConsentHandler {
	return &requestConsentHandler{logger: logger}
}

type ConsentPageData struct {
	ClientId       string
	Scope          string
	RedirectUri    string
	ResponseType   string
	ConsentPageURL string
}

func (h *requestConsentHandler) RequestConsent(w http.ResponseWriter, r *http.Request) {
	h.logger.Info("Received request for consent page",
		zap.String("client_id", r.URL.Query().Get("client_id")),
		zap.String("scope", r.URL.Query().Get("scope")))

	clientId := r.URL.Query().Get("client_id")
	scope := r.URL.Query().Get("scope")
	redirectUri := r.URL.Query().Get("redirect_uri")
	responseType := r.URL.Query().Get("response_type")

	// URL-encode parameters to include in the consent page link
	encodedClientId := url.QueryEscape(clientId)
	encodedScope := url.QueryEscape(scope)
	encodedRedirectUri := url.QueryEscape(redirectUri)
	encodedResponseType := url.QueryEscape(responseType)

	// Construct the consent page URL with query parameters
	consentPageURL := fmt.Sprintf("/consent/?client_id=%s&scope=%s&redirect_uri=%s&response_type=%s",
		encodedClientId, encodedScope, encodedRedirectUri, encodedResponseType)

	data := ConsentPageData{
		ClientId:       clientId,
		Scope:          scope,
		RedirectUri:    redirectUri,
		ResponseType:   responseType,
		ConsentPageURL: consentPageURL,
	}

	// Redirect the user to the consent page
	if err := requestTmpl.Execute(w, data); err != nil {
		h.logger.Error("Error rendering consent template", zap.Error(err))
		http.Error(w, "Error rendering template", http.StatusInternalServerError)
	}
}
