package wire

import (
	"testing"

	"github.com/google/wire"
	"github.com/stretchr/testify/mock"

	"github.com/danielmesquitta/flight-api/internal/app/server"
	"github.com/danielmesquitta/flight-api/internal/app/server/handler"
	"github.com/danielmesquitta/flight-api/internal/app/server/middleware"
	"github.com/danielmesquitta/flight-api/internal/app/server/router"
	"github.com/danielmesquitta/flight-api/internal/config/env"
	"github.com/danielmesquitta/flight-api/internal/domain/usecase/flight"
	"github.com/danielmesquitta/flight-api/internal/pkg/jwtutil"
	"github.com/danielmesquitta/flight-api/internal/pkg/validator"
	"github.com/danielmesquitta/flight-api/internal/provider/cache"
	"github.com/danielmesquitta/flight-api/internal/provider/cache/rediscache"
	"github.com/danielmesquitta/flight-api/internal/provider/flightapi"
	"github.com/danielmesquitta/flight-api/internal/provider/flightapi/amadeusapi"
	"github.com/danielmesquitta/flight-api/internal/provider/flightapi/mockflightapi"
)

func init() {
	_ = providers
	_ = devProviders
	_ = testProviders
	_ = stagingProviders
	_ = prodProviders
	_ = params
}

func params(
	v validator.Validator,
	e *env.Env,
	t *testing.T,
) {
}

var providers = []any{
	jwtutil.NewJWT,

	wire.Bind(new(cache.Cache), new(*rediscache.RedisCache)),
	rediscache.NewRedisCache,

	flight.NewSearchFlightUseCase,

	handler.NewDocHandler,
	handler.NewHealthHandler,
	handler.NewFlightHandler,

	middleware.NewMiddleware,

	router.NewRouter,

	server.Build,
}

var devProviders = []any{
	amadeusapi.NewAmadeusAPI,
	flightapi.NewFlightAPIs,
}

var testProviders = []any{
	wire.Bind(new(interface {
		mock.TestingT
		Cleanup(func())
	}), new(*testing.T)),
	mockflightapi.NewMockFlightAPI,
	mockflightapi.NewMockFlightAPIs,
}

var stagingProviders = []any{
	amadeusapi.NewAmadeusAPI,
	flightapi.NewFlightAPIs,
}

var prodProviders = []any{
	amadeusapi.NewAmadeusAPI,
	flightapi.NewFlightAPIs,
}
