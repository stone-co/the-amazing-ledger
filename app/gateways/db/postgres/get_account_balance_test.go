package postgres

import (
	"context"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"

	"github.com/stone-co/the-amazing-ledger/app/domain/entities"
	"github.com/stone-co/the-amazing-ledger/app/domain/vos"
)

func TestLedgerRepository_GetAccountBalance(t *testing.T) {
	r := NewLedgerRepository(pgDocker.DB, logrus.New())
	ctx := context.Background()

	_, err := pgDocker.DB.Exec(ctx, `insert into event (name) values ('defaults');`)
	assert.NoError(t, err)

	acc1, err := vos.NewAccountPath("liability.123.account1")
	assert.NoError(t, err)

	acc2, err := vos.NewAccountPath("liability.123.account2")
	assert.NoError(t, err)

	e1, _ := entities.NewEntry(
		uuid.New(),
		vos.DebitOperation,
		acc1.Name(),
		vos.NextAccountVersion,
		100,
	)
	e2, _ := entities.NewEntry(
		uuid.New(),
		vos.CreditOperation,
		acc2.Name(),
		vos.IgnoreAccountVersion,
		100,
	)

	tx, err := entities.NewTransaction(uuid.New(), e1, e2)
	assert.NoError(t, err)
	tx.Event = 1

	err = r.CreateTransaction(ctx, tx)
	assert.NoError(t, err)

	balance, err := r.GetAccountBalance(ctx, acc1)
	assert.NoError(t, err)
	assert.Equal(t, 0, balance.TotalCredit)
	assert.Equal(t, 100, balance.TotalDebit)

	snap, err := fetchSnapshot(ctx, pgDocker.DB, acc1)
	assert.NoError(t, err)
	assert.Equal(t, 0, snap.credit)
	assert.Equal(t, 100, snap.debit)
	assert.Equal(t, 1, snap.version)

	balance, err = r.GetAccountBalance(ctx, acc2)
	assert.NoError(t, err)
	assert.Equal(t, 100, balance.TotalCredit)
	assert.Equal(t, 0, balance.TotalDebit)

	_, err = fetchSnapshot(ctx, pgDocker.DB, acc2)
	assert.ErrorIs(t, pgx.ErrNoRows, err)

	e1, _ = entities.NewEntry(
		uuid.New(),
		vos.DebitOperation,
		acc1.Name(),
		vos.IgnoreAccountVersion,
		100,
	)
	e2, _ = entities.NewEntry(
		uuid.New(),
		vos.CreditOperation,
		acc2.Name(),
		vos.NextAccountVersion,
		100,
	)

	tx, err = entities.NewTransaction(uuid.New(), e1, e2)
	assert.NoError(t, err)
	tx.Event = 1

	err = r.CreateTransaction(ctx, tx)
	assert.NoError(t, err)

	balance, err = r.GetAccountBalance(ctx, acc1)
	assert.NoError(t, err)
	assert.Equal(t, 0, balance.TotalCredit)
	assert.Equal(t, 200, balance.TotalDebit)

	snap, err = fetchSnapshot(ctx, pgDocker.DB, acc1)
	assert.NoError(t, err)
	assert.Equal(t, 0, snap.credit)
	assert.Equal(t, 100, snap.debit)
	assert.Equal(t, 1, snap.version)

	balance, err = r.GetAccountBalance(ctx, acc2)
	assert.NoError(t, err)
	assert.Equal(t, 200, balance.TotalCredit)
	assert.Equal(t, 0, balance.TotalDebit)

	snap, err = fetchSnapshot(ctx, pgDocker.DB, acc2)
	assert.NoError(t, err)
	assert.Equal(t, 200, snap.credit)
	assert.Equal(t, 0, snap.debit)
	assert.Equal(t, 1, snap.version)

	e1, _ = entities.NewEntry(
		uuid.New(),
		vos.DebitOperation,
		acc1.Name(),
		vos.NextAccountVersion,
		100,
	)
	e2, _ = entities.NewEntry(
		uuid.New(),
		vos.CreditOperation,
		acc2.Name(),
		vos.NextAccountVersion,
		100,
	)

	tx, err = entities.NewTransaction(uuid.New(), e1, e2)
	assert.NoError(t, err)
	tx.Event = 1

	err = r.CreateTransaction(ctx, tx)
	assert.NoError(t, err)

	balance, err = r.GetAccountBalance(ctx, acc1)
	assert.NoError(t, err)
	assert.Equal(t, 0, balance.TotalCredit)
	assert.Equal(t, 300, balance.TotalDebit)

	snap, err = fetchSnapshot(ctx, pgDocker.DB, acc1)
	assert.NoError(t, err)
	assert.Equal(t, 0, snap.credit)
	assert.Equal(t, 300, snap.debit)
	assert.Equal(t, 2, snap.version)

	balance, err = r.GetAccountBalance(ctx, acc2)
	assert.NoError(t, err)
	assert.Equal(t, 300, balance.TotalCredit)
	assert.Equal(t, 0, balance.TotalDebit)

	snap, err = fetchSnapshot(ctx, pgDocker.DB, acc2)
	assert.NoError(t, err)
	assert.Equal(t, 300, snap.credit)
	assert.Equal(t, 0, snap.debit)
	assert.Equal(t, 2, snap.version)
}

type snapshot struct {
	credit  int
	debit   int
	version int
	date    time.Time
}

func fetchSnapshot(ctx context.Context, db *pgxpool.Pool, account vos.AccountPath) (snapshot, error) {
	const query = "select credit, debit, tx_date, tx_version from account_balance where account = $1;"

	var snap snapshot

	err := db.QueryRow(ctx, query, account.Name()).Scan(
		&snap.credit,
		&snap.debit,
		&snap.date,
		&snap.version,
	)

	return snap, err
}
