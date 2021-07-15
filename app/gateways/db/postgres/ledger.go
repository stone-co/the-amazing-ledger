package postgres

import (
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/sirupsen/logrus"

	"github.com/stone-co/the-amazing-ledger/app/domain"
	"github.com/stone-co/the-amazing-ledger/app/gateways/db/querybuilder"
)

const (
	collection = "entry"
)

var _ domain.Repository = &LedgerRepository{}

type LedgerRepository struct {
	db  *pgxpool.Pool
	log *logrus.Logger
	qb  querybuilder.QueryBuilder
}

func NewLedgerRepository(db *pgxpool.Pool, log *logrus.Logger) *LedgerRepository {
	qb := querybuilder.New(createTransactionQuery, numArgs)
	qb.Init(numDefaultQueries)

	return &LedgerRepository{
		db:  db,
		log: log,
		qb:  qb,
	}
}
