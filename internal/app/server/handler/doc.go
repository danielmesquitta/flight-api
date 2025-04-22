package handler

import (
	"net/http"
	"strings"

	root "github.com/danielmesquitta/flight-api"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/filesystem"
	fiberSwagger "github.com/swaggo/fiber-swagger"
)

type DocHandler struct{}

func NewDocHandler() *DocHandler {
	return &DocHandler{}
}

func (d *DocHandler) Get(c *fiber.Ctx) error {
	h := fiberSwagger.WrapHandler

	path := strings.Trim(c.Path(), "/")
	split := strings.Split(path, "/")
	lastPathParam := split[len(split)-1]

	files := map[string]struct{}{
		"openapi.yaml": {},
		"openapi.json": {},
	}

	if _, ok := files[lastPathParam]; ok {
		h = filesystem.New(filesystem.Config{
			Root:       http.FS(root.DocFiles),
			PathPrefix: "docs",
			Browse:     true,
		})
	}

	return h(c)
}
