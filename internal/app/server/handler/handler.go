package handler

import (
	"time"

	"github.com/itlightning/dateparse"

	"github.com/danielmesquitta/flight-api/internal/domain/errs"
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

func parseDateQueryParam(
	c *fiber.Ctx,
	param QueryParam,
) (time.Time, error) {
	date := c.Query(param)
	parsedDate, err := dateparse.ParseAny(date)
	if err != nil {
		return time.Time{}, errs.ErrInvalidDateFormat
	}
	return parsedDate, nil
}

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
