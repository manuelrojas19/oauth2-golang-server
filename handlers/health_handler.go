package handlers

import (
	"net/http"

	"go.uber.org/zap"
)

type HealthHandler interface {
	Health(w http.ResponseWriter, r *http.Request)
}

type healthHandler struct {
	log *zap.Logger
}

func NewHealthHandler(logger *zap.Logger) HealthHandler {
	return &healthHandler{log: logger}
}

func (h healthHandler) Health(w http.ResponseWriter, r *http.Request) {
	h.log.Info("Received health check request")
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte("OK"))
}
