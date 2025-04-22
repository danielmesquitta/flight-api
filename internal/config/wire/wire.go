package wire

import (
	"testing"

	"github.com/google/wire"

	"github.com/danielmesquitta/flight-api/internal/app/server"
	"github.com/danielmesquitta/flight-api/internal/app/server/handler"
	"github.com/danielmesquitta/flight-api/internal/app/server/middleware"
	"github.com/danielmesquitta/flight-api/internal/app/server/router"
	"github.com/danielmesquitta/flight-api/internal/config/env"
	"github.com/danielmesquitta/flight-api/internal/pkg/jwtutil"
	"github.com/danielmesquitta/flight-api/internal/pkg/validator"
	"github.com/danielmesquitta/flight-api/internal/provider/cache"
	"github.com/danielmesquitta/flight-api/internal/provider/cache/rediscache"
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

	handler.NewDocHandler,
	handler.NewHealthHandler,

	middleware.NewMiddleware,

	router.NewRouter,

	server.Build,
}

var devProviders = []any{
	// Specific dev providers
}

var testProviders = []any{
	// Specific test providers
}

var stagingProviders = []any{
	// Specific staging providers
}

var prodProviders = []any{
	// Specific prod providers
}
