package postgres

import (
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/sirupsen/logrus"

	"github.com/stone-co/the-amazing-ledger/app/domain"
)

const (
	collection = "entry"
)

var _ domain.Repository = &LedgerRepository{}

type LedgerRepository struct {
	db  *pgxpool.Pool
	log *logrus.Logger
}

func NewLedgerRepository(db *pgxpool.Pool, log *logrus.Logger) *LedgerRepository {
	initCreateTransactioQueryMap()

	return &LedgerRepository{
		db:  db,
		log: log,
	}
}

func initCreateTransactioQueryMap() {
	createTransactionQueryMap = make(map[int]string)

	for i := 2; i < maxQueriesDefault; i++ {
		createTransactionQueryMap[i] = buildQuery(i)
	}
}
