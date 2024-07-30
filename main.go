package main

import (
	"fmt"
	"github.com/joho/godotenv"
	"github.com/manuelrojas19/go-oauth2-server/configuration"
	"github.com/manuelrojas19/go-oauth2-server/handlers"
	"github.com/manuelrojas19/go-oauth2-server/idp"
	"github.com/manuelrojas19/go-oauth2-server/services"
	"github.com/manuelrojas19/go-oauth2-server/session"
	"github.com/manuelrojas19/go-oauth2-server/store"
	"log"
	"net/http"
)

func main() {
	// Load Google Secrets
	err := godotenv.Load()
	if err != nil {
		fmt.Println("Error loading .env file")
		return
	}

	configuration.LoadGoogleSecrets()
	configuration.LoadDbSecrets()
	configuration.LoadRedisSecrets()

	// Initialize database connection
	db, err := configuration.InitDatabaseConnection()
	if err != nil {
		log.Fatalf("Failed to initialize database connection: %v", err)
		return
	}
	log.Println("Database connection initialized successfully")

	err = configuration.Initialize()
	if err != nil {
		log.Fatalf("Failed to initialize keys: %v", err)
	}

	redisClient := session.NewRedisClient()
	userSessionService := session.NewSessionService(redisClient)

	// Initialize repositories and services
	oauthClientRepository := store.NewOauthClientRepository(db)
	accessTokenRepository := store.NewAccessTokenRepository(db)
	refreshTokenRepository := store.NewRefreshTokenRepository(db)
	userConsentRepository := store.NewUserConsentRepository(db)
	authorizationRepository := store.NewAuthorizationRepository(db)
	userRepository := store.NewUserRepository(db)
	userConsentService := services.NewUserConsentService(userConsentRepository)
	oauthClientService := services.NewOauthClientService(oauthClientRepository)
	tokenService := services.NewTokenService(accessTokenRepository, refreshTokenRepository, authorizationRepository, oauthClientService)
	authorizationService := services.NewAuthorizationService(oauthClientService, userConsentService, authorizationRepository, userSessionService, userRepository)
	wellKnownService := services.NewWellKnownService()
	registerHandler := handlers.NewRegisterHandler(oauthClientService)
	tokenHandler := handlers.NewTokenHandler(tokenService)
	jwksHandler := handlers.NewJwksHandler(wellKnownService)
	authorizeHandler := handlers.NewAuthorizeHandler(authorizationService)
	requestConsentHandler := handlers.NewRequestConsentHandler()
	googleAuthorizeCallbackHandler := idp.NewGoogleAuthorizeCallbackHandler(userSessionService, userRepository)
	loginHandler := handlers.NewLoginHandler()
	log.Println("Services and handlers initialized successfully")

	// Setup HTTP handler
	http.HandleFunc("/oauth/register", registerHandler.Handler)
	http.HandleFunc("/oauth/token", tokenHandler.Handler)
	http.HandleFunc("/oauth/authorize", authorizeHandler.Handler)
	http.HandleFunc("/oauth/consent", requestConsentHandler.Handler)
	http.HandleFunc("/oauth/login", loginHandler.Handler)
	http.HandleFunc("/.well-known/jwks.json", jwksHandler.Handler)
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
