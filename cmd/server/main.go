package main

import (
	"github.com/manuelrojas19/go-oauth2-server/configuration"
	"github.com/manuelrojas19/go-oauth2-server/handlers"
	"github.com/manuelrojas19/go-oauth2-server/internal"
	"github.com/manuelrojas19/go-oauth2-server/services"
	"github.com/manuelrojas19/go-oauth2-server/session"
	"github.com/manuelrojas19/go-oauth2-server/store/repositories"
	"go.uber.org/fx"
)

// main is the entry point of the application. It uses the fx framework to set up and run the application with various modules and dependencies.
func main() {
	// Create a new fx.App with the specified modules and invoke functions
	app := fx.New(
		// Include the configuration module which provides application configuration settings
		configuration.Module,

		// Include the session module which handles session management
		session.Module,

		// Include the repositories module which provides data access layer (DAL) functionality
		repositories.Module,

		// Include the services module which contains business logic and services
		services.Module,

		// Include the handlers module which defines the HTTP handlers for the application
		handlers.Module,

		// Invoke the SetupHTTPServer function from the internal package to set up the HTTP server
		fx.Invoke(internal.SetupHTTPServer),
	)

	// Run the fx.App, starting the application and its dependencies
	app.Run()
}
