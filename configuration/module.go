package configuration

import (
	"go.uber.org/fx"
	"go.uber.org/fx/fxevent"
	"go.uber.org/zap"
)

var Module = fx.Options(
	fx.WithLogger(func(log *zap.Logger) fxevent.Logger {
		return &fxevent.ZapLogger{Logger: log}
	}),
	fx.Provide(
		NewDatabaseConnection,
		NewRedisClient,
		NewServeMux,
		NewLogger,
	),
	fx.Invoke(
		LoadSecrets,
		InitializeKeys,
	),
)
