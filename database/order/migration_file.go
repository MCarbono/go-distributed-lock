package order

import (
	"embed"
)

//go:embed migrations/*.sql
var MigrationsFS embed.FS
