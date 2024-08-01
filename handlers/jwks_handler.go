package handlers

import (
	"github.com/manuelrojas19/go-oauth2-server/services"
	"github.com/manuelrojas19/go-oauth2-server/utils"
	"net/http"
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
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if jwk == nil {
		http.Error(w, "No JWK jwk available", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	utils.RespondWithJSON(w, http.StatusOK, jwk)
}
