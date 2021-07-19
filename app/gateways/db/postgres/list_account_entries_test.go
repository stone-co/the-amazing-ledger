package postgres

import (
	"context"
	"encoding/json"
	"testing"
	"time"

	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"

	"github.com/stone-co/the-amazing-ledger/app"
	"github.com/stone-co/the-amazing-ledger/app/domain/entities"
	"github.com/stone-co/the-amazing-ledger/app/domain/vos"
	"github.com/stone-co/the-amazing-ledger/app/pagination"
	"github.com/stone-co/the-amazing-ledger/app/tests"
)

func Test_generateListAccountEntriesQuery(t *testing.T) {
	account, err := vos.NewSingleAccount("liability.test.account1")
	assert.NoError(t, err)

	size := 10

	end := time.Now().UTC().Round(time.Microsecond)
	start := end.Add(-10 * time.Second)

	version := vos.Version(1)

	testCases := []struct {
		name          string
		req           func() vos.AccountEntryRequest
		expectedQuery string
		expectedArgs  []interface{}
		expectedErr   error
	}{
		{
			name: "valid - without pagination",
			req: func() vos.AccountEntryRequest {
				return vos.AccountEntryRequest{
					Account:   account,
					StartDate: start,
					EndDate:   end,
					Page: pagination.Page{
						Size:   size,
						Cursor: nil,
					},
				}
			},
			expectedQuery: _accountEntriesQueryPrefix + _accountEntriesQuerySuffix,
			expectedArgs:  []interface{}{account.Value(), start, end, size + 1},
			expectedErr:   nil,
		},
		{
			name: "valid - with pagination",
			req: func() vos.AccountEntryRequest {
				cursor, _ := pagination.NewCursor(listAccountEntriesCursor{
					CompetenceDate: end,
					Version:        1,
				})

				return vos.AccountEntryRequest{
					Account:   account,
					StartDate: start,
					EndDate:   end,
					Page: pagination.Page{
						Size:   size,
						Cursor: cursor,
					},
				}
			},
			expectedQuery: _accountEntriesQueryPrefix + _accountEntriesQueryPagination + _accountEntriesQuerySuffix,
			expectedArgs:  []interface{}{account.Value(), start, end, size + 1, end, version.AsInt64()},
			expectedErr:   nil,
		},
		{
			name: "invalid page	token",
			req: func() vos.AccountEntryRequest {
				cursor, _ := pagination.NewCursor(map[string]interface{}{"version": "none"})

				return vos.AccountEntryRequest{
					Account:   account,
					StartDate: start,
					EndDate:   end,
					Page: pagination.Page{
						Size:   size,
						Cursor: cursor,
					},
				}
			},
			expectedQuery: "",
			expectedArgs:  nil,
			expectedErr:   app.ErrInvalidPageCursor,
		},
	}
	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			got, got1, err := generateListAccountEntriesQuery(tt.req())
			assert.ErrorIs(t, err, tt.expectedErr)
			assert.Equal(t, tt.expectedQuery, got)
			assert.EqualValues(t, tt.expectedArgs, got1)
		})
	}
}

func TestLedgerRepository_ListAccountEntries(t *testing.T) {
	type w struct {
		entries []vos.AccountEntry
		cursor  pagination.Cursor
	}

	const (
		account1 = "liability.abc.account1"
		account2 = "liability.abc.account2"
		amount   = 100
	)

	testCases := []struct {
		name         string
		seedRepo     func(*testing.T, context.Context, *LedgerRepository) []entities.Transaction
		setupRequest func(*testing.T, []entities.Transaction) vos.AccountEntryRequest
		want         func(*testing.T, []entities.Transaction) w
	}{
		{
			name: "no exiting entries case",
			seedRepo: func(t *testing.T, ctx context.Context, r *LedgerRepository) []entities.Transaction {
				e1 := createEntry(t, vos.DebitOperation, account1, vos.Version(1), amount)
				e2 := createEntry(t, vos.CreditOperation, account2, vos.IgnoreAccountVersion, amount)

				tx := createTransaction(t, ctx, r, e1, e2)

				return []entities.Transaction{tx}
			},
			setupRequest: func(t *testing.T, _ []entities.Transaction) vos.AccountEntryRequest {
				account, err := vos.NewSingleAccount("liability.abc.account3")
				assert.NoError(t, err)

				now := time.Now()

				return vos.AccountEntryRequest{
					Account:   account,
					StartDate: now.Add(-10 * time.Second),
					EndDate:   now.Add(10 * time.Second),
					Page: pagination.Page{
						Size:   10,
						Cursor: nil,
					},
				}
			},
			want: func(_ *testing.T, _ []entities.Transaction) w {
				return w{
					entries: []vos.AccountEntry{},
					cursor:  nil,
				}
			},
		},
		{
			name: "return all entries",
			seedRepo: func(t *testing.T, ctx context.Context, r *LedgerRepository) []entities.Transaction {
				e1 := createEntry(t, vos.DebitOperation, account1, vos.Version(1), amount)
				e2 := createEntry(t, vos.CreditOperation, account2, vos.IgnoreAccountVersion, amount)

				tx := createTransaction(t, ctx, r, e1, e2)

				return []entities.Transaction{tx}
			},
			setupRequest: func(t *testing.T, _ []entities.Transaction) vos.AccountEntryRequest {
				account, err := vos.NewSingleAccount(account1)
				assert.NoError(t, err)

				now := time.Now()

				return vos.AccountEntryRequest{
					Account:   account,
					StartDate: now.Add(-10 * time.Second),
					EndDate:   now.Add(10 * time.Second),
					Page: pagination.Page{
						Size:   10,
						Cursor: nil,
					},
				}
			},
			want: func(_ *testing.T, txs []entities.Transaction) w {
				entries := accountEntriesFromTransaction(t, txs[0], account1)

				return w{
					entries: entries,
					cursor:  nil,
				}
			},
		},
		{
			name: "return first page",
			seedRepo: func(t *testing.T, ctx context.Context, r *LedgerRepository) []entities.Transaction {
				e1 := createEntry(t, vos.DebitOperation, account1, vos.Version(1), amount)
				e2 := createEntry(t, vos.CreditOperation, account2, vos.IgnoreAccountVersion, amount)

				tx1 := createTransaction(t, ctx, r, e1, e2)

				e1 = createEntry(t, vos.DebitOperation, account1, vos.Version(2), amount)
				e2 = createEntry(t, vos.CreditOperation, account2, vos.IgnoreAccountVersion, amount)

				tx2 := createTransaction(t, ctx, r, e1, e2)

				return []entities.Transaction{tx1, tx2}
			},
			setupRequest: func(t *testing.T, _ []entities.Transaction) vos.AccountEntryRequest {
				account, err := vos.NewSingleAccount(account1)
				assert.NoError(t, err)

				now := time.Now()

				return vos.AccountEntryRequest{
					Account:   account,
					StartDate: now.Add(-10 * time.Second),
					EndDate:   now.Add(10 * time.Second),
					Page: pagination.Page{
						Size:   1,
						Cursor: nil,
					},
				}
			},
			want: func(t *testing.T, txs []entities.Transaction) w {
				entries := accountEntriesFromTransaction(t, txs[1], account1)
				cur := cursorFromTransaction(t, txs[0], account1)

				return w{
					entries: entries,
					cursor:  cur,
				}
			},
		},
		{
			name: "return second page",
			seedRepo: func(t *testing.T, ctx context.Context, r *LedgerRepository) []entities.Transaction {
				e1 := createEntry(t, vos.DebitOperation, account1, vos.Version(1), amount)
				e2 := createEntry(t, vos.CreditOperation, account2, vos.IgnoreAccountVersion, amount)

				tx1 := createTransaction(t, ctx, r, e1, e2)

				e1 = createEntry(t, vos.DebitOperation, account1, vos.Version(2), amount)
				e2 = createEntry(t, vos.CreditOperation, account2, vos.IgnoreAccountVersion, amount)

				tx2 := createTransaction(t, ctx, r, e1, e2)

				return []entities.Transaction{tx1, tx2}
			},
			setupRequest: func(t *testing.T, txs []entities.Transaction) vos.AccountEntryRequest {
				account, err := vos.NewSingleAccount(account1)
				assert.NoError(t, err)

				cur := cursorFromTransaction(t, txs[0], account1)

				now := time.Now()

				return vos.AccountEntryRequest{
					Account:   account,
					StartDate: now.Add(-10 * time.Second),
					EndDate:   now.Add(10 * time.Second),
					Page: pagination.Page{
						Size:   1,
						Cursor: cur,
					},
				}
			},
			want: func(t *testing.T, txs []entities.Transaction) w {
				entries := accountEntriesFromTransaction(t, txs[0], account1)

				return w{
					entries: entries,
					cursor:  nil,
				}
			},
		},
	}

	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			r := NewLedgerRepository(pgDocker.DB, logrus.New())
			ctx := context.Background()

			tests.TruncateTables(ctx, pgDocker.DB, "account_version, entry")

			txs := tt.seedRepo(t, ctx, r)

			req := tt.setupRequest(t, txs)

			resp, cur, err := r.ListAccountEntries(ctx, req)
			want := tt.want(t, txs)
			got := w{entries: resp, cursor: cur}
			assert.NoError(t, err)
			assert.Equal(t, want, got)
		})
	}
}

func accountEntriesFromTransaction(t *testing.T, tx entities.Transaction, account string) []vos.AccountEntry {
	t.Helper()

	act := make([]vos.AccountEntry, 0, len(tx.Entries))
	for _, et := range tx.Entries {
		if et.Account.Value() != account {
			continue
		}

		var mt map[string]interface{}
		err := json.Unmarshal(et.Metadata, &mt)
		assert.NoError(t, err)

		act = append(act, vos.AccountEntry{
			ID:             et.ID,
			Version:        et.Version,
			Operation:      et.Operation,
			Amount:         et.Amount,
			Event:          int(tx.Event),
			CompetenceDate: tx.CompetenceDate.Round(time.Microsecond),
			Metadata:       mt,
		})
	}

	return act
}

func cursorFromTransaction(t *testing.T, tx entities.Transaction, account string) pagination.Cursor {
	t.Helper()

	var et entities.Entry
	for _, entry := range tx.Entries {
		if entry.Account.Value() == account {
			et = entry
			break
		}
	}
	assert.NotEmpty(t, et)

	cur, err := pagination.NewCursor(listAccountEntriesCursor{
		CompetenceDate: tx.CompetenceDate.Round(time.Microsecond),
		Version:        et.Version.AsInt64(),
	})
	assert.NoError(t, err)

	return cur
}
