package postgres

import (
	"embed"
	"errors"
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
		return nil, err
	}

	return migrate.NewWithSourceInstance("httpfs", source, dbUrl)
}

func RunMigrations(dbUrl string) error {
	m, err := GetMigrationHandler(dbUrl)
	if err != nil {
		return err
	}

	defer m.Close()
	if err := m.Up(); err != nil && !errors.Is(err, migrate.ErrNoChange) {
		return err
	}
	return nil
}
