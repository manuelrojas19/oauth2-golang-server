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

type TokenResponse struct {
	AccessToken  string `json:"access_token"`
	TokenType    string `json:"token_type"`
	RefreshToken string `json:"refresh_token"`
	Expiry       int64  `json:"expires_in"`
	IDToken      string `json:"id_token"`
}

type UserInfo struct {
	ID    string `json:"sub"`
	Email string `json:"email"`
	Name  string `json:"name"`
}

type authorizeCallbackHandler struct {
	userSessionService services.SessionService
	userRepository     repositories.UserRepository
	logger             *zap.Logger
}

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

func (g authorizeCallbackHandler) ProcessCallback(writer http.ResponseWriter, request *http.Request) {
	g.logger.Info("Received authorize callback request")

	state := request.URL.Query().Get("state")
	if state == "" {
		g.logger.Error("State parameter missing in callback request")
		http.Error(writer, "State parameter missing", http.StatusBadRequest)
		return
	}

	originalParams, err := decodeState(state)
	if err != nil {
		g.logger.Error("Failed to decode state", zap.Error(err))
		http.Error(writer, fmt.Sprintf("Failed to decode state: %v", err), http.StatusInternalServerError)
		return
	}

	code := request.FormValue("code")
	token, err := exchangeCodeForToken(code, g.logger)
	if err != nil {
		g.logger.Error("Failed to exchange code for token", zap.Error(err))
		http.Error(writer, fmt.Sprintf("Failed to exchange code for token: %s", err.Error()), http.StatusInternalServerError)
		return
	}

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

	user, err = g.userRepository.Save(user)

	if err != nil {
		g.logger.Error("Failed to save user", zap.Error(err))
		http.Error(writer, fmt.Sprintf("Failed to save user: %v", err), http.StatusInternalServerError)
		return
	}

	sessionId, err := g.userSessionService.CreateSession(user.Id, user.Email)

	if err != nil {
		g.logger.Error("Failed to create session", zap.Error(err))
		http.Error(writer, "Failed to create session", http.StatusInternalServerError)
		return
	}

	// Set session cookie
	http.SetCookie(writer, &http.Cookie{
		Name:     "session_id",
		Value:    sessionId,
		Path:     "/",
		HttpOnly: true, // Prevent client-side scripts from accessing the cookie
		Secure:   true, // Use HTTPS in production
		SameSite: http.SameSiteStrictMode,
		Expires:  time.Now().Add(1 * time.Hour), // Adjust expiration as needed
	})

	// Construct final redirect URL with original parameters
	redirectURL := buildRedirectURL(originalParams)
	g.logger.Info("Redirecting to original URL", zap.String("redirectURL", redirectURL))
	http.Redirect(writer, request, redirectURL, http.StatusSeeOther)
}

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
		return nil, err
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		logger.Error("Failed to perform token exchange request", zap.Error(err))
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		bodyContent, _ := ioutil.ReadAll(resp.Body)
		logger.Error("Token exchange failed with non-OK status",
			zap.Int("statusCode", resp.StatusCode),
			zap.ByteString("responseBody", bodyContent))
		return nil, fmt.Errorf("token exchange failed: %s", string(bodyContent))
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		logger.Error("Failed to read token exchange response body", zap.Error(err))
		return nil, err
	}

	var tokenResp TokenResponse
	if err := json.Unmarshal(body, &tokenResp); err != nil {
		logger.Error("Failed to unmarshal token exchange response", zap.Error(err))
		return nil, err
	}
	logger.Info("Successfully exchanged code for token")

	return &tokenResp, nil
}

func getUserInfo(accessToken string, logger *zap.Logger) (*UserInfo, error) {
	logger.Debug("Getting user info", zap.String("accessToken", accessToken))
	req, err := http.NewRequest("GET", configuration.GoogleUserInfoURL, nil)
	if err != nil {
		logger.Error("Failed to create user info request", zap.Error(err))
		return nil, err
	}
	req.Header.Set("Authorization", "Bearer "+accessToken)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		logger.Error("Failed to perform user info request", zap.Error(err))
		return nil, err
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			logger.Error("Error closing response body", zap.Error(err))
		}
	}(resp.Body)

	if resp.StatusCode != http.StatusOK {
		bodyContent, _ := ioutil.ReadAll(resp.Body)
		logger.Error("User info request failed with non-OK status",
			zap.Int("statusCode", resp.StatusCode),
			zap.ByteString("responseBody", bodyContent))
		return nil, fmt.Errorf("user info request failed: %s", string(bodyContent))
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		logger.Error("Failed to read user info response body", zap.Error(err))
		return nil, err
	}

	var userInfo UserInfo
	if err := json.Unmarshal(body, &userInfo); err != nil {
		logger.Error("Failed to unmarshal user info response", zap.Error(err))
		return nil, err
	}
	logger.Info("Successfully retrieved user info", zap.String("userId", userInfo.ID))

	return &userInfo, nil
}

func decodeState(encodedState string) (map[string]string, error) {
	stateJSON, err := base64.URLEncoding.DecodeString(encodedState)
	if err != nil {
		return nil, err
	}

	var state map[string]string
	if err := json.Unmarshal(stateJSON, &state); err != nil {
		return nil, err
	}

	return state, nil
}

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
