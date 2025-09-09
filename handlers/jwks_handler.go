package handlers

import (
	"net/http"

	"github.com/manuelrojas19/go-oauth2-server/api"
	"github.com/manuelrojas19/go-oauth2-server/services"
	"github.com/manuelrojas19/go-oauth2-server/utils"
	"go.uber.org/zap"
)

type jwksHandler struct {
	wellKnownService services.WellKnownService
	logger           *zap.Logger
}

func NewJwksHandler(wellKnownService services.WellKnownService, logger *zap.Logger) JwksHandler {
	return &jwksHandler{
		wellKnownService: wellKnownService,
		logger:           logger,
	}
}

func (j jwksHandler) Jwks(w http.ResponseWriter, r *http.Request) {

	j.logger.Info("Received JWKS request", zap.String("method", r.Method), zap.String("url", r.URL.String()))

	if r.Method != http.MethodGet {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	jwk, err := j.wellKnownService.GetJwk()
	if err != nil {
		j.logger.Error("Error retrieving JWK", zap.Error(err))
		utils.HandleErrorResponse(w, j.logger, err)
		return
	}

	if jwk == nil {
		j.logger.Warn("JWK not found or empty")
		utils.RespondWithJSON(w, http.StatusInternalServerError, api.ErrorResponseBody(api.ErrServerError, "JWK not found or empty"))
		return
	}

	w.Header().Set("Content-Type", "application/json")
	j.logger.Info("Successfully retrieved JWK")
	utils.RespondWithJSON(w, http.StatusOK, jwk)
}
