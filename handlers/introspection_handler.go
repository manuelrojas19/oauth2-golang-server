package handlers

import (
	"net/http"

	"github.com/manuelrojas19/go-oauth2-server/api"
	"github.com/manuelrojas19/go-oauth2-server/services"
	"github.com/manuelrojas19/go-oauth2-server/utils"
	"go.uber.org/zap"
)

type introspectionHandler struct {
	introspectionService services.IntrospectionService
	log                  *zap.Logger
}

func NewIntrospectionHandler(introspectionService services.IntrospectionService, logger *zap.Logger) IntrospectionHandler {
	return &introspectionHandler{
		introspectionService: introspectionService,
		log:                  logger,
	}
}

func (h *introspectionHandler) Introspect(w http.ResponseWriter, r *http.Request) {
	h.log.Info("Entered Introspection handler", zap.String("method", r.Method), zap.String("url", r.URL.String()))

	if r.Method != http.MethodPost {
		h.log.Warn("Invalid request method for Introspection endpoint", zap.String("method", r.Method))
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
		h.log.Warn("Missing token parameter in introspection request")
		utils.RespondWithJSON(w, http.StatusBadRequest, api.ErrorResponseBody(api.ErrInvalidRequest))
		return
	}

	tokenTypeHint := r.Form.Get("token_type_hint")

	command := &services.IntrospectCommand{
		Token:         token,
		TokenTypeHint: tokenTypeHint,
	}

	introspectionResponse, err := h.introspectionService.Introspect(command)
	if err != nil {
		h.log.Error("Error performing introspection", zap.Error(err))
		utils.RespondWithJSON(w, http.StatusInternalServerError, api.ErrorResponseBody(api.ErrServerError))
		return
	}

	h.log.Info("Introspection successful", zap.Any("response", introspectionResponse))
	utils.RespondWithJSON(w, http.StatusOK, introspectionResponse)
}
