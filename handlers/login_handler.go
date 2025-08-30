package handlers

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"net/url"

	"github.com/manuelrojas19/go-oauth2-server/configuration"
	"go.uber.org/zap"
)

var loginTmpl *template.Template

func init() {
	var err error
	loginTmpl, err = template.ParseFiles("templates/login.html")
	if err != nil {
		log.Fatalf("Error parsing template: %v", err)
	}

}

type loginHandler struct {
	log *zap.Logger
}

func NewLoginHandler(logger *zap.Logger) LoginHandler {
	return &loginHandler{log: logger}
}

type LoginData struct {
	GoogleURL string
}

func (l loginHandler) Login(writer http.ResponseWriter, request *http.Request) {
	// Redirect the user to the consent page
	originalParams := map[string]string{
		"client_id":     request.URL.Query().Get("client_id"),
		"scope":         request.URL.Query().Get("scope"),
		"redirect_uri":  request.URL.Query().Get("redirect_uri"),
		"response_type": request.URL.Query().Get("response_type"),
	}

	params, err := encodeState(originalParams)
	if err != nil {
		http.Error(writer, "Failed to encode state", http.StatusInternalServerError)
		return
	}

	googleAuthURL := fmt.Sprintf("%s?client_id=%s&redirect_uri=%s&response_type=code&scope=%s&state=%s",
		configuration.GoogleAuthURL, configuration.GoogleClientID, url.QueryEscape(configuration.GoogleRedirectURL), configuration.Scopes, params)

	data := LoginData{
		GoogleURL: googleAuthURL,
	}

	if err := loginTmpl.Execute(writer, data); err != nil {
		http.Error(writer, "Error rendering template", http.StatusInternalServerError)
	}
}

func encodeState(state map[string]string) (string, error) {
	stateJSON, err := json.Marshal(state)
	if err != nil {
		return "", err
	}
	return base64.URLEncoding.EncodeToString(stateJSON), nil
}
