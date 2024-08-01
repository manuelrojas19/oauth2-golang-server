package handlers

import (
	"net/http"
)

type HealthHandler interface {
	Health(w http.ResponseWriter, r *http.Request)
}

type healthHandler struct {
}

func NewHealthHandler() HealthHandler {
	return &healthHandler{}
}

func (h healthHandler) Health(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte("OK"))
}
