package internal

import (
	"context"
	"github.com/manuelrojas19/go-oauth2-server/handlers"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"go.uber.org/fx"
)

func SetupHTTPServer(
	lc fx.Lifecycle,
	mux *http.ServeMux,
	registerHandler handlers.RegisterHandler,
	tokenHandler handlers.TokenHandler,
	authorizeHandler handlers.AuthorizeHandler,
	requestConsentHandler handlers.RequestConsentHandler,
	jwksHandler handlers.JwksHandler,
	authorizeCallbackHandler handlers.AuthorizeCallbackHandler,
	loginHandler handlers.LoginHandler,
) {
	// Register routes
	mux.HandleFunc("/oauth/register", registerHandler.Register)
	mux.HandleFunc("/oauth/token", tokenHandler.Token)
	mux.HandleFunc("/oauth/authorize", authorizeHandler.Authorize)
	mux.HandleFunc("/oauth/consent", requestConsentHandler.RequestConsent)
	mux.HandleFunc("/oauth/login", loginHandler.Login)
	mux.HandleFunc("/.well-known/jwks.json", jwksHandler.Jwks)
	mux.HandleFunc("/google/authorize/callback", authorizeCallbackHandler.ProcessCallback)
	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	})

	// Register lifecycle hooks
	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			log.Println("Starting HTTP server on :8080")
			server := &http.Server{
				Addr:    ":8080",
				Handler: mux,
			}

			go func() {
				if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
					log.Fatalf("Error starting HTTP server: %v", err)
				}
			}()

			go func() {
				stop := make(chan os.Signal, 1)
				signal.Notify(stop, os.Interrupt)
				<-stop

				log.Println("Received interrupt signal, shutting down gracefully...")
				ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
				defer cancel()

				if err := server.Shutdown(ctx); err != nil {
					log.Fatalf("Server forced to shutdown: %v", err)
				}
			}()

			return nil
		},
		OnStop: func(ctx context.Context) error {
			log.Println("Stopping HTTP server")
			// Here you can add additional cleanup logic if necessary
			return nil
		},
	})
}
