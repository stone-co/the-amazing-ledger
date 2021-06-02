package tests

import (
	"context"
	"fmt"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/ory/dockertest"
)

type PostgresDocker struct {
	DB       *pgxpool.Pool
	Pool     *dockertest.Pool
	Resource *dockertest.Resource
}

func SetupTest(migrationsPath string) *PostgresDocker {
	var conn *pgxpool.Pool

	pool, err := dockertest.NewPool("")
	if err != nil {
		panic(fmt.Errorf("failed to connect to docker: %w", err))
	}

	database := "dev"

	resource, err := pool.Run(
		"postgres",
		"13.2",
		[]string{"POSTGRES_PASSWORD=postgres", "POSTGRES_DB=" + database},
	)
	if err != nil {
		panic(fmt.Errorf("failed to start resource: %w", err))
	}

	connString := fmt.Sprintf(
		"postgres://postgres:postgres@localhost:%s/%s?sslmode=disable",
		resource.GetPort("5432/tcp"),
		database)

	if err = pool.Retry(func() error {
		var rErr error
		ctx := context.Background()
		conn, rErr = pgxpool.Connect(ctx, connString)
		if rErr != nil {
			return rErr
		}
		_, rErr = conn.Acquire(ctx)
		if rErr != nil {
			return rErr
		}

		return nil
	}); err != nil {
		panic(fmt.Errorf("failed to connect to docker: %w", err))
	}

	if err := runMigrations(migrationsPath, connString); err != nil {
		panic(fmt.Errorf("failed to run migrations: %w", err))
	}

	return &PostgresDocker{
		DB:       conn,
		Pool:     pool,
		Resource: resource,
	}
}

func RemoveContainer(pgDocker *PostgresDocker) {
	if err := pgDocker.Pool.Purge(pgDocker.Resource); err != nil {
		panic(fmt.Errorf("failed to remove container: %w", err))
	}
}

func TruncateTables(ctx context.Context, db *pgx.Conn, tables ...string) {
	if _, err := db.Exec(ctx, "truncate entry, event, account_version"); err != nil {
		panic(fmt.Errorf("failed to truncate table(s) %v: %w", tables, err))
	}
}

func runMigrations(migrationsPath, connString string) error {
	if migrationsPath != "" {
		mig, err := migrate.New("file://"+migrationsPath, connString)
		if err != nil {
			return fmt.Errorf("failed to start migrate struct: %s", err.Error())
		}
		defer mig.Close()
		if err = mig.Up(); err != nil {
			return fmt.Errorf("failed to run migration: %s", err.Error())
		}
	}

	return nil
}
