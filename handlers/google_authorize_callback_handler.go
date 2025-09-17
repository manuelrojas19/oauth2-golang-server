package handlers

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"time"

	"github.com/manuelrojas19/go-oauth2-server/configuration"
	"github.com/manuelrojas19/go-oauth2-server/services"
	"github.com/manuelrojas19/go-oauth2-server/store"
	"github.com/manuelrojas19/go-oauth2-server/store/repositories"

	"go.uber.org/zap"
)

// TokenResponse represents the structure of the token response from Google.
type TokenResponse struct {
	AccessToken  string `json:"access_token"`
	TokenType    string `json:"token_type"`
	RefreshToken string `json:"refresh_token"`
	Expiry       int64  `json:"expires_in"`
	IDToken      string `json:"id_token"`
}

// UserInfo represents the structure of the user information received from Google.
type UserInfo struct {
	ID    string `json:"sub"`
	Email string `json:"email"`
	Name  string `json:"name"`
}

// authorizeCallbackHandler handles the callback from the external OAuth2 provider (e.g., Google) after user authentication.
type authorizeCallbackHandler struct {
	userSessionService services.SessionService
	userRepository     repositories.UserRepository
	logger             *zap.Logger
}

// NewAuthorizeCallbackHandler creates and returns a new instance of authorizeCallbackHandler.
func NewAuthorizeCallbackHandler(
	userSessionService services.SessionService,
	userRepository repositories.UserRepository,
	logger *zap.Logger,
) AuthorizeCallbackHandler {
	return &authorizeCallbackHandler{
		userSessionService: userSessionService,
		userRepository:     userRepository,
		logger:             logger,
	}
}

// ProcessCallback handles the HTTP GET request to the OAuth2 redirect URI.
// It exchanges the authorization code for tokens, retrieves user information,
// creates a user session, sets a session cookie, and redirects the user.
func (g authorizeCallbackHandler) ProcessCallback(writer http.ResponseWriter, request *http.Request) {
	g.logger.Info("Received authorize callback request")

	state := request.URL.Query().Get("state")
	g.logger.Debug("State parameter received", zap.String("state", state))
	if state == "" {
		g.logger.Error("State parameter missing in callback request")
		http.Error(writer, "State parameter missing", http.StatusBadRequest)
		return
	}

	originalParams, err := decodeState(state)
	g.logger.Debug("Decoded original parameters from state", zap.Any("originalParams", originalParams))
	if err != nil {
		g.logger.Error("Failed to decode state", zap.Error(err))
		http.Error(writer, fmt.Sprintf("Failed to decode state: %v", err), http.StatusInternalServerError)
		return
	}

	code := request.FormValue("code")
	g.logger.Debug("Authorization code received", zap.String("code", code))

	token, err := exchangeCodeForToken(code, g.logger)
	if err != nil {
		g.logger.Error("Failed to exchange code for token", zap.Error(err))
		http.Error(writer, fmt.Sprintf("Failed to exchange code for token: %s", err.Error()), http.StatusInternalServerError)
		return
	}
	g.logger.Debug("Successfully exchanged code for token", zap.Any("token", token))

	userInfo, err := getUserInfo(token.AccessToken, g.logger)
	if err != nil {
		g.logger.Error("Failed to get user info", zap.Error(err))
		http.Error(writer, fmt.Sprintf("Failed to get user info: %s", err.Error()), http.StatusInternalServerError)
		return
	}

	g.logger.Info("User info received", zap.Any("userInfo", userInfo))

	user := store.NewUserBuilder().
		WithID(userInfo.ID).
		WithName(userInfo.Name).
		WithEmail(userInfo.Email).
		WithIdpName("Google Authorize").
		Build()
	g.logger.Debug("User entity built from user info (before save)", zap.Any("user", user))

	user, err = g.userRepository.Save(user)

	if err != nil {
		g.logger.Error("Failed to save user", zap.Error(err))
		http.Error(writer, fmt.Sprintf("Failed to save user: %v", err), http.StatusInternalServerError)
		return
	}
	g.logger.Debug("User entity saved (after save)", zap.Any("savedUser", user))

	sessionId, err := g.userSessionService.CreateSession(user.Id, user.Email)

	if err != nil {
		g.logger.Error("Failed to create session", zap.Error(err))
		http.Error(writer, "Failed to create session", http.StatusInternalServerError)
		return
	}
	g.logger.Debug("Session created successfully", zap.String("sessionId", sessionId))

	// Set session cookie
	http.SetCookie(writer, &http.Cookie{
		Name:     "session_id",
		Value:    sessionId,
		Path:     "/",
		HttpOnly: true,  // Prevent client-side scripts from accessing the cookie
		Secure:   false, // Explicitly set to false for local HTTP development
		SameSite: http.SameSiteLaxMode,
		Expires:  time.Now().Add(1 * time.Hour), // Adjust expiration as needed
	})

	g.logger.Debug("Session cookie set",
		zap.String("cookieName", "session_id"),
		zap.String("cookieValue", sessionId),
		zap.String("path", "/"),
		zap.Bool("httpOnly", true),
		zap.Bool("secure", false),
		zap.String("sameSite", "Lax"),
		zap.Time("expires", time.Now().Add(1*time.Hour)),
	)

	// Construct final redirect URL with original parameters
	redirectURL := buildRedirectURL(originalParams)
	g.logger.Debug("Constructed redirect URL", zap.String("redirectURL", redirectURL))
	g.logger.Info("Redirecting to original URL", zap.String("redirectURL", redirectURL))
	http.Redirect(writer, request, redirectURL, http.StatusSeeOther)
}

// exchangeCodeForToken exchanges an authorization code for an access token and other tokens from Google's token endpoint.
func exchangeCodeForToken(code string, logger *zap.Logger) (*TokenResponse, error) {
	data := url.Values{}
	data.Set("code", code)
	data.Set("client_id", configuration.GoogleClientID)
	data.Set("client_secret", configuration.GoogleClientSecret)
	data.Set("redirect_uri", configuration.GoogleRedirectURL)
	data.Set("grant_type", "authorization_code")

	logger.Debug("Exchanging code for token", zap.String("code", code))
	req, err := http.NewRequest("POST", configuration.GoogleTokenURL, bytes.NewBufferString(data.Encode()))
	if err != nil {
		logger.Error("Failed to create token exchange request", zap.Error(err))
		return nil, fmt.Errorf("failed to create token exchange request: %w", err)
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		logger.Error("Failed to perform token exchange request", zap.Error(err))
		return nil, fmt.Errorf("failed to perform token exchange request: %w", err)
	}
	defer func() {
		if cerr := resp.Body.Close(); cerr != nil {
			logger.Error("Error closing token exchange response body", zap.Error(cerr))
		}
	}()

	if resp.StatusCode != http.StatusOK {
		bodyContent, readErr := io.ReadAll(resp.Body)
		if readErr != nil {
			logger.Error("Failed to read token exchange error response body", zap.Error(readErr))
			return nil, fmt.Errorf("token exchange failed: non-OK status %d, but also failed to read response body: %w", resp.StatusCode, readErr)
		}
		logger.Error("Token exchange failed with non-OK status",
			zap.Int("statusCode", resp.StatusCode),
			zap.ByteString("responseBody", bodyContent))
		return nil, fmt.Errorf("token exchange failed: %s", string(bodyContent))
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		logger.Error("Failed to read token exchange successful response body", zap.Error(err))
		return nil, fmt.Errorf("failed to read token exchange successful response body: %w", err)
	}

	var tokenResp TokenResponse
	if err := json.Unmarshal(body, &tokenResp); err != nil {
		logger.Error("Failed to unmarshal token exchange response", zap.Error(err))
		return nil, fmt.Errorf("failed to unmarshal token exchange response: %w", err)
	}

	logger.Info("Successfully exchanged code for token")

	return &tokenResp, nil
}

// getUserInfo retrieves user profile information from Google's user info endpoint using the provided access token.
func getUserInfo(accessToken string, logger *zap.Logger) (*UserInfo, error) {
	logger.Debug("Getting user info", zap.String("accessToken", accessToken))
	req, err := http.NewRequest("GET", configuration.GoogleUserInfoURL, nil)
	if err != nil {
		logger.Error("Failed to create user info request", zap.Error(err))
		return nil, fmt.Errorf("failed to create user info request: %w", err)
	}
	req.Header.Set("Authorization", "Bearer "+accessToken)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		logger.Error("Failed to perform user info request", zap.Error(err))
		return nil, fmt.Errorf("failed to perform user info request: %w", err)
	}
	defer func() {
		if cerr := resp.Body.Close(); cerr != nil {
			logger.Error("Error closing user info response body", zap.Error(cerr))
		}
	}()

	if resp.StatusCode != http.StatusOK {
		bodyContent, readErr := ioutil.ReadAll(resp.Body)
		if readErr != nil {
			logger.Error("Failed to read user info error response body", zap.Error(readErr))
			return nil, fmt.Errorf("user info request failed: non-OK status %d, but also failed to read response body: %w", resp.StatusCode, readErr)
		}
		logger.Error("User info request failed with non-OK status",
			zap.Int("statusCode", resp.StatusCode),
			zap.ByteString("responseBody", bodyContent))
		return nil, fmt.Errorf("user info request failed: %s", string(bodyContent))
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		logger.Error("Failed to read user info successful response body", zap.Error(err))
		return nil, fmt.Errorf("failed to read user info successful response body: %w", err)
	}

	var userInfo UserInfo
	if err := json.Unmarshal(body, &userInfo); err != nil {
		logger.Error("Failed to unmarshal user info response", zap.Error(err))
		return nil, fmt.Errorf("failed to unmarshal user info response: %w", err)
	}
	logger.Info("Successfully retrieved user info", zap.String("userId", userInfo.ID))

	return &userInfo, nil
}

// decodeState decodes a base64 URL-encoded JSON string into a map.
// It is used to retrieve the original authorization request parameters from the state parameter.
func decodeState(encodedState string) (map[string]string, error) {
	stateJSON, err := base64.URLEncoding.DecodeString(encodedState)
	if err != nil {
		return nil, fmt.Errorf("failed to base64 decode state: %w", err)
	}

	var state map[string]string
	if err := json.Unmarshal(stateJSON, &state); err != nil {
		return nil, fmt.Errorf("failed to unmarshal state JSON: %w", err)
	}

	return state, nil
}

// buildRedirectURL constructs the final redirect URL for the client application
// by appending the original authorization parameters to the base URL.
func buildRedirectURL(params map[string]string) string {
	baseURL := "/oauth/authorize" // Change this to your final redirect endpoint

	// Encode parameters for query string
	queryParams := url.Values{}
	for key, value := range params {
		queryParams.Add(key, value)
	}

	// Construct full URL with query parameters
	return fmt.Sprintf("%s?%s", baseURL, queryParams.Encode())
}
