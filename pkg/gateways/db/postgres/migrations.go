package postgres

//go:generate go-bindata -prefix migrations/ -o migrations.gen.go -pkg postgres migrations

import (
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	bindata "github.com/golang-migrate/migrate/v4/source/go_bindata"
)

func GetMigrationHandler(dbUrl string) (*migrate.Migrate, error) {
	s := bindata.Resource(AssetNames(),
		func(name string) ([]byte, error) {
			return Asset(name)
		})

	d, err := bindata.WithInstance(s)
	if err != nil {
		return nil, err
	}
	return migrate.NewWithSourceInstance("go-bindata", d, dbUrl)
}

func RunMigrations(dbUrl string) error {
	m, err := GetMigrationHandler(dbUrl)
	if err != nil {
		return err
	}
	defer m.Close()
	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		return err
	}
	return nil
}
