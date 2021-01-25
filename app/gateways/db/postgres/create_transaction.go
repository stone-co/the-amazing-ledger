package postgres

import (
	"context"

	"github.com/jackc/pgx/v4"
	"github.com/newrelic/go-agent/v3/newrelic"
	"github.com/stone-co/the-amazing-ledger/app"
	"github.com/stone-co/the-amazing-ledger/app/domain/entities"
)

func (r *LedgerRepository) CreateTransaction(ctx context.Context, transaction *entities.Transaction) error {
	operation := "Repository.CreateTransaction"

	query := `
		INSERT INTO
			entries (
				id,
				account_class,
				account_group,
				account_subgroup,
				account_id,
	  			operation,
				amount,
				version,
				transaction_id
			) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
	`

	txn := newrelic.FromContext(ctx)
	seg := newrelic.DatastoreSegment{
		Product:            newrelic.DatastorePostgres,
		Collection:         "entries",
		Operation:          operation,
		ParameterizedQuery: query,
	}
	seg.StartTime = txn.StartSegmentNow()
	defer seg.End()

	tx, err := r.db.Begin(ctx)
	if err != nil {
		return err
	}

	defer func() {
		_ = tx.Rollback(ctx)
	}()

	var batch pgx.Batch

	for _, entry := range transaction.Entries {
		batch.Queue(
			query,
			entry.ID,
			entry.Account.Class.String(),
			entry.Account.Group,
			entry.Account.Subgroup,
			entry.Account.ID,
			entry.Operation.String(),
			entry.Amount,
			entry.Version,
			transaction.ID,
		)
	}

	br := tx.SendBatch(ctx, &batch)
	err = br.Close()
	if err != nil {
		// TODO: assuming that is duplicate key.
		return app.ErrIdempotencyKeyViolation
	}

	err = tx.Commit(ctx)
	if err != nil {
		return err
	}

	return nil
}
