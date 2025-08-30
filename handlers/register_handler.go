package handlers

import (
	"log"
	"net/http"

	"github.com/manuelrojas19/go-oauth2-server/api"
	"github.com/manuelrojas19/go-oauth2-server/services"
	"github.com/manuelrojas19/go-oauth2-server/utils"
)

type registerHandler struct {
	oauthClientService services.OauthClientService
}

// NewRegisterHandler creates a new instance of RegisterHandler.
func NewRegisterHandler(oauthClientService services.OauthClientService) RegisterHandler {
	return &registerHandler{oauthClientService: oauthClientService}
}

// Register processes the registration of a new OAuth client.
func (handler *registerHandler) Register(w http.ResponseWriter, r *http.Request) {
	log.Println("Received registration request")

	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	var req api.RegisterClientRequest

	if err := utils.DecodeJSON(r, &req); err != nil {
		log.Printf("Error decoding request body: %v", err)
		utils.RespondWithJSON(w, http.StatusBadRequest, api.ErrorResponseBody(api.ErrInvalidRequest))
		return
	}

	// Validate the request data
	if err := req.Validate(); err != nil {
		log.Println(err);
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
		log.Printf("Error creating OAuth client: %v", err)
		utils.RespondWithJSON(w, http.StatusBadRequest, api.ErrorResponseBody(api.ErrServerError))
		return
	}

	log.Println("OAuth client created successfully")
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
