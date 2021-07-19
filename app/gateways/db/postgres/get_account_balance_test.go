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
	"github.com/stone-co/the-amazing-ledger/app/tests/testdata"
)

func TestLedgerRepository_GetAccountBalanceSuccess(t *testing.T) {
	acc1, err := vos.NewSingleAccount(testdata.GenerateAccountPath())
	assert.NoError(t, err)

	acc2, err := vos.NewSingleAccount(testdata.GenerateAccountPath())
	assert.NoError(t, err)

	type accountValues struct {
		acc1Credit int
		acc1Debit  int
		acc2Credit int
		acc2Debit  int
	}

	type wants struct {
		total    accountValues
		snapshot accountValues
		err      error
	}

	testCases := []struct {
		name     string
		repoSeed func(t *testing.T, ctx context.Context, r *LedgerRepository)
		wants    wants
	}{
		{
			name:     "should get account balance successfully when is the first request",
			repoSeed: func(t *testing.T, ctx context.Context, r *LedgerRepository) {},
			wants: wants{
				total: accountValues{
					acc1Credit: 0,
					acc1Debit:  100,
					acc2Credit: 100,
					acc2Debit:  0,
				},
				err: pgx.ErrNoRows,
			},
		},
		{
			name: "should get account balance successfully when is the second request",
			repoSeed: func(t *testing.T, ctx context.Context, r *LedgerRepository) {
				e1 := createEntry(t, vos.DebitOperation, acc1.Value(), vos.NextAccountVersion, 100)
				e2 := createEntry(t, vos.CreditOperation, acc2.Value(), vos.NextAccountVersion, 100)

				createTransaction(t, ctx, r, e1, e2)

				_, err = r.GetAccountBalance(ctx, acc1)
				assert.NoError(t, err)

				_, err = r.GetAccountBalance(ctx, acc2)
				assert.NoError(t, err)
			},
			wants: wants{
				total: accountValues{
					acc1Credit: 0,
					acc1Debit:  200,
					acc2Credit: 200,
					acc2Debit:  0,
				},
				snapshot: accountValues{
					acc1Credit: 0,
					acc1Debit:  100,
					acc2Credit: 100,
					acc2Debit:  0,
				},
			},
		},
		{
			name: "should get account balance successfully when is the third request",
			repoSeed: func(t *testing.T, ctx context.Context, r *LedgerRepository) {
				e1 := createEntry(t, vos.DebitOperation, acc1.Value(), vos.NextAccountVersion, 100)
				e2 := createEntry(t, vos.CreditOperation, acc2.Value(), vos.NextAccountVersion, 100)

				createTransaction(t, ctx, r, e1, e2)

				_, err = r.GetAccountBalance(ctx, acc1)
				assert.NoError(t, err)

				_, err = r.GetAccountBalance(ctx, acc2)
				assert.NoError(t, err)

				e3 := createEntry(t, vos.DebitOperation, acc1.Value(), vos.NextAccountVersion, 100)
				e4 := createEntry(t, vos.CreditOperation, acc2.Value(), vos.NextAccountVersion, 100)

				createTransaction(t, ctx, r, e3, e4)

				_, err = r.GetAccountBalance(ctx, acc1)
				assert.NoError(t, err)

				_, err = r.GetAccountBalance(ctx, acc2)
				assert.NoError(t, err)
			},
			wants: wants{
				total: accountValues{
					acc1Credit: 0,
					acc1Debit:  300,
					acc2Credit: 300,
					acc2Debit:  0,
				},
				snapshot: accountValues{
					acc1Credit: 0,
					acc1Debit:  200,
					acc2Credit: 200,
					acc2Debit:  0,
				},
			},
		},
	}

	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.Background()
			r := NewLedgerRepository(pgDocker.DB, logrus.New())

			tt.repoSeed(t, ctx, r)

			defer tests.TruncateTables(ctx, pgDocker.DB, "entry", "account_version", "account_balance")

			e1 := createEntry(t, vos.DebitOperation, acc1.Value(), vos.NextAccountVersion, 100)
			e2 := createEntry(t, vos.CreditOperation, acc2.Value(), vos.NextAccountVersion, 100)

			createTransaction(t, ctx, r, e1, e2)

			balance, err := r.GetAccountBalance(ctx, acc1)
			assert.NoError(t, err)
			assert.Equal(t, tt.wants.total.acc1Credit, balance.TotalCredit)
			assert.Equal(t, tt.wants.total.acc1Debit, balance.TotalDebit)

			balance, err = r.GetAccountBalance(ctx, acc2)
			assert.NoError(t, err)
			assert.Equal(t, tt.wants.total.acc2Credit, balance.TotalCredit)
			assert.Equal(t, tt.wants.total.acc2Debit, balance.TotalDebit)

			if tt.wants.err != nil {
				_, err = fetchSnapshot(ctx, pgDocker.DB, acc1)
				assert.ErrorIs(t, err, pgx.ErrNoRows)

				_, err = fetchSnapshot(ctx, pgDocker.DB, acc2)
				assert.ErrorIs(t, err, pgx.ErrNoRows)
			} else {
				snap, err := fetchSnapshot(ctx, pgDocker.DB, acc1)
				assert.NoError(t, err)
				assert.Equal(t, tt.wants.snapshot.acc1Credit, snap.credit)
				assert.Equal(t, tt.wants.snapshot.acc1Debit, snap.debit)

				snap, err = fetchSnapshot(ctx, pgDocker.DB, acc2)
				assert.NoError(t, err)
				assert.Equal(t, tt.wants.snapshot.acc2Credit, snap.credit)
				assert.Equal(t, tt.wants.snapshot.acc2Debit, snap.debit)
			}
		})
	}
}

func TestLedgerRepository_GetAccountBalanceFailure(t *testing.T) {
	t.Run("should return an error if account does not exist", func(t *testing.T) {
		r := NewLedgerRepository(pgDocker.DB, logrus.New())

		acc, err := vos.NewSingleAccount(testdata.GenerateAccountPath())
		assert.NoError(t, err)

		_, err = r.GetAccountBalance(context.Background(), acc)
		assert.ErrorIs(t, app.ErrAccountNotFound, err)
	})
}

type snapshot struct {
	credit int
	debit  int
	date   time.Time
}

func fetchSnapshot(ctx context.Context, db *pgxpool.Pool, account vos.Account) (snapshot, error) {
	const query = "select credit, debit, tx_date from account_balance where account = $1;"

	var snap snapshot

	err := db.QueryRow(ctx, query, account.Value()).Scan(
		&snap.credit,
		&snap.debit,
		&snap.date,
	)
	if err != nil {
		return snapshot{}, fmt.Errorf("failed to fetch snapshot: %w", err)
	}

	return snap, nil
}
