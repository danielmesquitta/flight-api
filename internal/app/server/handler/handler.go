package handler

import (
	"time"

	"github.com/danielmesquitta/flight-api/internal/domain/errs"
	"github.com/danielmesquitta/flight-api/internal/pkg/jwtutil"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
)

type QueryParam = string

const (
	QueryParamOrigin      QueryParam = "origin"
	QueryParamDestination QueryParam = "destination"
	QueryParamStartDate   QueryParam = "start_date"
	QueryParamEndDate     QueryParam = "end_date"
)

type PathParam = string

const (
	pathParamCategoryID    PathParam = "category_id"
	pathParamTransactionID PathParam = "transaction_id"
	pathParamAIChatID      PathParam = "ai_chat_id"
)

func parseDateQueryParam(
	c *fiber.Ctx,
	param QueryParam,
) (time.Time, error) {
	paramValue := c.Query(param)
	if paramValue == "" {
		return time.Time{}, nil
	}

	date, err := time.Parse(time.RFC3339, paramValue)
	if err != nil {
		return time.Time{}, errs.ErrInvalidDate
	}

	return date, nil
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
