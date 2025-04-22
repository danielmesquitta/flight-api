package router

import (
	"github.com/gofiber/fiber/v2"

	_ "github.com/danielmesquitta/flight-api/docs" // swagger docs
	"github.com/danielmesquitta/flight-api/internal/app/server/handler"
	"github.com/danielmesquitta/flight-api/internal/app/server/middleware"
	"github.com/danielmesquitta/flight-api/internal/config/env"
)

type Router struct {
	e  *env.Env
	m  *middleware.Middleware
	hh *handler.HealthHandler
	dh *handler.DocHandler
}

func NewRouter(
	e *env.Env,
	m *middleware.Middleware,
	hh *handler.HealthHandler,
	dh *handler.DocHandler,
) *Router {
	return &Router{
		e:  e,
		m:  m,
		hh: hh,
		dh: dh,
	}
}

func (r *Router) Register(
	app *fiber.App,
) {
	basePath := "/api"
	api := app.Group(basePath)

	api.Get("/health", r.hh.Health)
	api.Use("/docs", r.dh.Get)

	apiV1 := app.Group(basePath + "/v1")

	_ = apiV1
}
