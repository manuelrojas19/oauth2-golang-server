package handlers

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"html/template"
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
		// Log the error using fmt.Printf for visibility during startup,
		// then panic as the application cannot proceed without the template.
		panic(fmt.Sprintf("FATAL: Error parsing login template: %v", err))
	}

}

// loginHandler handles requests related to user authentication and login.
type loginHandler struct {
	log *zap.Logger // Logger for logging messages within the handler.
}

// NewLoginHandler creates and returns a new instance of loginHandler.
// It takes a *zap.Logger as a dependency for structured logging.
func NewLoginHandler(logger *zap.Logger) LoginHandler {
	return &loginHandler{log: logger}
}

// LoginData holds data that will be passed to the login HTML template.
type LoginData struct {
	GoogleURL string // The URL to initiate Google OAuth authentication.
}

// Login handles the HTTP GET request for the login page.
// It constructs the Google OAuth authentication URL and renders the login template.
func (l loginHandler) Login(writer http.ResponseWriter, request *http.Request) {
	// Extract original authorization parameters from the request URL query.
	originalParams := map[string]string{
		"client_id":     request.URL.Query().Get("client_id"),
		"scope":         request.URL.Query().Get("scope"),
		"redirect_uri":  request.URL.Query().Get("redirect_uri"),
		"response_type": request.URL.Query().Get("response_type"),
		"state":         request.URL.Query().Get("state"),
	}

	// Encode the original parameters into a state string for round-tripping through Google OAuth.
	params, err := encodeState(originalParams)
	if err != nil {
		l.log.Error("Failed to encode state for login page", zap.Error(err))
		http.Error(writer, "Failed to encode state parameters for Google OAuth", http.StatusInternalServerError)
		return
	}

	// Construct the Google OAuth authentication URL.
	googleAuthURL := fmt.Sprintf("%s?client_id=%s&redirect_uri=%s&response_type=code&scope=%s&state=%s",
		configuration.GoogleAuthURL, configuration.GoogleClientID, url.QueryEscape(configuration.GoogleRedirectURL), configuration.Scopes, params)

	// Prepare data to be passed to the login template.
	data := LoginData{
		GoogleURL: googleAuthURL,
	}

	// Execute the login template, rendering the HTML response.
	if err := loginTmpl.Execute(writer, data); err != nil {
		l.log.Error("Error rendering login template", zap.Error(err))
		http.Error(writer, "Error rendering the login page", http.StatusInternalServerError)
	}
}

// encodeState marshals a map of string key-value pairs into a JSON string,
// and then base64 URL-encodes it. This is used to safely pass state information
// through the OAuth flow.
func encodeState(state map[string]string) (string, error) {
	stateJSON, err := json.Marshal(state)
	if err != nil {
		return "", fmt.Errorf("failed to marshal state to JSON: %w", err)
	}
	return base64.URLEncoding.EncodeToString(stateJSON), nil
}
