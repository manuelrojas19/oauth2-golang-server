package main

import (
	"fmt"
	"github.com/joho/godotenv"
	"github.com/manuelrojas19/go-oauth2-server/api/handlers"
	"github.com/manuelrojas19/go-oauth2-server/configuration/database"
	"github.com/manuelrojas19/go-oauth2-server/configuration/googleauth"
	"github.com/manuelrojas19/go-oauth2-server/configuration/keymanager"
	"github.com/manuelrojas19/go-oauth2-server/configuration/session"
	"github.com/manuelrojas19/go-oauth2-server/services"
	"github.com/manuelrojas19/go-oauth2-server/store/repositories"
	"log"
	"net/http"
)

func main() {

	err := godotenv.Load()
	if err != nil {
		fmt.Printf("Error loading .env file: %v\n", err)
		return
	}
	// Load Google Secrets
	googleauth.LoadSecrets()
	// Initialize database connection
	db, err := database.InitDatabaseConnection()
	if err != nil {
		log.Fatalf("Failed to initialize database connection: %v", err)
		return
	}
	log.Println("Database connection initialized successfully")

	err = keymanager.Initialize()
	if err != nil {
		log.Fatalf("Failed to initialize keys: %v", err)
	}

	redisClient := session.NewRedisClient()

	userSessionService := services.NewSessionService(redisClient)

	// Initialize repositories and services
	oauthClientRepository := repositories.NewOauthClientRepository(db)
	accessTokenRepository := repositories.NewAccessTokenRepository(db)
	refreshTokenRepository := repositories.NewRefreshTokenRepository(db)
	userConsentRepository := repositories.NewUserConsentRepository(db)
	authorizationRepository := repositories.NewAuthorizationRepository(db)
	userConsentService := services.NewUserConsentService(userConsentRepository)
	oauthClientService := services.NewOauthClientService(oauthClientRepository)
	tokenService := services.NewTokenService(accessTokenRepository, refreshTokenRepository, authorizationRepository, oauthClientService)
	authorizationService := services.NewAuthorizationService(oauthClientService, userConsentService, authorizationRepository, userSessionService)
	wellKnownService := services.NewWellKnownService()
	registerHandler := handlers.NewRegisterHandler(oauthClientService)
	tokenHandler := handlers.NewTokenHandler(tokenService)
	jwksHandler := handlers.NewJwksHandler(wellKnownService)
	authorizeHandler := handlers.NewAuthorizeHandler(authorizationService)
	requestConsentHandler := handlers.NewRequestConsentHandler()
	googleLoginHandler := handlers.NewGoogleLoginHandler()
	googleAuthorizeCallbackHandler := handlers.NewGoogleAuthorizeCallbackHandler(userSessionService)
	log.Println("Services and handlers initialized successfully")

	// Setup HTTP handler
	http.HandleFunc("/oauth/register", registerHandler.Handler)
	http.HandleFunc("/oauth/token", tokenHandler.Handler)
	http.HandleFunc("/oauth/authorize", authorizeHandler.Handler)
	http.HandleFunc("/.well-known/jwks.json", jwksHandler.Handler)
	http.HandleFunc("/oauth/consent", requestConsentHandler.Handler)
	http.HandleFunc("/google/authorize", googleLoginHandler.Handler)
	http.HandleFunc("/google/authorize/callback", googleAuthorizeCallbackHandler.Handler)

	log.Println("HTTP handler for /register is set up")
	log.Println("HTTP handler for /oauth/token is set up")
	log.Println("HTTP handler for /.well-known/jwks.json is set up")

	// Start HTTP server
	log.Println("Starting HTTP server on :8080")
	err = http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatalf("Failed to start HTTP server: %v", err)
	}
}
