package handlers

import (
	"net/http"

	"github.com/manuelrojas19/go-oauth2-server/api"
	"github.com/manuelrojas19/go-oauth2-server/services"
	"github.com/manuelrojas19/go-oauth2-server/utils"
)

type jwksHandler struct {
	wellKnownService services.WellKnownService
}

func NewJwksHandler(wellKnownService services.WellKnownService) JwksHandler {
	return &jwksHandler{
		wellKnownService: wellKnownService,
	}
}

func (j jwksHandler) Jwks(w http.ResponseWriter, _ *http.Request) {
	jwk, err := j.wellKnownService.GetJwk()
	if err != nil {
		utils.RespondWithJSON(w, http.StatusInternalServerError, api.ErrorResponseBody(api.ErrServerError))
		return
	}

	if jwk == nil {
		utils.RespondWithJSON(w, http.StatusNotFound, api.ErrorResponseBody(api.ErrServerError))
		return
	}

	w.Header().Set("Content-Type", "application/json")

	utils.RespondWithJSON(w, http.StatusOK, jwk)
}
