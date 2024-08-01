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

func main() {
	app := fx.New(

		configuration.Module,
		session.Module,
		repositories.Module,
		services.Module,
		handlers.Module,
		fx.Invoke(internal.SetupHTTPServer),
	)
	app.Run()
}
