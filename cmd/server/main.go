package main

import (
	"log"

	"github.com/danielmesquitta/flight-api/internal/app/server"
	"github.com/danielmesquitta/flight-api/internal/config"
	"github.com/danielmesquitta/flight-api/internal/config/env"
	"github.com/danielmesquitta/flight-api/internal/pkg/validator"
)

// @title Flight API
// @version 1.0
// @description Flight API
// @contact.name Daniel Mesquita
// @contact.email danielmesquitta123@gmail.com
// @BasePath /api
// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @description Type "Bearer" followed by a space and JWT token.
// @securityDefinitions.basic BasicAuth
func main() {
	v := validator.New()
	e := config.LoadConfig(v)

	var app *server.App
	switch e.Environment {
	case env.EnvironmentProduction:
		app = server.NewProd(v, e, nil)

	case env.EnvironmentTest:
		app = server.NewTest(v, e, nil)

	case env.EnvironmentStaging:
		app = server.NewStaging(v, e, nil)

	default:
		app = server.NewDev(v, e, nil)
	}

	if err := app.Listen(":" + e.Port); err != nil {
		log.Fatalf("failed to start server: %v", err)
	}
}
