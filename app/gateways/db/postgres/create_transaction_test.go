package postgres

import (
	"context"
	"encoding/json"
	"errors"
	"strings"
	"testing"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/stone-co/the-amazing-ledger/app"
	"github.com/stone-co/the-amazing-ledger/app/domain/entities"
	"github.com/stone-co/the-amazing-ledger/app/domain/vos"
	"github.com/stone-co/the-amazing-ledger/app/tests"
)

func TestLedgerRepository_CreateTransactionSuccess(t *testing.T) {
	testCases := []struct {
		name                   string
		repoSetup              func(t *testing.T, ctx context.Context, r *LedgerRepository)
		entriesSetup           func(t *testing.T) []entities.Entry
		expectedEntryVersion   vos.Version
		expectedAccountVersion vos.Version
	}{
		{
			name:      "insert transaction successfully with no previous versions - auto version",
			repoSetup: func(t *testing.T, ctx context.Context, r *LedgerRepository) {},
			entriesSetup: func(t *testing.T) []entities.Entry {
				e1 := createEntry(t, vos.DebitOperation, "liability.abc.account1", vos.NextAccountVersion)
				e2 := createEntry(t, vos.CreditOperation, "liability.abc.account2", vos.IgnoreAccountVersion)

				return []entities.Entry{e1, e2}
			},
			expectedEntryVersion:   vos.Version(1),
			expectedAccountVersion: vos.Version(1),
		},
		{
			name:      "insert transaction successfully with no previous versions - manual version",
			repoSetup: func(t *testing.T, ctx context.Context, r *LedgerRepository) {},
			entriesSetup: func(t *testing.T) []entities.Entry {
				e1 := createEntry(t, vos.DebitOperation, "liability.abc.account1", vos.Version(1))
				e2 := createEntry(t, vos.CreditOperation, "liability.abc.account2", vos.IgnoreAccountVersion)

				return []entities.Entry{e1, e2}
			},
			expectedEntryVersion:   vos.Version(1),
			expectedAccountVersion: vos.Version(1),
		},
		{
			name: "insert transaction successfully with existing versions - auto version",
			repoSetup: func(t *testing.T, ctx context.Context, r *LedgerRepository) {
				e1 := createEntry(t, vos.DebitOperation, "liability.abc.account1", vos.NextAccountVersion)
				e2 := createEntry(t, vos.CreditOperation, "liability.abc.account2", vos.IgnoreAccountVersion)

				err := createTransaction(t, ctx, r, e1, e2)
				assert.NoError(t, err)
			},
			entriesSetup: func(t *testing.T) []entities.Entry {
				e1 := createEntry(t, vos.DebitOperation, "liability.abc.account1", vos.NextAccountVersion)
				e2 := createEntry(t, vos.CreditOperation, "liability.abc.account2", vos.IgnoreAccountVersion)

				return []entities.Entry{e1, e2}
			},
			expectedEntryVersion:   vos.Version(2),
			expectedAccountVersion: vos.Version(2),
		},
		{
			name: "insert transaction successfully with existing versions - manual version",
			repoSetup: func(t *testing.T, ctx context.Context, r *LedgerRepository) {
				e1 := createEntry(t, vos.DebitOperation, "liability.abc.account1", vos.NextAccountVersion)
				e2 := createEntry(t, vos.CreditOperation, "liability.abc.account2", vos.IgnoreAccountVersion)

				err := createTransaction(t, ctx, r, e1, e2)
				assert.NoError(t, err)
			},
			entriesSetup: func(t *testing.T) []entities.Entry {
				e1 := createEntry(t, vos.DebitOperation, "liability.abc.account1", vos.Version(2))
				e2 := createEntry(t, vos.CreditOperation, "liability.abc.account2", vos.IgnoreAccountVersion)

				return []entities.Entry{e1, e2}
			},
			expectedEntryVersion:   vos.Version(2),
			expectedAccountVersion: vos.Version(2),
		},
	}

	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			r := NewLedgerRepository(pgDocker.DB, logrus.New())
			ctx := context.Background()

			defer tests.TruncateTables(ctx, pgDocker.DB, "entry", "account_version")

			tt.repoSetup(t, ctx, r)

			entries := tt.entriesSetup(t)
			e1, e2 := entries[0], entries[1]

			err := createTransaction(t, ctx, r, e1, e2)
			assert.NoError(t, err)

			assertMetadata(t, ctx, pgDocker.DB, e1.ID, e1.Metadata)
			assertMetadata(t, ctx, pgDocker.DB, e2.ID, e2.Metadata)

			assertEntryVersion(t, ctx, pgDocker.DB, e1.ID, tt.expectedEntryVersion)
			assertEntryVersion(t, ctx, pgDocker.DB, e2.ID, vos.IgnoreAccountVersion)

			assertAccountVersion(t, ctx, pgDocker.DB, e1.Account, tt.expectedAccountVersion)
			assertAccountVersion(t, ctx, pgDocker.DB, e2.Account, vos.Version(0))
		})
	}
}

func TestLedgerRepository_CreateTransactionFailure(t *testing.T) {
	e1 := createEntry(t, vos.DebitOperation, "liability.abc.account1", vos.NextAccountVersion)
	e2 := createEntry(t, vos.CreditOperation, "liability.abc.account2", vos.IgnoreAccountVersion)

	testCases := []struct {
		name                   string
		repoSetup              func(t *testing.T, ctx context.Context, r *LedgerRepository)
		entriesSetup           func(t *testing.T) []entities.Entry
		expectedErr            error
		expectedAccountVersion vos.Version
	}{
		{
			name: "return error when sending same version",
			repoSetup: func(t *testing.T, ctx context.Context, r *LedgerRepository) {
				err := createTransaction(t, ctx, r, e1, e2)
				assert.NoError(t, err)
			},
			entriesSetup: func(t *testing.T) []entities.Entry {
				e3 := createEntry(t, vos.DebitOperation, "liability.abc.account1", vos.Version(1))
				e4 := createEntry(t, vos.CreditOperation, "liability.abc.account2", vos.IgnoreAccountVersion)

				return []entities.Entry{e3, e4}
			},
			expectedErr:            app.ErrInvalidVersion,
			expectedAccountVersion: vos.Version(1),
		},
		{
			name: "return error when sending lower version",
			repoSetup: func(t *testing.T, ctx context.Context, r *LedgerRepository) {
				err := createTransaction(t, ctx, r, e1, e2)
				assert.NoError(t, err)

				e3 := createEntry(t, vos.DebitOperation, "liability.abc.account1", vos.NextAccountVersion)
				e4 := createEntry(t, vos.CreditOperation, "liability.abc.account2", vos.IgnoreAccountVersion)

				err = createTransaction(t, ctx, r, e3, e4)
				assert.NoError(t, err)
			},
			entriesSetup: func(t *testing.T) []entities.Entry {
				e5 := createEntry(t, vos.DebitOperation, "liability.abc.account1", vos.Version(1))
				e6 := createEntry(t, vos.CreditOperation, "liability.abc.account2", vos.IgnoreAccountVersion)

				return []entities.Entry{e5, e6}
			},
			expectedErr:            app.ErrInvalidVersion,
			expectedAccountVersion: vos.Version(2),
		},
		{
			name: "return error when sending random high version",
			repoSetup: func(t *testing.T, ctx context.Context, r *LedgerRepository) {
				err := createTransaction(t, ctx, r, e1, e2)
				assert.NoError(t, err)
			},
			entriesSetup: func(t *testing.T) []entities.Entry {
				e3 := createEntry(t, vos.DebitOperation, "liability.abc.account1", vos.Version(30))
				e4 := createEntry(t, vos.CreditOperation, "liability.abc.account2", vos.IgnoreAccountVersion)

				return []entities.Entry{e3, e4}
			},
			expectedErr:            app.ErrInvalidVersion,
			expectedAccountVersion: vos.Version(1),
		},
		{
			name: "return error when reusing entry id",
			repoSetup: func(t *testing.T, ctx context.Context, r *LedgerRepository) {
				err := createTransaction(t, ctx, r, e1, e2)
				assert.NoError(t, err)
			},
			entriesSetup: func(t *testing.T) []entities.Entry {
				return []entities.Entry{e1, e2}
			},
			expectedErr:            app.ErrIdempotencyKeyViolation,
			expectedAccountVersion: vos.Version(1),
		},
	}

	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			r := NewLedgerRepository(pgDocker.DB, logrus.New())
			ctx := context.Background()

			defer tests.TruncateTables(ctx, pgDocker.DB, "entry", "account_version")

			tt.repoSetup(t, ctx, r)

			entries := tt.entriesSetup(t)
			e1, e2 := entries[0], entries[1]

			err := createTransaction(t, ctx, r, e1, e2)
			assert.ErrorIs(t, err, tt.expectedErr)

			assertAccountVersion(t, ctx, pgDocker.DB, e1.Account, tt.expectedAccountVersion)
			assertAccountVersion(t, ctx, pgDocker.DB, e2.Account, vos.Version(0))
		})
	}
}

func assertAccountVersion(t *testing.T, ctx context.Context, db *pgxpool.Pool, account vos.AccountPath, want vos.Version) {
	t.Helper()

	const query = `select coalesce(version, 0) from account_version where account = $1;`

	var version vos.Version

	err := db.QueryRow(ctx, query, account.Name()).Scan(&version)
	if err != nil && !errors.Is(err, pgx.ErrNoRows) {
		t.Fatalf("unexpected error: %v", err)
	} else if errors.Is(err, pgx.ErrNoRows) {
		assert.Equal(t, vos.Version(0), version)
	}

	assert.Equal(t, want, version)
}

func assertEntryVersion(t *testing.T, ctx context.Context, db *pgxpool.Pool, id uuid.UUID, want vos.Version) {
	t.Helper()

	const query = `select version from entry where id = $1;`

	var version vos.Version

	err := db.QueryRow(ctx, query, id).Scan(&version)
	require.NoError(t, err)

	assert.Equal(t, want, version)
}

func assertMetadata(t *testing.T, ctx context.Context, db *pgxpool.Pool, id uuid.UUID, want json.RawMessage) {
	t.Helper()

	const query = `select metadata from entry where id = $1`

	var metadata json.RawMessage

	err := db.QueryRow(ctx, query, id).Scan(&metadata)
	require.NoError(t, err)

	assert.Equal(t, string(want), string(metadata))
}

func Test_buildQuery(t *testing.T) {
	testCases := []struct {
		name     string
		size     int
		expected string
	}{
		{
			name: "should create query with 2 entries successfully",
			size: 2,
			expected: `
				insert into entry (id, tx_id, event, operation, version, amount, competence_date, account, company, metadata)
				values ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10),
				($11, $12, $13, $14, $15, $16, $17, $18, $19, $20);`,
		},
		{
			name: "should create query with 3 entries successfully",
			size: 3,
			expected: `
				insert into entry (id, tx_id, event, operation, version, amount, competence_date, account, company, metadata)
				values ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10),
				($11, $12, $13, $14, $15, $16, $17, $18, $19, $20),
				($21, $22, $23, $24, $25, $26, $27, $28, $29, $30);`,
		},
		{
			name: "should create query with 4 entries successfully",
			size: 4,
			expected: `
				insert into entry (id, tx_id, event, operation, version, amount, competence_date, account, company, metadata)
				values ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10),
				($11, $12, $13, $14, $15, $16, $17, $18, $19, $20),
				($21, $22, $23, $24, $25, $26, $27, $28, $29, $30),
				($31, $32, $33, $34, $35, $36, $37, $38, $39, $40);`,
		},
		{
			name: "should create query with 5 entries successfully",
			size: 5,
			expected: `
				insert into entry (id, tx_id, event, operation, version, amount, competence_date, account, company, metadata)
				values ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10),
				($11, $12, $13, $14, $15, $16, $17, $18, $19, $20),
				($21, $22, $23, $24, $25, $26, $27, $28, $29, $30),
				($31, $32, $33, $34, $35, $36, $37, $38, $39, $40),
				($41, $42, $43, $44, $45, $46, $47, $48, $49, $50);`,
		},
	}

	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			got := buildQuery(tt.size)
			want := strings.ReplaceAll(tt.expected, "\t", "")

			assert.Equal(t, want, got)
		})
	}
}
