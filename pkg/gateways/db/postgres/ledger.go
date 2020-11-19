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

func (r *LedgerRepository) CreateAccount(a *entities.Account) (entities.Account, error) {
	a.ID = uuid.New().String()
	row := r.db.QueryRow(context.Background(), `INSERT INTO
		accounts (
			id,
			owner,
			name,
			owner_id,
			type,
			metadata,
			balance
		) VALUES ($1, $2, $3, $4, $5, $6, $7)
		returning
		id,
		type,
		owner_id,
		owner,
		name,
		metadata,
		balance,
		created_at,
		updated_at`,
		a.ID,
		a.Owner,
		a.Name,
		a.OwnerID,
		a.Type,
		a.Metadata,
		a.Balance,
	)

	if err := row.Scan(
		&a.ID,
		&a.Type,
		&a.OwnerID,
		&a.Owner,
		&a.Name,
		&a.Metadata,
		&a.Balance,
		&a.CreatedAt,
		&a.UpdatedAt,
	); err != nil {
		return entities.Account{}, err
	}

	return *a, nil
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
