package postgres

import (
	"context"
	"errors"
	"strings"
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

func TestLedgerRepository_CreateTransaction(t *testing.T) {
	event := uint32(1)
	company := "abc"
	competenceDate := time.Now()

	r := NewLedgerRepository(pgDocker.DB, logrus.New())
	ctx := context.Background()

	_, err := pgDocker.DB.Exec(ctx, `insert into event (name) values ('default');`)
	assert.NoError(t, err)

	t.Run("insert transaction successfully with no previous versions - auto version", func(t *testing.T) {
		e1, _ := entities.NewEntry(
			uuid.New(),
			vos.DebitOperation,
			"liability.abc.account1",
			vos.NextAccountVersion,
			100,
		)
		e2, _ := entities.NewEntry(
			uuid.New(),
			vos.CreditOperation,
			"liability.abc.account2",
			vos.IgnoreAccountVersion,
			100,
		)

		tx, err := entities.NewTransaction(uuid.New(), event, company, competenceDate, e1, e2)
		assert.NoError(t, err)

		err = r.CreateTransaction(ctx, tx)
		assert.NoError(t, err)

		ev1, err := fetchEntryVersion(ctx, pgDocker.DB, e1.ID)
		assert.NoError(t, err)
		assert.Equal(t, vos.Version(1), ev1)

		ev2, err := fetchEntryVersion(ctx, pgDocker.DB, e2.ID)
		assert.NoError(t, err)
		assert.Equal(t, vos.IgnoreAccountVersion, ev2)

		v1, err := fetchAccountVersion(ctx, pgDocker.DB, e1.Account)
		assert.NoError(t, err)
		assert.Equal(t, vos.Version(1), v1)

		v2, err := fetchAccountVersion(ctx, pgDocker.DB, e2.Account)
		assert.NoError(t, err)
		assert.Equal(t, vos.Version(0), v2)
	})

	t.Run("insert transaction successfully with no previous versions - manual version", func(t *testing.T) {
		e1, _ := entities.NewEntry(
			uuid.New(),
			vos.DebitOperation,
			"liability.abc.account3",
			vos.Version(3),
			100,
		)
		e2, _ := entities.NewEntry(
			uuid.New(),
			vos.CreditOperation,
			"liability.abc.account4",
			vos.IgnoreAccountVersion,
			100,
		)

		tx, err := entities.NewTransaction(uuid.New(), event, company, competenceDate, e1, e2)
		assert.NoError(t, err)

		err = r.CreateTransaction(ctx, tx)
		assert.NoError(t, err)

		ev1, err := fetchEntryVersion(ctx, pgDocker.DB, e1.ID)
		assert.NoError(t, err)
		assert.Equal(t, vos.Version(1), ev1)

		ev2, err := fetchEntryVersion(ctx, pgDocker.DB, e2.ID)
		assert.NoError(t, err)
		assert.Equal(t, vos.IgnoreAccountVersion, ev2)

		v1, err := fetchAccountVersion(ctx, pgDocker.DB, e1.Account)
		assert.NoError(t, err)
		assert.Equal(t, vos.Version(1), v1)

		v2, err := fetchAccountVersion(ctx, pgDocker.DB, e2.Account)
		assert.NoError(t, err)
		assert.Equal(t, vos.Version(0), v2)
	})

	t.Run("insert transaction successfully with existing versions - auto version", func(t *testing.T) {
		e1, _ := entities.NewEntry(
			uuid.New(),
			vos.DebitOperation,
			"liability.abc.account1",
			vos.NextAccountVersion,
			100,
		)
		e2, _ := entities.NewEntry(
			uuid.New(),
			vos.CreditOperation,
			"liability.abc.account2",
			vos.IgnoreAccountVersion,
			100,
		)

		tx, err := entities.NewTransaction(uuid.New(), event, company, competenceDate, e1, e2)
		assert.NoError(t, err)

		err = r.CreateTransaction(ctx, tx)
		assert.NoError(t, err)

		ev1, err := fetchEntryVersion(ctx, pgDocker.DB, e1.ID)
		assert.NoError(t, err)
		assert.Equal(t, vos.Version(2), ev1)

		ev2, err := fetchEntryVersion(ctx, pgDocker.DB, e2.ID)
		assert.NoError(t, err)
		assert.Equal(t, vos.IgnoreAccountVersion, ev2)

		v1, err := fetchAccountVersion(ctx, pgDocker.DB, e1.Account)
		assert.NoError(t, err)
		assert.Equal(t, vos.Version(2), v1)

		v2, err := fetchAccountVersion(ctx, pgDocker.DB, e2.Account)
		assert.NoError(t, err)
		assert.Equal(t, vos.Version(0), v2)
	})

	t.Run("insert transaction successfully with existing versions - manual version", func(t *testing.T) {
		e1, _ := entities.NewEntry(
			uuid.New(),
			vos.DebitOperation,
			"liability.abc.account1",
			vos.Version(3),
			100,
		)
		e2, _ := entities.NewEntry(
			uuid.New(),
			vos.CreditOperation,
			"liability.abc.account2",
			vos.IgnoreAccountVersion,
			100,
		)

		tx, err := entities.NewTransaction(uuid.New(), event, company, competenceDate, e1, e2)
		assert.NoError(t, err)

		err = r.CreateTransaction(ctx, tx)
		assert.NoError(t, err)

		ev1, err := fetchEntryVersion(ctx, pgDocker.DB, e1.ID)
		assert.NoError(t, err)
		assert.Equal(t, vos.Version(3), ev1)

		ev2, err := fetchEntryVersion(ctx, pgDocker.DB, e2.ID)
		assert.NoError(t, err)
		assert.Equal(t, vos.IgnoreAccountVersion, ev2)

		v1, err := fetchAccountVersion(ctx, pgDocker.DB, e1.Account)
		assert.NoError(t, err)
		assert.Equal(t, vos.Version(3), v1)

		v2, err := fetchAccountVersion(ctx, pgDocker.DB, e2.Account)
		assert.NoError(t, err)
		assert.Equal(t, vos.Version(0), v2)
	})

	t.Run("return error when sending same version", func(t *testing.T) {
		e1, _ := entities.NewEntry(
			uuid.New(),
			vos.DebitOperation,
			"liability.abc.account1",
			vos.Version(3),
			100,
		)
		e2, _ := entities.NewEntry(
			uuid.New(),
			vos.CreditOperation,
			"liability.abc.account2",
			vos.IgnoreAccountVersion,
			100,
		)

		tx, err := entities.NewTransaction(uuid.New(), event, company, competenceDate, e1, e2)
		assert.NoError(t, err)

		err = r.CreateTransaction(ctx, tx)
		assert.ErrorIs(t, err, app.ErrInvalidVersion)

		v1, err := fetchAccountVersion(ctx, pgDocker.DB, e1.Account)
		assert.NoError(t, err)
		assert.Equal(t, vos.Version(3), v1)

		v2, err := fetchAccountVersion(ctx, pgDocker.DB, e2.Account)
		assert.NoError(t, err)
		assert.Equal(t, vos.Version(0), v2)
	})

	t.Run("return error when sending lower version", func(t *testing.T) {
		e1, _ := entities.NewEntry(
			uuid.New(),
			vos.DebitOperation,
			"liability.abc.account1",
			vos.Version(1),
			100,
		)
		e2, _ := entities.NewEntry(
			uuid.New(),
			vos.CreditOperation,
			"liability.abc.account2",
			vos.IgnoreAccountVersion,
			100,
		)

		tx, err := entities.NewTransaction(uuid.New(), event, company, competenceDate, e1, e2)
		assert.NoError(t, err)

		err = r.CreateTransaction(ctx, tx)
		assert.ErrorIs(t, err, app.ErrInvalidVersion)

		v1, err := fetchAccountVersion(ctx, pgDocker.DB, e1.Account)
		assert.NoError(t, err)
		assert.Equal(t, vos.Version(3), v1)

		v2, err := fetchAccountVersion(ctx, pgDocker.DB, e2.Account)
		assert.NoError(t, err)
		assert.Equal(t, vos.Version(0), v2)
	})

	t.Run("return error when sending random high version", func(t *testing.T) {
		e1, _ := entities.NewEntry(
			uuid.New(),
			vos.DebitOperation,
			"liability.abc.account1",
			vos.Version(30),
			100,
		)
		e2, _ := entities.NewEntry(
			uuid.New(),
			vos.CreditOperation,
			"liability.abc.account2",
			vos.IgnoreAccountVersion,
			100,
		)

		tx, err := entities.NewTransaction(uuid.New(), event, company, competenceDate, e1, e2)
		assert.NoError(t, err)

		err = r.CreateTransaction(ctx, tx)
		assert.ErrorIs(t, err, app.ErrInvalidVersion)

		v1, err := fetchAccountVersion(ctx, pgDocker.DB, e1.Account)
		assert.NoError(t, err)
		assert.Equal(t, vos.Version(3), v1)

		v2, err := fetchAccountVersion(ctx, pgDocker.DB, e2.Account)
		assert.NoError(t, err)
		assert.Equal(t, vos.Version(0), v2)
	})

	t.Run("return error when reusing entry id", func(t *testing.T) {
		e1, _ := entities.NewEntry(
			uuid.New(),
			vos.DebitOperation,
			"liability.abc.account1",
			vos.NextAccountVersion,
			100,
		)
		e2, _ := entities.NewEntry(
			uuid.New(),
			vos.CreditOperation,
			"liability.abc.account2",
			vos.IgnoreAccountVersion,
			100,
		)

		tx, err := entities.NewTransaction(uuid.New(), event, company, competenceDate, e1, e2)
		assert.NoError(t, err)

		err = r.CreateTransaction(ctx, tx)
		assert.NoError(t, err)

		ev1, err := fetchEntryVersion(ctx, pgDocker.DB, e1.ID)
		assert.NoError(t, err)
		assert.Equal(t, vos.Version(4), ev1)

		ev2, err := fetchEntryVersion(ctx, pgDocker.DB, e2.ID)
		assert.NoError(t, err)
		assert.Equal(t, vos.IgnoreAccountVersion, ev2)

		v1, err := fetchAccountVersion(ctx, pgDocker.DB, e1.Account)
		assert.NoError(t, err)
		assert.Equal(t, vos.Version(4), v1)

		v2, err := fetchAccountVersion(ctx, pgDocker.DB, e2.Account)
		assert.NoError(t, err)
		assert.Equal(t, vos.Version(0), v2)

		err = r.CreateTransaction(ctx, tx)
		assert.ErrorIs(t, err, app.ErrIdempotencyKeyViolation)

		v1, err = fetchAccountVersion(ctx, pgDocker.DB, e1.Account)
		assert.NoError(t, err)
		assert.Equal(t, vos.Version(4), v1)

		v2, err = fetchAccountVersion(ctx, pgDocker.DB, e2.Account)
		assert.NoError(t, err)
		assert.Equal(t, vos.Version(0), v2)
	})
}

func fetchAccountVersion(ctx context.Context, db *pgxpool.Pool, account vos.AccountPath) (vos.Version, error) {
	const query = `select coalesce(version, 0) from account_version where account = $1;`

	var version int64
	err := db.QueryRow(ctx, query, account.Name()).Scan(&version)
	if err != nil && err != pgx.ErrNoRows {
		return vos.Version(0), err
	} else if errors.Is(err, pgx.ErrNoRows) {
		return vos.Version(0), nil
	}

	return vos.Version(version), nil
}

func fetchEntryVersion(ctx context.Context, db *pgxpool.Pool, id uuid.UUID) (vos.Version, error) {
	const query = `select version from entry where id = $1;`

	var version int64
	if err := db.QueryRow(ctx, query, id).Scan(&version); err != nil {
		return 0, err
	}

	return vos.Version(version), nil
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
				insert into entry (id, tx_id, event, operation, version, amount, competence_date, account, company)
				values ($1, $2, $3, $4, $5, $6, $7, $8, $9),
				($10, $11, $12, $13, $14, $15, $16, $17, $18);`,
		},
		{
			name: "should create query with 3 entries successfully",
			size: 3,
			expected: `
				insert into entry (id, tx_id, event, operation, version, amount, competence_date, account, company)
				values ($1, $2, $3, $4, $5, $6, $7, $8, $9),
				($10, $11, $12, $13, $14, $15, $16, $17, $18),
				($19, $20, $21, $22, $23, $24, $25, $26, $27);`,
		},
		{
			name: "should create query with 4 entries successfully",
			size: 4,
			expected: `
				insert into entry (id, tx_id, event, operation, version, amount, competence_date, account, company)
				values ($1, $2, $3, $4, $5, $6, $7, $8, $9),
				($10, $11, $12, $13, $14, $15, $16, $17, $18),
				($19, $20, $21, $22, $23, $24, $25, $26, $27),
				($28, $29, $30, $31, $32, $33, $34, $35, $36);`,
		},
		{
			name: "should create query with 5 entries successfully",
			size: 5,
			expected: `
				insert into entry (id, tx_id, event, operation, version, amount, competence_date, account, company)
				values ($1, $2, $3, $4, $5, $6, $7, $8, $9),
				($10, $11, $12, $13, $14, $15, $16, $17, $18),
				($19, $20, $21, $22, $23, $24, $25, $26, $27),
				($28, $29, $30, $31, $32, $33, $34, $35, $36),
				($37, $38, $39, $40, $41, $42, $43, $44, $45);`,
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
