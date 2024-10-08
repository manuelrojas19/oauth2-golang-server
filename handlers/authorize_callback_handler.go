package handlers

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/manuelrojas19/go-oauth2-server/configuration"
	"github.com/manuelrojas19/go-oauth2-server/services"
	"github.com/manuelrojas19/go-oauth2-server/store"
	"github.com/manuelrojas19/go-oauth2-server/store/repositories"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"time"
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
}

func NewAuthorizeCallbackHandler(
	userSessionService services.SessionService,
	userRepository repositories.UserRepository,
) AuthorizeCallbackHandler {
	return &authorizeCallbackHandler{
		userSessionService: userSessionService,
		userRepository:     userRepository}
}

func (g authorizeCallbackHandler) ProcessCallback(writer http.ResponseWriter, request *http.Request) {

	state := request.URL.Query().Get("state")
	if state == "" {
		http.Error(writer, "State parameter missing", http.StatusBadRequest)
		return
	}

	originalParams, err := decodeState(state)
	if err != nil {
		http.Error(writer, fmt.Sprintf("Failed to decode state: %v", err), http.StatusInternalServerError)
		return
	}

	code := request.FormValue("code")
	token, err := exchangeCodeForToken(code)
	if err != nil {
		http.Error(writer, fmt.Sprintf("Failed to exchange code for token: %s", err.Error()), http.StatusInternalServerError)
		return
	}

	userInfo, err := getUserInfo(token.AccessToken)
	if err != nil {
		http.Error(writer, fmt.Sprintf("Failed to get user info: %s", err.Error()), http.StatusInternalServerError)
		return
	}

	log.Printf("User info: %v", userInfo)

	user := store.NewUserBuilder().
		WithID(userInfo.ID).
		WithName(userInfo.Name).
		WithEmail(userInfo.Email).
		WithIdpName("Google Authorize").
		Build()

	user, err = g.userRepository.Save(user)

	if err != nil {
		http.Error(writer, fmt.Sprintf("Failed to save user: %v", err), http.StatusInternalServerError)
		return
	}

	sessionId, err := g.userSessionService.CreateSession(user.Id, user.Email)

	if err != nil {
		log.Printf("Failed to create session: %s", err.Error())
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
	log.Printf("Redirect URL: %v", redirectURL)
	http.Redirect(writer, request, redirectURL, http.StatusSeeOther)
}

func exchangeCodeForToken(code string) (*TokenResponse, error) {
	data := url.Values{}
	data.Set("code", code)
	data.Set("client_id", configuration.GoogleClientID)
	data.Set("client_secret", configuration.GoogleClientSecret)
	data.Set("redirect_uri", configuration.GoogleRedirectURL)
	data.Set("grant_type", "authorization_code")

	req, err := http.NewRequest("POST", configuration.GoogleTokenURL, bytes.NewBufferString(data.Encode()))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var tokenResp TokenResponse
	if err := json.Unmarshal(body, &tokenResp); err != nil {
		return nil, err
	}

	return &tokenResp, nil
}

func getUserInfo(accessToken string) (*UserInfo, error) {
	req, err := http.NewRequest("GET", configuration.GoogleUserInfoURL, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", "Bearer "+accessToken)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {

		}
	}(resp.Body)

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var userInfo UserInfo
	if err := json.Unmarshal(body, &userInfo); err != nil {
		return nil, err
	}

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
