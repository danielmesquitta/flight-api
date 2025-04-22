package handler

import (
	"github.com/danielmesquitta/flight-api/internal/pkg/jwtutil"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
)

type QueryParam = string

const (
	QueryParamOrigin      QueryParam = "origin"
	QueryParamDestination QueryParam = "destination"
	QueryParamDate        QueryParam = "date"
)

func GetClaims(
	c *fiber.Ctx,
) *jwtutil.UserClaims {
	token, ok := c.Locals(jwtutil.ClaimsKey).(*jwt.Token)
	if !ok {
		return nil
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil
	}

	issuer, _ := claims.GetIssuer()
	issuedAt, _ := claims.GetIssuedAt()
	expiresAt, _ := claims.GetExpirationTime()

	return &jwtutil.UserClaims{
		Issuer:    issuer,
		IssuedAt:  issuedAt.Time,
		ExpiresAt: expiresAt.Time,
	}
}
