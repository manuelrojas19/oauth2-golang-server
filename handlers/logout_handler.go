package handlers

import (
	"net/http"

	"github.com/manuelrojas19/go-oauth2-server/services"
	"go.uber.org/zap"
)

type logoutHandler struct {
	sessionService services.SessionService
	log            *zap.Logger
}

func NewLogoutHandler(sessionService services.SessionService, logger *zap.Logger) LogoutHandler {
	return &logoutHandler{
		sessionService: sessionService,
		log:            logger,
	}
}

func (h *logoutHandler) Logout(w http.ResponseWriter, r *http.Request) {
	h.log.Info("Entered Logout handler", zap.String("method", r.Method), zap.String("url", r.URL.String()))

	if r.Method != http.MethodGet {
		h.log.Warn("Invalid request method for Logout endpoint", zap.String("method", r.Method))
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	cookie, err := r.Cookie("session_id")
	if err != nil || cookie.Value == "" {
		h.log.Info("No session cookie found, user is already logged out or never logged in")
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	sessionId := cookie.Value
	err = h.sessionService.DeleteSession(sessionId)
	if err != nil {
		h.log.Error("Error deleting session", zap.Error(err), zap.String("sessionId", sessionId))
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:     "session_id",
		Value:    "",
		Path:     "/",
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteStrictMode,
		MaxAge:   -1, // Immediately expire the cookie
	})

	h.log.Info("User logged out successfully", zap.String("sessionId", sessionId))
	http.Redirect(w, r, "/", http.StatusSeeOther)
}
