package migrations

import "embed"

// FS bundles all .sql files in the migrations directory.
//go:embed *.sql
var FS embed.FS
