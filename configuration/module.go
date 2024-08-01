package configuration

import (
	"go.uber.org/fx"
)

var Module = fx.Options(
	fx.Provide(
		NewDatabaseConnection,
		NewRedisClient,
		NewServeMux,
	),
	fx.Invoke(
		LoadSecrets,
		InitializeKeys,
	),
)
