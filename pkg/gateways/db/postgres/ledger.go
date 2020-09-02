package postgres

import (
	"context"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/sirupsen/logrus"
	"github.com/stone-co/the-amazing-ledger/pkg/command-handler/domain/ledger"
	"github.com/stone-co/the-amazing-ledger/pkg/command-handler/domain/ledger/entities"
)

var _ ledger.Repository = &LedgerRepository{}

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

func (r *LedgerRepository) CreateAccount(a *entities.Account) error {
	a.ID = uuid.New().String()
	if err := r.db.QueryRow(context.Background(), `INSERT INTO
		accounts (
			id,
			owner,
			name,
			owner_id,
			type,
			metadata,
			balance
		) VALUES ($1, $2, $3, $4, $5, $6, $7)
		returning created_at`,
		a.ID,
		a.Owner,
		a.Name,
		a.OwnerID,
		a.Type,
		a.Metadata,
		a.Balance,
	).Scan(&a.CreatedAt); err != nil {
		return err
	}

	return nil
}

func (r *LedgerRepository) GetAccount(id string) (entities.Account, error) {
	var account = entities.Account{}
	row := r.db.QueryRow(context.Background(),
		`SELECT
			id, owner, name, owner_id, type, metadata, balance
		FROM accounts where id = $1`, id)

	err := row.Scan(
		&account.ID,
		&account.Owner,
		&account.Name,
		&account.OwnerID,
		&account.Type,
		&account.Metadata,
		&account.Balance,
	)
	return account, err
}

func (r *LedgerRepository) SearchAccount(accountType string, accountOwnerID string, accountOwner string, accountName string, accountMetadata string) (entities.Account, error) {
	var account = entities.Account{}
	row := r.db.QueryRow(context.Background(),
		`select
			id, owner, name, owner_id, type, metadata, balance
		from accounts 
		where 
		accountType = $1,
		and accountOwnerID = $2,
		and accountOwner = $3,
		and accountName = $4,
		and accountMetadata = $5`,
		accountType,
		accountOwnerID,
		accountOwner,
		accountName,
		accountMetadata)

	err := row.Scan(
		&account.ID,
		&account.Owner,
		&account.Name,
		&account.OwnerID,
		&account.Type,
		&account.Metadata,
		&account.Balance,
	)
	return account, err
}

func (r *LedgerRepository) UpdateBalance(id string, balance int) error {
	if _, err := r.db.Exec(context.Background(),
		`UPDATE accounts set balance = $1 where id = $2`, balance, id); err != nil {
		return err
	}

	return nil
}

func (r *LedgerRepository) CreateTransaction(o *[]entities.Entry) error {
	transactionId := uuid.New().String()

	tx, err := r.db.Begin(context.Background())
	if err != nil {
		return err
	}

	for _, e := range *o {
		entryId := uuid.New().String()
		if err = tx.QueryRow(context.Background(), `INSERT INTO
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
			return err
		}
	}

	err = tx.Commit(context.Background())
	if err != nil {
		return err
	}

	return err
}
