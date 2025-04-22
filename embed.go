package root

import (
	"embed"
)

//go:embed .env*
var Env embed.FS

//go:embed docs/openapi.yaml docs/openapi.json
var DocFiles embed.FS
