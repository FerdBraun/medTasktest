package postgres

import (
	"context"
	"fmt"
	"log/slog"
	"sort"
	"strings"

	"example.com/taskservice/migrations"
	"github.com/jackc/pgx/v5/pgxpool"
)

// Migrate runs all .sql files in the migrations directory in alphabetical order.
// The SQL scripts should be idempotent (using IF NOT EXISTS).
func Migrate(ctx context.Context, pool *pgxpool.Pool, logger *slog.Logger) error {
	entries, err := migrations.FS.ReadDir(".")
	if err != nil {
		return fmt.Errorf("read migrations dir: %w", err)
	}

	var files []string
	for _, entry := range entries {
		if !entry.IsDir() && strings.HasSuffix(entry.Name(), ".sql") {
			files = append(files, entry.Name())
		}
	}
	sort.Strings(files)

	for _, file := range files {
		logger.Info("applying migration", "file", file)
		content, err := migrations.FS.ReadFile(file)
		if err != nil {
			return fmt.Errorf("read migration file %s: %w", file, err)
		}

		if _, err := pool.Exec(ctx, string(content)); err != nil {
			return fmt.Errorf("execute migration %s: %w", file, err)
		}
	}

	return nil
}
