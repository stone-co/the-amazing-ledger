package postgres

import (
	"context"

	"github.com/jackc/pgx/v4/log/logrusadapter"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/sirupsen/logrus"
)

func ConnectPool(dbURL string, log *logrus.Logger) (*pgxpool.Pool, error) {
	config, err := pgxpool.ParseConfig(dbURL)
	if err != nil {
		return nil, err
	}
	if log != nil {
		config.ConnConfig.Logger = logrusadapter.NewLogger(log)
	}
	db, err := pgxpool.ConnectConfig(context.Background(), config)
	return db, err
}
