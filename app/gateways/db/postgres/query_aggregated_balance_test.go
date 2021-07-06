package postgres

import (
	"context"
	"encoding/json"
	"fmt"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"

	"github.com/stone-co/the-amazing-ledger/app"
	"github.com/stone-co/the-amazing-ledger/app/domain/entities"
	"github.com/stone-co/the-amazing-ledger/app/domain/vos"
)

func TestLedgerRepository_QueryAggregatedBalance(t *testing.T) {
	r := NewLedgerRepository(pgDocker.DB, logrus.New())
	ctx := context.Background()

	metadata := json.RawMessage(`{}`)

	_, err := pgDocker.DB.Exec(ctx, `insert into event (id, name) values (3, 'query_aggregated_balance');`)
	assert.NoError(t, err)

	query, err := vos.NewAccountQuery("liability.agg.*")
	assert.NoError(t, err)

	acc1, err := vos.NewAccountPath("liability.agg.agg1")
	assert.NoError(t, err)

	acc2, err := vos.NewAccountPath("liability.agg.agg2")
	assert.NoError(t, err)

	acc3, err := vos.NewAccountPath("liability.abc.agg3")
	assert.NoError(t, err)

	_, err = r.QueryAggregatedBalance(ctx, query)
	assert.ErrorIs(t, err, app.ErrAccountNotFound)

	e1, _ := entities.NewEntry(
		uuid.New(),
		vos.DebitOperation,
		acc1.Name(),
		vos.NextAccountVersion,
		100,
		metadata,
	)
	e2, _ := entities.NewEntry(
		uuid.New(),
		vos.CreditOperation,
		acc2.Name(),
		vos.IgnoreAccountVersion,
		100,
		metadata,
	)

	tx, err := entities.NewTransaction(uuid.New(), 1, "company", time.Now(), e1, e2)
	assert.NoError(t, err)

	err = r.CreateTransaction(ctx, tx)
	assert.NoError(t, err)

	balance, err := r.QueryAggregatedBalance(ctx, query)
	assert.NoError(t, err)
	assert.Equal(t, 0, balance.Balance)

	_, err = fetchQuerySnapshot(ctx, pgDocker.DB, query)
	assert.ErrorIs(t, err, pgx.ErrNoRows)

	e1, _ = entities.NewEntry(
		uuid.New(),
		vos.DebitOperation,
		acc1.Name(),
		vos.IgnoreAccountVersion,
		100,
		metadata,
	)
	e3, _ := entities.NewEntry(
		uuid.New(),
		vos.CreditOperation,
		acc3.Name(),
		vos.NextAccountVersion,
		100,
		metadata,
	)

	tx, err = entities.NewTransaction(uuid.New(), 1, "company", time.Now(), e1, e3)
	assert.NoError(t, err)
	tx.Event = 1

	err = r.CreateTransaction(ctx, tx)
	assert.NoError(t, err)

	balance, err = r.QueryAggregatedBalance(ctx, query)
	assert.NoError(t, err)
	assert.Equal(t, -100, balance.Balance)

	snap, err := fetchQuerySnapshot(ctx, pgDocker.DB, query)
	assert.NoError(t, err)
	assert.Equal(t, 0, snap.balance)

	e1, _ = entities.NewEntry(
		uuid.New(),
		vos.CreditOperation,
		acc1.Name(),
		vos.NextAccountVersion,
		200,
		metadata,
	)
	e3, _ = entities.NewEntry(
		uuid.New(),
		vos.DebitOperation,
		acc3.Name(),
		vos.NextAccountVersion,
		200,
		metadata,
	)

	tx, err = entities.NewTransaction(uuid.New(), 1, "company", time.Now(), e1, e3)
	assert.NoError(t, err)
	tx.Event = 1

	err = r.CreateTransaction(ctx, tx)
	assert.NoError(t, err)

	balance, err = r.QueryAggregatedBalance(ctx, query)
	assert.NoError(t, err)
	assert.Equal(t, 100, balance.Balance)

	snap, err = fetchQuerySnapshot(ctx, pgDocker.DB, query)
	assert.NoError(t, err)
	assert.Equal(t, -100, snap.balance)
}

type querySnapshot struct {
	balance int
	date    time.Time
}

func fetchQuerySnapshot(ctx context.Context, db *pgxpool.Pool, query vos.AccountQuery) (querySnapshot, error) {
	const cmd = "select balance, tx_date from aggregated_query_balance where query = $1;"

	var snap querySnapshot

	err := db.QueryRow(ctx, cmd, query.Value()).Scan(&snap.balance, &snap.date)
	if err != nil {
		return querySnapshot{}, fmt.Errorf("failed to fetch query snapshot: %w", err)
	}

	return snap, nil
}
