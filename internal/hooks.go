package internal

import (
	"context"
	"go.uber.org/zap"
	"net/http"
	"os"
	"os/signal"
	"time"

	"go.uber.org/fx"
)

func SetupHTTPServer(
	lc fx.Lifecycle,
	mux *http.ServeMux,
	log *zap.Logger,
	routes map[string]http.HandlerFunc,
) {
	// Register routes dynamically
	for path, handler := range routes {
		log.Info("Adding route", zap.String("path", path))
		mux.HandleFunc(path, handler)
	}

	// Register lifecycle hooks
	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			log.Info("Starting HTTP server on :8080")
			server := &http.Server{
				Addr:    ":8080",
				Handler: mux,
			}

			go func() {
				if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
					log.Error("Error starting HTTP server", zap.Error(err))
				}
			}()

			go func() {
				stop := make(chan os.Signal, 1)
				signal.Notify(stop, os.Interrupt)
				<-stop

				log.Info("Received interrupt signal, shutting down gracefully...")
				ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
				defer cancel()

				if err := server.Shutdown(ctx); err != nil {
					log.Error("Server forced to shutdown", zap.Error(err))
				}
			}()

			return nil
		},
		OnStop: func(ctx context.Context) error {
			log.Info("Stopping HTTP server")
			// Here you can add additional cleanup logic if necessary
			return nil
		},
	})
}
