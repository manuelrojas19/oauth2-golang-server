package handlers

import (
	"errors"
	"github.com/manuelrojas19/go-oauth2-server/api/dto/request"
	"github.com/manuelrojas19/go-oauth2-server/api/dto/response"
	"github.com/manuelrojas19/go-oauth2-server/services"
	"github.com/manuelrojas19/go-oauth2-server/services/commands"
	"github.com/manuelrojas19/go-oauth2-server/utils"
	"log"
	"net/http"
)

type tokenHandler struct {
	tokenService services.TokenService // A service for generating and managing tokens
}

// NewTokenHandler creates a new instance of the handler.
func NewTokenHandler(tokenService services.TokenService) Handler {
	return &tokenHandler{
		tokenService: tokenService,
	}
}

// Handler processes the request for an access token using the Client Credentials Grant flow.
func (handler *tokenHandler) Handler(w http.ResponseWriter, r *http.Request) {
	log.Println("Received token request")

	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	var req request.TokenRequest
	if err := utils.DecodeTokenRequest(r, &req); err != nil {
		log.Printf("Error decoding request body: %v", err)
		utils.RespondWithJSON(w, http.StatusBadRequest, utils.ErrorResponseBody(err))
		return
	}

	// Validate the request data
	if err := req.Validate(); err != nil {
		utils.RespondWithJSON(w, http.StatusBadRequest, utils.ErrorResponseBody(err))
		return
	}

	// Validate the grant_type
	if !(req.GrantType == "client_credentials" || req.GrantType == "refresh_token") {
		utils.RespondWithJSON(w, http.StatusBadRequest, utils.ErrorResponseBody(errors.New("invalid grant_type")))
		return
	}

	grantAccessTokenCommand := commands.NewGrantAccessTokenCommand(req.ClientId, req.ClientSecret, req.GrantType, req.RefreshToken)

	// Generate an access token
	token, err := handler.tokenService.GrantAccessToken(grantAccessTokenCommand)
	if err != nil {
		log.Printf("Error generating token: %v", err)
		utils.RespondWithJSON(w, http.StatusInternalServerError, utils.ErrorResponseBody(err))
		return
	}

	res := response.NewTokenResponse(token.AccessToken,
		"Bearer",
		services.AccessTokenDuration.String(),
		token.RefreshToken)

	// Send the response with the token
	utils.RespondWithJSON(w, http.StatusOK, res)
}
