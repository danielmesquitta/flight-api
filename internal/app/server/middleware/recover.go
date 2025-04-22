package middleware

import (
	"errors"
	"fmt"

	"github.com/danielmesquitta/flight-api/internal/domain/errs"
	"github.com/gofiber/fiber/v2"
)

func (m *Middleware) Recover() fiber.Handler {
	return func(c *fiber.Ctx) error {
		defer func() {
			if r := recover(); r != nil {
				var err error
				switch x := r.(type) {
				case string:
					err = errors.New(x)
				case error:
					err = x
				default:
					err = fmt.Errorf("%v", x)
				}

				appErr := errs.New(err)

				_ = m.handleInternalServerError(c, appErr)
			}
		}()

		return c.Next()
	}
}
