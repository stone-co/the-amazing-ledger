package postgres

import (
	"context"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/sirupsen/logrus"
	"github.com/stone-co/the-amazing-ledger/pkg/command-handler/domain/entries"
	"github.com/stone-co/the-amazing-ledger/pkg/command-handler/domain/entries/entities"
)

var _ entries.Repository = &EntriesRepository{}

type EntriesRepository struct {
	db  *pgxpool.Pool
	log *logrus.Logger
}

func NewEntriesRepository(db *pgxpool.Pool, log *logrus.Logger) *EntriesRepository {
	return &EntriesRepository{
		db:  db,
		log: log,
	}
}

func (r *EntriesRepository) Create(o *[]entities.Entry) error {
	var err error = nil
	transactionId := uuid.New().String()

	for _, e := range *o {
		entryId := uuid.New().String()
		if err = r.db.QueryRow(context.Background(), `INSERT INTO
		entries (
			id,
			account_id,
			transaction_id,
			request_id,
			amount,
			balance_after
		) VALUES ($1, $2, $3, $4, $5, $6)
		returning created_at`,
			entryId,
			e.AccountID,
			transactionId,
			e.RequestID,
			e.Amount,
			e.BalanceAfter,
		).Scan(&e.CreatedAt); err != nil {
			break
		}
	}

	return err
}
