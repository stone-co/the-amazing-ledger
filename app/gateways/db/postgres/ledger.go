package postgres

import (
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/sirupsen/logrus"
	"github.com/stone-co/the-amazing-ledger/app/domain"
)

const (
	collection = "entries"
)

var _ domain.Repository = &LedgerRepository{}

type LedgerRepository struct {
	db  *pgxpool.Pool
	log *logrus.Logger
}

func NewLedgerRepository(db *pgxpool.Pool, log *logrus.Logger) *LedgerRepository {
	return &LedgerRepository{
		db:  db,
		log: log,
	}
}
