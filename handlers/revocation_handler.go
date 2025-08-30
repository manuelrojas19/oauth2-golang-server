package handlers

import (
	"net/http"

	"github.com/manuelrojas19/go-oauth2-server/api"
	"github.com/manuelrojas19/go-oauth2-server/services"
	"github.com/manuelrojas19/go-oauth2-server/utils"
	"go.uber.org/zap"
)

type revocationHandler struct {
	revocationService services.RevocationService
	log               *zap.Logger
}

func NewRevocationHandler(revocationService services.RevocationService, logger *zap.Logger) RevocationHandler {
	return &revocationHandler{
		revocationService: revocationService,
		log:               logger,
	}
}

func (h *revocationHandler) Revoke(w http.ResponseWriter, r *http.Request) {
	h.log.Info("Entered Revocation handler", zap.String("method", r.Method), zap.String("url", r.URL.String()))

	if r.Method != http.MethodPost {
		h.log.Warn("Invalid request method for Revocation endpoint", zap.String("method", r.Method))
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	if err := r.ParseForm(); err != nil {
		h.log.Error("Failed to parse form data", zap.Error(err))
		utils.RespondWithJSON(w, http.StatusBadRequest, api.ErrorResponseBody(api.ErrInvalidRequest))
		return
	}

	token := r.Form.Get("token")
	if token == "" {
		h.log.Warn("Missing token parameter in revocation request")
		utils.RespondWithJSON(w, http.StatusBadRequest, api.ErrorResponseBody(api.ErrInvalidRequest))
		return
	}

	tokenTypeHint := r.Form.Get("token_type_hint")

	command := &services.RevokeCommand{
		Token:         token,
		TokenTypeHint: tokenTypeHint,
	}

	err := h.revocationService.Revoke(command)
	if err != nil {
		h.log.Error("Error revoking token", zap.Error(err))
		utils.RespondWithJSON(w, http.StatusInternalServerError, api.ErrorResponseBody(api.ErrServerError))
		return
	}

	h.log.Info("Token revoked successfully")
	w.WriteHeader(http.StatusNoContent)
}
