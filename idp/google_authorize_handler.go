package idp

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/manuelrojas19/go-oauth2-server/configuration"
	"github.com/manuelrojas19/go-oauth2-server/handlers"
	"net/http"
	"net/url"
)

type googleAuthorizeHandler struct {
}

func NewGoogleLoginHandler() handlers.Handler {
	return &googleAuthorizeHandler{}
}

func (g googleAuthorizeHandler) Handler(writer http.ResponseWriter, request *http.Request) {
	originalParams := map[string]string{
		"client_id":     request.URL.Query().Get("client_id"),
		"scope":         request.URL.Query().Get("scope"),
		"redirect_uri":  request.URL.Query().Get("redirect_uri"),
		"response_type": request.URL.Query().Get("response_type"),
	}

	encodedParams, err := encodeState(originalParams)
	if err != nil {
		http.Error(writer, "Failed to encode state", http.StatusInternalServerError)
		return
	}

	authURL := fmt.Sprintf("%s?client_id=%s&redirect_uri=%s&response_type=code&scope=%s&state=%s",
		configuration.GoogleAuthURL, configuration.GoogleClientID, url.QueryEscape(configuration.GoogleRedirectURL), configuration.Scopes, encodedParams)
	http.Redirect(writer, request, authURL, http.StatusTemporaryRedirect)
}

func encodeState(state map[string]string) (string, error) {
	stateJSON, err := json.Marshal(state)
	if err != nil {
		return "", err
	}
	return base64.URLEncoding.EncodeToString(stateJSON), nil
}
