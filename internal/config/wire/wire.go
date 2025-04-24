package wire

import (
	"testing"

	"github.com/google/wire"

	"github.com/danielmesquitta/flight-api/internal/app/server"
	"github.com/danielmesquitta/flight-api/internal/app/server/handler"
	"github.com/danielmesquitta/flight-api/internal/app/server/middleware"
	"github.com/danielmesquitta/flight-api/internal/app/server/router"
	"github.com/danielmesquitta/flight-api/internal/config/env"
	"github.com/danielmesquitta/flight-api/internal/domain/usecase/auth"
	"github.com/danielmesquitta/flight-api/internal/domain/usecase/flight"
	"github.com/danielmesquitta/flight-api/internal/pkg/jwtutil"
	"github.com/danielmesquitta/flight-api/internal/pkg/validator"
	"github.com/danielmesquitta/flight-api/internal/provider/cache"
	"github.com/danielmesquitta/flight-api/internal/provider/cache/rediscache"
	"github.com/danielmesquitta/flight-api/internal/provider/flightapi"
	"github.com/danielmesquitta/flight-api/internal/provider/flightapi/amadeusapi"
	"github.com/danielmesquitta/flight-api/internal/provider/flightapi/duffelapi"
	"github.com/danielmesquitta/flight-api/internal/provider/flightapi/serpapi"
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

	amadeusapi.NewAmadeusAPI,
	serpapi.NewSerpAPI,
	duffelapi.NewDuffelAPI,
	flightapi.NewFlightAPIs,

	wire.Bind(new(cache.Cache), new(*rediscache.RedisCache)),
	rediscache.NewRedisCache,

	flight.NewSearchFlightsUseCase,
	auth.NewLoginUseCase,

	handler.NewDocHandler,
	handler.NewHealthHandler,
	handler.NewFlightHandler,
	handler.NewAuthHandler,

	middleware.NewMiddleware,

	router.NewRouter,

	server.Build,
}

var devProviders = []any{
	// Add any development-specific providers here
}

var testProviders = []any{
	// Add any test-specific providers here
}

var stagingProviders = []any{
	// Add any staging-specific providers here
}

var prodProviders = []any{
	// Add any production-specific providers here
}
