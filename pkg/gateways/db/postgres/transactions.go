package postgres

import (
	"context"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/sirupsen/logrus"
	"github.com/stone-co/the-amazing-ledger/pkg/command-handler/domain/transactions"
	"github.com/stone-co/the-amazing-ledger/pkg/command-handler/domain/transactions/entities"
)

var _ transactions.Repository = &TransactionsRepository{}

type TransactionsRepository struct {
	db  *pgxpool.Pool
	log *logrus.Logger
}

func NewTransactionsRepository(db *pgxpool.Pool, log *logrus.Logger) *TransactionsRepository {
	return &TransactionsRepository{
		db:  db,
		log: log,
	}
}

func (r *TransactionsRepository) Create(o *[]entities.Transaction) error {
	var err error = nil
	operationId := uuid.New().String()

	for _, t := range *o {
		transactionId := uuid.New().String()
		if err = r.db.QueryRow(context.Background(), `INSERT INTO
		transactions (
			id,
			account_id,
			operation_id,
			request_id,
			amount,
			balance_after
		) VALUES ($1, $2, $3, $4, $5, $6)
		returning created_at`,
			transactionId,
			t.AccountID,
			operationId,
			t.RequestID,
			t.Amount,
			t.BalanceAfter,
		).Scan(&t.CreatedAt); err != nil {
			break
		}
	}

	return err
}
