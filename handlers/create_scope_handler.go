package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/manuelrojas19/go-oauth2-server/services"
	"go.uber.org/zap"
)

type scopeHandler struct {
	scopeService services.ScopeService
	logger       *zap.Logger
}

func NewScopeHandler(scopeService services.ScopeService, logger *zap.Logger) ScopeHandler {
	return &scopeHandler{
		scopeService: scopeService,
		logger:       logger,
	}
}

type CreateScopeRequest struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}

type CreateScopeResponse struct {
	Id          string `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
}

// CreateScope handles the creation of a new scope
func (h *scopeHandler) CreateScope(w http.ResponseWriter, r *http.Request) {
	var req CreateScopeRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.logger.Error("Failed to decode request body", zap.Error(err))
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	h.logger.Info("Creating new scope", zap.String("name", req.Name), zap.String("description", req.Description))

	scope, err := h.scopeService.Save(req.Name, req.Description)
	if err != nil {
		h.logger.Error("Failed to create scope", zap.Error(err))
		http.Error(w, "Failed to create scope", http.StatusInternalServerError)
		return
	}

	resp := CreateScopeResponse{
		Id:          scope.Id,
		Name:        scope.Name,
		Description: scope.Description,
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(resp); err != nil {
		h.logger.Error("Failed to encode response", zap.Error(err))
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
		return
	}

	h.logger.Info("Successfully created scope", zap.String("scopeID", scope.Id))
}
