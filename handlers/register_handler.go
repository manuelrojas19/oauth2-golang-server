package handlers

import (
	"github.com/manuelrojas19/go-oauth2-server/models/request"
	"github.com/manuelrojas19/go-oauth2-server/services"
	"github.com/manuelrojas19/go-oauth2-server/utils"
	"log"
	"net/http"
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

	var requestBody request.RegisterClientRequest
	if err := utils.Decode(r, &requestBody); err != nil {
		log.Printf("Error decoding request body: %v", err)
		utils.RespondWithJSON(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
		return
	}

	client, err := handler.oauthClientService.CreateOauthClient(requestBody.RedirectUris)
	if err != nil {
		log.Printf("Error creating OAuth client: %v", err)
		utils.RespondWithJSON(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
		return
	}

	log.Println("OAuth client created successfully")
	utils.RespondWithJSON(w, http.StatusCreated, client)
}
