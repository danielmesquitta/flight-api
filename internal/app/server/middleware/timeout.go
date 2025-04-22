package middleware

import (
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/timeout"
)

func (m *Middleware) Timeout(t time.Duration) fiber.Handler {
	return timeout.NewWithContext(
		func(c *fiber.Ctx) (err error) { return c.Next() },
		t,
	)
}
