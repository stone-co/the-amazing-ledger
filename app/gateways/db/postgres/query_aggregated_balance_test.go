package postgres

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"

	"github.com/stone-co/the-amazing-ledger/app"
	"github.com/stone-co/the-amazing-ledger/app/domain/vos"
	"github.com/stone-co/the-amazing-ledger/app/tests"
)

func TestLedgerRepository_QueryAggregatedBalanceFailure(t *testing.T) {
	t.Run("should return an error if accounts do not exist", func(t *testing.T) {
		r := NewLedgerRepository(pgDocker.DB, logrus.New())
		ctx := context.Background()

		query, err := vos.NewAccountQuery("liability.agg.*")
		assert.NoError(t, err)

		_, err = r.QueryAggregatedBalance(ctx, query)
		assert.ErrorIs(t, err, app.ErrAccountNotFound)
	})
}

func TestLedgerRepository_QueryAggregatedBalanceSuccess(t *testing.T) {
	acc1, err := vos.NewAccountPath("liability.agg.agg1")
	assert.NoError(t, err)

	acc2, err := vos.NewAccountPath("liability.agg.agg2")
	assert.NoError(t, err)

	acc3, err := vos.NewAccountPath("liability.abc.agg3")
	assert.NoError(t, err)

	query, err := vos.NewAccountQuery("liability.agg.*")
	assert.NoError(t, err)

	type args struct {
		acc1   vos.AccountPath
		acc2   vos.AccountPath
		debit  int
		credit int
	}

	type wants struct {
		balance     int
		snapErr     error
		snapBalance int
	}

	testCases := []struct {
		name     string
		repoSeed func(t *testing.T, ctx context.Context, r *LedgerRepository)
		args     args
		wants    wants
	}{
		{
			name:     "should query aggregated balance involving two accounts",
			repoSeed: func(t *testing.T, ctx context.Context, r *LedgerRepository) {},
			args: args{
				acc1:   acc1,
				acc2:   acc2,
				debit:  100,
				credit: 100,
			},
			wants: wants{
				balance: 0,
				snapErr: pgx.ErrNoRows,
			},
		},
		{
			name: "should query aggregated balance involving three accounts (first snapshot)",
			repoSeed: func(t *testing.T, ctx context.Context, r *LedgerRepository) {
				e1 := createEntry(t, vos.DebitOperation, acc1.Name(), vos.NextAccountVersion, 100)
				e2 := createEntry(t, vos.CreditOperation, acc2.Name(), vos.IgnoreAccountVersion, 100)

				createTransaction(t, ctx, r, e1, e2)

				_, err = r.QueryAggregatedBalance(ctx, query)
				assert.NoError(t, err)
			},
			args: args{
				acc1:   acc1,
				acc2:   acc3,
				debit:  100,
				credit: 100,
			},
			wants: wants{
				balance:     -100,
				snapBalance: 0,
			},
		},
		{
			name: "should query aggregated balance involving three accounts (second snapshot)",
			repoSeed: func(t *testing.T, ctx context.Context, r *LedgerRepository) {
				e1 := createEntry(t, vos.DebitOperation, acc1.Name(), vos.NextAccountVersion, 100)
				e2 := createEntry(t, vos.CreditOperation, acc2.Name(), vos.IgnoreAccountVersion, 100)

				createTransaction(t, ctx, r, e1, e2)

				_, err = r.QueryAggregatedBalance(ctx, query)
				assert.NoError(t, err)

				e1 = createEntry(t, vos.DebitOperation, acc1.Name(), vos.IgnoreAccountVersion, 100)
				e3 := createEntry(t, vos.CreditOperation, acc3.Name(), vos.NextAccountVersion, 100)

				createTransaction(t, ctx, r, e1, e3)

				_, err = r.QueryAggregatedBalance(ctx, query)
				assert.NoError(t, err)
			},
			args: args{
				acc1:   acc1,
				acc2:   acc3,
				debit:  200,
				credit: 200,
			},
			wants: wants{
				balance:     -300,
				snapBalance: -100,
			},
		},
	}

	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.Background()
			r := NewLedgerRepository(pgDocker.DB, logrus.New())

			tt.repoSeed(t, ctx, r)

			defer tests.TruncateTables(ctx, pgDocker.DB, "entry", "account_version", "aggregated_query_balance")

			e1 := createEntry(t, vos.DebitOperation, tt.args.acc1.Name(), vos.NextAccountVersion, tt.args.debit)
			e2 := createEntry(t, vos.CreditOperation, tt.args.acc2.Name(), vos.IgnoreAccountVersion, tt.args.credit)

			createTransaction(t, ctx, r, e1, e2)

			balance, err := r.QueryAggregatedBalance(ctx, query)
			assert.NoError(t, err)
			assert.Equal(t, tt.wants.balance, balance.Balance)

			if tt.wants.snapErr != nil {
				_, err = fetchQuerySnapshot(ctx, pgDocker.DB, query)
				assert.ErrorIs(t, err, pgx.ErrNoRows)
			} else {
				snap, err := fetchQuerySnapshot(ctx, pgDocker.DB, query)
				assert.NoError(t, err)
				assert.Equal(t, tt.wants.snapBalance, snap.balance)
			}
		})
	}
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
