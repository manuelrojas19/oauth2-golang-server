package handlers

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"net/url"
)

var requestTmpl *template.Template

func init() {
	var err error
	requestTmpl, err = template.ParseFiles("templates/authorize.html")
	if err != nil {
		log.Fatalf("Error parsing template: %v", err)
	}
}

type RequestConsentHandler struct {
}

func NewRequestConsentHandler() Handler {
	return &RequestConsentHandler{}
}

type ConsentPageData struct {
	ClientId       string
	Scope          string
	RedirectUri    string
	ResponseType   string
	ConsentPageURL string
}

func (h *RequestConsentHandler) Handler(w http.ResponseWriter, r *http.Request) {
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
		http.Error(w, "Error rendering template", http.StatusInternalServerError)
	}
}
