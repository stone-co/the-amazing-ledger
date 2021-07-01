package postgres

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgconn"
	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v4"

	"github.com/stone-co/the-amazing-ledger/app"
	"github.com/stone-co/the-amazing-ledger/app/domain/entities"
	"github.com/stone-co/the-amazing-ledger/app/shared/instrumentation/newrelic"
)

const createTransactionQuery = `
insert into entry (id, tx_id, event, operation, version, amount, competence_date, account, company)
values ($1, $2, $3, $4, $5, $6, $7, $8, $9);
`

func (r LedgerRepository) CreateTransaction(ctx context.Context, transaction entities.Transaction) error {
	const operation = "Repository.CreateTransaction"

	defer newrelic.NewDatastoreSegment(ctx, collection, operation, createTransactionQuery).End()

	tx, err := r.db.Begin(ctx)
	if err != nil {
		return fmt.Errorf("failed to begin db transaction: %w", err)
	}

	defer func(tx pgx.Tx, ctx context.Context) {
		_ = tx.Rollback(ctx)
	}(tx, ctx)

	var batch pgx.Batch
	for _, entry := range transaction.Entries {
		batch.Queue(
			createTransactionQuery,
			entry.ID,
			transaction.ID,
			transaction.Event,
			entry.Operation,
			entry.Version,
			entry.Amount,
			transaction.CompetenceDate,
			entry.Account.Name(),
			transaction.Company,
		)
	}

	br := tx.SendBatch(ctx, &batch)
	err = br.Close()
	if err == nil {
		_ = tx.Commit(ctx) // TODO: double check
		return nil
	}

	var pgErr *pgconn.PgError
	if ok := errors.As(err, &pgErr); !ok {
		return fmt.Errorf("failed to create transaction: %w", err)
	}

	if pgErr.Code == pgerrcode.RaiseException {
		return app.ErrInvalidVersion
	} else if pgErr.Code == pgerrcode.UniqueViolation {
		return app.ErrIdempotencyKeyViolation
	}

	return fmt.Errorf("failed to create transaction: %w", pgErr)
}
