package postgres

import (
	"embed"
	"errors"
	"fmt"
	"net/http"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	"github.com/golang-migrate/migrate/v4/source/httpfs"
)

//go:embed migrations
var migrations embed.FS

func GetMigrationHandler(dbUrl string) (*migrate.Migrate, error) {
	// use httpFS until go-migrate implements ioFS (see https://github.com/golang-migrate/migrate/issues/480#issuecomment-731518493)
	source, err := httpfs.New(http.FS(migrations), "migrations")
	if err != nil {
		return nil, fmt.Errorf("failed to init httpfs: %w", err)
	}

	mgt, err := migrate.NewWithSourceInstance("httpfs", source, dbUrl)
	if err != nil {
		return nil, fmt.Errorf("failed to create migrate source: %w", err)
	}

	return mgt, nil
}

func RunMigrations(dbUrl string) error {
	m, err := GetMigrationHandler(dbUrl)
	if err != nil {
		return fmt.Errorf("failed to get migration handler: %w", err)
	}

	defer m.Close()

	if err := m.Up(); err != nil && !errors.Is(err, migrate.ErrNoChange) {
		return fmt.Errorf("failed to run migrations: %w", err)
	}

	return nil
}
