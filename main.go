package main

import (
	"github.com/manuelrojas19/go-oauth2-server/api/handlers"
	"github.com/manuelrojas19/go-oauth2-server/configuration/database"
	"github.com/manuelrojas19/go-oauth2-server/services"
	"log"
	"net/http"

	"github.com/manuelrojas19/go-oauth2-server/store/repositories"
)

func main() {
	// Initialize database connection
	db, err := database.InitDatabaseConnection()
	if err != nil {
		log.Fatalf("Failed to initialize database connection: %v", err)
		return
	}
	log.Println("Database connection initialized successfully")

	// Initialize repositories and services
	oauthClientRepository := repositories.NewOauthClientRepository(db)
	accessTokenRepository := repositories.NewAccessTokenRepository(db)
	oauthClientService := services.NewOauthClientService(oauthClientRepository)
	tokenService := services.NewTokenService(accessTokenRepository, oauthClientService)
	registerHandler := handlers.NewRegisterHandler(oauthClientService)
	tokenHandler := handlers.NewTokenHandler(tokenService)
	log.Println("Services and handlers initialized successfully")

	// Setup HTTP handler
	http.HandleFunc("/register", registerHandler.Handler)
	http.HandleFunc("/token", tokenHandler.Handler)
	log.Println("HTTP handler for /register is set up")

	// Start HTTP server
	log.Println("Starting HTTP server on :8080")
	err = http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatalf("Failed to start HTTP server: %v", err)
	}
}
