package handlers

import (
	"net/http"
	"strings"

	"github.com/manuelrojas19/go-oauth2-server/api"
	"github.com/manuelrojas19/go-oauth2-server/services"
	"github.com/manuelrojas19/go-oauth2-server/utils"
	"go.uber.org/zap"
)

type userinfoHandler struct {
	userinfoService services.UserinfoService
	log             *zap.Logger
}

func NewUserinfoHandler(userinfoService services.UserinfoService, logger *zap.Logger) UserinfoHandler {
	return &userinfoHandler{
		userinfoService: userinfoService,
		log:             logger,
	}
}

func (h *userinfoHandler) Userinfo(w http.ResponseWriter, r *http.Request) {
	h.log.Info("Entered Userinfo handler", zap.String("method", r.Method), zap.String("url", r.URL.String()))

	if r.Method != http.MethodGet {
		h.log.Warn("Invalid request method for Userinfo endpoint", zap.String("method", r.Method))
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	authHeader := r.Header.Get("Authorization")
	if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
		h.log.Warn("Missing or invalid Authorization header")
		utils.RespondWithJSON(w, http.StatusUnauthorized, api.ErrorResponseBody(api.ErrInvalidRequest))
		return
	}

	accessToken := strings.TrimPrefix(authHeader, "Bearer ")
	h.log.Debug("Extracted access token from header", zap.String("accessToken", accessToken))

	command := &services.GetUserinfoCommand{
		AccessToken: accessToken,
	}
	h.log.Debug("Created GetUserinfoCommand", zap.Any("command", command))

	userinfo, err := h.userinfoService.GetUserinfo(command)
	if err != nil {
		h.log.Error("Error retrieving user info", zap.Error(err), zap.String("accessToken", accessToken))
		utils.RespondWithJSON(w, http.StatusUnauthorized, api.ErrorResponseBody(api.ErrInvalidToken))
		return
	}

	h.log.Info("Userinfo retrieved successfully", zap.Any("userinfo", userinfo))
	utils.RespondWithJSON(w, http.StatusOK, userinfo)
}
