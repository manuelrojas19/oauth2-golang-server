package handlers

import (
	"net/http"

	"github.com/manuelrojas19/go-oauth2-server/api"
	"github.com/manuelrojas19/go-oauth2-server/services"
	"github.com/manuelrojas19/go-oauth2-server/utils"
	"go.uber.org/zap"
)

type tokenHandler struct {
	tokenService services.TokenService // A wellKnownService for generating and managing tokens
	logger       *zap.Logger
}

// NewTokenHandler creates a new instance of the handler.
func NewTokenHandler(tokenService services.TokenService, logger *zap.Logger) TokenHandler {
	return &tokenHandler{
		tokenService: tokenService,
		logger:       logger,
	}
}

// Token processes the request for an access token using the Client Credentials Grant flow.
func (handler *tokenHandler) Token(w http.ResponseWriter, r *http.Request) {
	handler.logger.Info("Received token request")

	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	var req api.TokenRequest
	if err := api.DecodeTokenRequest(r, &req); err != nil {
		handler.logger.Error("Error decoding request body", zap.Error(err))
		utils.RespondWithJSON(w, http.StatusBadRequest, api.ErrorResponseBody(api.ErrInvalidRequest))
		return
	}

	// Validate the request data
	if err := req.Validate(); err != nil {
		handler.logger.Error("Invalid token request data", zap.Error(err))
		utils.RespondWithJSON(w, http.StatusBadRequest, api.ErrorResponseBody(api.ErrInvalidRequest))
		return
	}

	grantAccessTokenCommand := services.NewGrantAccessTokenCommand(req.ClientId,
		req.ClientSecret,
		req.GrantType,
		req.RefreshToken,
		req.AuthCode,
		req.RedirectUri,
		req.CodeVerifier,
	)

	// Generate an access token
	token, err := handler.tokenService.GrantAccessToken(grantAccessTokenCommand)
	if err != nil {
		utils.HandleErrorResponse(w, handler.logger, err)
		return
	}

	handler.logger.Info("Access token granted successfully")
	res := api.NewTokenResponse(
		token.AccessToken,
		token.TokenType,
		token.AccessTokenExpiresIn,
		token.RefreshToken,
		utils.JoinStringSlice(token.Scope, " "))

	// Send the response with the token
	utils.RespondWithJSON(w, http.StatusOK, res)
}
