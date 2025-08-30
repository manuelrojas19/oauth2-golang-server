package handlers

import (
	"net/http"

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

	// Validate the request data
	if err := req.Validate(); err != nil {
		handler.logger.Error("Invalid registration request data", zap.Error(err))
		utils.RespondWithJSON(w, http.StatusBadRequest, api.ErrorResponseBody(api.ErrInvalidRequest))
		return
	}

	command := &services.RegisterOauthClientCommand{
		ClientName:              req.ClientName,
		GrantTypes:              req.GrantTypes,
		ResponseTypes:           req.ResponseTypes,
		TokenEndpointAuthMethod: req.TokenEndpointAuthMethod,
		RedirectUris:            req.RedirectUris,
		Scopes:                  req.Scopes,
	}

	client, err := handler.oauthClientService.CreateOauthClient(command)
	if err != nil {
		handler.logger.Error("Error creating OAuth client", zap.Error(err))
		utils.RespondWithJSON(w, http.StatusBadRequest, api.ErrorResponseBody(api.ErrServerError))
		return
	}

	handler.logger.Info("OAuth client created successfully", zap.String("clientId", client.ClientId))
	res := &api.RegisterClientResponse{
		ClientId:                client.ClientId,
		ClientSecret:            client.ClientSecret,
		ClientName:              client.ClientName,
		GrantTypes:              client.GrantTypes,
		ResponseTypes:           client.ResponseTypes,
		TokenEndpointAuthMethod: client.TokenEndpointAuthMethod,
		RedirectUris:            client.RedirectUris,
		Scopes:                  client.Scopes,
	}

	utils.RespondWithJSON(w, http.StatusCreated, res)
}
