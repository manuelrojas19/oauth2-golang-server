package session

import (
	"go.uber.org/fx"
)

var Module = fx.Options(
	fx.Provide(
		NewRedisClient,
		NewSessionService,
	),
)
