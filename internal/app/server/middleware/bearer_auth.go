package middleware

import (
	"github.com/danielmesquitta/flight-api/internal/pkg/jwtutil"
	jwtware "github.com/gofiber/contrib/jwt"
	"github.com/gofiber/fiber/v2"
)

func (m *Middleware) BearerAuthAccessToken() fiber.Handler {
	return jwtware.New(jwtware.Config{
		ContextKey: jwtutil.ClaimsKey,
		SigningKey: jwtware.SigningKey{
			Key: []byte(m.e.JWTAccessTokenSecretKey),
		},
	})
}
