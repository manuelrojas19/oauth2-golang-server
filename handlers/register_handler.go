package handlers

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/manuelrojas19/go-oauth2-server/api"
	"github.com/manuelrojas19/go-oauth2-server/services"
	"github.com/manuelrojas19/go-oauth2-server/utils"
	"go.uber.org/zap"
)

type registerHandler struct {
	oauthClientService services.OauthClientService
	logger             *zap.Logger
}

// NewRegisterHandler creates a new instance of RegisterHandler.
func NewRegisterHandler(oauthClientService services.OauthClientService, logger *zap.Logger) RegisterHandler {
	return &registerHandler{oauthClientService: oauthClientService, logger: logger}
}

// Register processes the registration of a new OAuth client.
func (handler *registerHandler) Register(w http.ResponseWriter, r *http.Request) {
	handler.logger.Info("Received registration request")

	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	var req api.RegisterClientRequest

	if err := utils.DecodeJSON(r, &req); err != nil {
		handler.logger.Error("Error decoding request body", zap.Error(err))
		utils.RespondWithJSON(w, http.StatusBadRequest, api.ErrorResponseBody(api.ErrInvalidRequest))
		return
	}

	req.Sanitize()

	// Validate the request data
	if err := req.Validate(); err != nil {
		handler.logger.Error("Invalid registration request data", zap.Error(err))

		if strings.Contains(err.Error(), "malformed redirect_uri") {
			utils.RespondWithJSON(w, http.StatusBadRequest, api.ErrorResponseBody(api.ErrInvalidRedirectUri, "One or more redirect URIs are invalid or missing"))
		} else {
			utils.RespondWithJSON(w, http.StatusBadRequest, api.ErrorResponseBody(api.ErrInvalidRequest))
		}
		return
	}

	command := services.RegisterOauthClientCommand{
		ClientName:              req.ClientName,
		GrantTypes:              req.GrantTypes,
		ResponseTypes:           req.ResponseTypes,
		TokenEndpointAuthMethod: req.TokenEndpointAuthMethod,
		RedirectUris:            req.RedirectUris,
		Scopes:                  req.Scopes,
	}

	client, err := handler.oauthClientService.CreateOauthClient(&command)
	if err != nil {
		handler.logger.Error("Error creating OAuth client", zap.Error(err))

		if strings.Contains(err.Error(), "client with name") && strings.Contains(err.Error(), "already exists") {
			utils.RespondWithJSON(w, http.StatusConflict, api.ErrorResponseBody(api.ErrClientAlreadyExists, fmt.Sprintf("Client with name '%s' already exists", command.ClientName)))
		} else if strings.Contains(err.Error(), "invalid_scope") {
			utils.RespondWithJSON(w, http.StatusBadRequest, api.ErrorResponseBody(api.ErrInvalidScope, err.Error()))
		} else {
			utils.RespondWithJSON(w, http.StatusInternalServerError, api.ErrorResponseBody(api.ErrServerError, "An unexpected error occurred during client registration."))
		}
		return
	}

	handler.logger.Info("OAuth client created successfully", zap.String("clientId", client.ClientId))
	res := &api.RegisterClientResponse{
		ClientId:                client.ClientId,
		ClientSecret:            client.ClientSecret,
		ClientIdIssuedAt:        fmt.Sprintf("%d", client.ClientIdIssuedAt),
		ClientSecretExpiresAt:   fmt.Sprintf("%d", client.ClientSecretExpiresAt),
		ClientName:              client.ClientName,
		GrantTypes:              client.GrantTypes,
		ResponseTypes:           client.ResponseTypes,
		TokenEndpointAuthMethod: client.TokenEndpointAuthMethod,
		RedirectUris:            client.RedirectUris,
		Scopes:                  client.Scopes,
	}

	utils.RespondWithJSON(w, http.StatusCreated, res)
}
