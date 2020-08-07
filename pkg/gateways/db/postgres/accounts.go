package postgres

import (
	"context"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/sirupsen/logrus"
	"github.com/stone-co/the-amazing-ledger/pkg/command-handler/domain/accounts"
	"github.com/stone-co/the-amazing-ledger/pkg/command-handler/domain/accounts/entities"
)

var _ accounts.Repository = &AccountsRepository{}

type AccountsRepository struct {
	db  *pgxpool.Pool
	log *logrus.Logger
}

func NewAccountsRepository(db *pgxpool.Pool, log *logrus.Logger) *AccountsRepository {
	return &AccountsRepository{
		db:  db,
		log: log,
	}
}

func (r *AccountsRepository) Create(a *entities.Account) error {
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

func (r *AccountsRepository) Get(id *string) (entities.Account, error) {
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

func (r *AccountsRepository) Update(id *string, balance *int) error {
	if _, err := r.db.Exec(context.Background(),
		`UPDATE accounts set balance = $1 where id = $2`, balance, id); err != nil {
		return err
	}

	return nil
}
