package postgres

import (
	"context"

	"github.com/jackc/pgx/v4"

	"github.com/stone-co/the-amazing-ledger/app"
	"github.com/stone-co/the-amazing-ledger/app/domain/entities"
	"github.com/stone-co/the-amazing-ledger/app/shared/instrumentation/newrelic"
)

const createTransactionQuery = `
insert into entry (id, tx_id, version, operation, company, event, amount, competence_date, account)
values ($1, $2, $3, $4, $5, $6, $7, $8, $9);
`

func (r LedgerRepository) CreateTransaction(ctx context.Context, transaction entities.Transaction) error {
	const operation = "Repository.CreateTransaction"

	defer newrelic.NewDatastoreSegment(ctx, collection, operation, createTransactionQuery).End()

	tx, err := r.db.Begin(ctx)
	if err != nil {
		return err
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
			entry.Version,
			entry.Operation.String(),
			transaction.Company,
			transaction.Event,
			entry.Amount,
			transaction.CompetenceDate,
			entry.Account.Name(),
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
