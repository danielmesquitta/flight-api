package middleware

import (
	"github.com/danielmesquitta/flight-api/internal/config/env"
	"github.com/danielmesquitta/flight-api/internal/pkg/jwtutil"
)

type Middleware struct {
	e *env.Env
	j *jwtutil.JWT
}

func NewMiddleware(
	e *env.Env,
	j *jwtutil.JWT,
) *Middleware {
	return &Middleware{
		e: e,
		j: j,
	}
}
