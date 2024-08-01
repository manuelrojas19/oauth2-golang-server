package configuration

import (
	"go.uber.org/fx"
	"go.uber.org/fx/fxevent"
	"go.uber.org/zap"
)

// Module this module defines a set of options for the fx application, including logger configuration,
// provided dependencies, and functions to be invoked at startup.
var Module = fx.Options(
	// Configure the logger for the fx application using a custom Zap logger
	fx.WithLogger(func(log *zap.Logger) fxevent.Logger {
		// Return an fxevent.Logger that uses the provided Zap logger
		return &fxevent.ZapLogger{Logger: log}
	}),

	// Provide various dependencies to the fx application
	fx.Provide(
		// Provide a new database connection
		NewDatabaseConnection,

		// Provide a new Redis client
		NewRedisClient,

		// Provide a new HTTP ServeMux
		NewServeMux,

		// Provide a new logger
		NewLogger,
	),

	// Invoke functions to be executed at application startup
	fx.Invoke(
		// Load secrets required by the application
		LoadSecrets,

		// Initialize cryptographic keys
		InitializeKeys,
	),
)
