package middleware

import (
	"errors"
	"log/slog"
	"net/http"

	"github.com/danielmesquitta/flight-api/internal/app/server/dto"
	"github.com/danielmesquitta/flight-api/internal/app/server/handler"
	"github.com/danielmesquitta/flight-api/internal/domain/errs"
	"github.com/gofiber/fiber/v2"
)

type requestIDContextKey string

const RequestIDContextKey requestIDContextKey = "requestid"

var mapAppErrToHTTPError = map[errs.Code]int{
	errs.ErrCodeForbidden:    http.StatusForbidden,
	errs.ErrCodeUnauthorized: http.StatusUnauthorized,
	errs.ErrCodeValidation:   http.StatusBadRequest,
	errs.ErrCodeUnknown:      http.StatusInternalServerError,
	errs.ErrCodeNotFound:     http.StatusNotFound,
}

func (m *Middleware) ErrorHandler(ctx *fiber.Ctx, err error) error {
	var appErr *errs.Err
	if errors.As(err, &appErr) {
		code, ok := mapAppErrToHTTPError[appErr.Code]
		if !ok {
			code = http.StatusInternalServerError
		}

		if code >= 500 {
			return m.handleInternalServerError(ctx, appErr)
		}

		return ctx.Status(code).JSON(
			dto.ErrorResponse{Message: appErr.Message},
		)
	}

	var fiberErr *fiber.Error
	if errors.As(err, &fiberErr) {
		code := fiberErr.Code

		if code >= 500 {
			return m.handleInternalServerError(ctx, errs.New(fiberErr))
		}

		return ctx.Status(code).JSON(
			dto.ErrorResponse{Message: fiberErr.Message},
		)
	}

	return nil
}

func (m *Middleware) handleInternalServerError(
	c *fiber.Ctx,
	appErr *errs.Err,
) error {
	statusCode := mapAppErrToHTTPError[appErr.Code]

	args := []any{
		"method", c.Method(),
		"url", c.BaseURL(),
	}

	queries := c.Queries()
	if len(queries) > 0 {
		args = append(args, "query", queries)
	}

	requestId := c.Locals(RequestIDContextKey)
	if requestId != "" {
		args = append(args, "request_id", requestId)
	}

	userId := ""
	claims := handler.GetClaims(c)
	if claims != nil {
		userId = claims.Issuer
	}
	if userId != "" {
		args = append(args, "user_id", userId)
	}

	requestData := map[string]any{}
	_ = c.BodyParser(&requestData)
	if len(requestData) > 0 {
		args = append(args, "body", requestData)
	}

	args = append(args, "stacktrace", appErr.StackTrace)

	slog.Error(
		appErr.Error(),
		args...,
	)

	return c.Status(statusCode).JSON(
		dto.ErrorResponse{Message: "internal server error"},
	)
}
