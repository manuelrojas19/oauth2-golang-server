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
	routes handlers.Routes,
) {
	// Register routes dynamically
	for path, handler := range routes.Routes() {
		mux.HandleFunc(path, handler)
	}

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
