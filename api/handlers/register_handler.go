package handlers

import (
	"github.com/manuelrojas19/go-oauth2-server/api/request"
	"github.com/manuelrojas19/go-oauth2-server/mappers"
	"github.com/manuelrojas19/go-oauth2-server/services"
	"github.com/manuelrojas19/go-oauth2-server/utils"
	"log"
	"net/http"
)

type registerHandler struct {
	oauthClientService services.OauthClientService
}

// NewRegisterHandler creates a new instance of Handler.
func NewRegisterHandler(oauthClientService services.OauthClientService) Handler {
	return &registerHandler{oauthClientService: oauthClientService}
}

// Handler processes the registration of a new OAuth client.
func (handler *registerHandler) Handler(w http.ResponseWriter, r *http.Request) {
	log.Println("Received registration request")

	var req request.RegisterClientRequest

	if err := utils.Decode(r, &req); err != nil {
		log.Printf("Error decoding request body: %v", err)
		utils.RespondWithJSON(w, http.StatusBadRequest, utils.ErrorResponseBody(err))
		return
	}

	command := mappers.NewCreateOauthClientCommandFromRequest(&req)

	client, err := handler.oauthClientService.CreateOauthClient(command)
	if err != nil {
		log.Printf("Error creating OAuth client: %v", err)
		utils.RespondWithJSON(w, http.StatusBadRequest, utils.ErrorResponseBody(err))
		return
	}

	log.Println("OAuth client created successfully")
	res := mappers.NewRegisterClientResponseFromClientModel(client)
	utils.RespondWithJSON(w, http.StatusCreated, res)
}
