package postgres

import (
	"context"

	"github.com/jackc/pgx/v4"
	"github.com/stone-co/the-amazing-ledger/pkg/command-handler/domain/ledger/entities"
)

func (r *LedgerRepository) CreateTransaction(ctx context.Context, transaction *entities.Transaction) error {
	tx, err := r.db.Begin(ctx)
	if err != nil {
		return err
	}

	defer tx.Rollback(ctx)

	query := `
		INSERT INTO
			entries (
				id,
				account_id,
	  			operation,
				amount,
				version,
				transaction_id,
				created_at
			) VALUES ($1, $2, $3, $4, $5, $6, $7)
	`

	var batch pgx.Batch

	for _, entry := range transaction.Entries {
		batch.Queue(
			query,
			entry.ID,
			entry.AccountID,
			entry.Operation.String(),
			entry.Amount,
			entry.Version,
			transaction.ID,
			transaction.CreatedAt,
		)
	}

	br := tx.SendBatch(ctx, &batch)
	err = br.Close()
	if err != nil {
		// TODO: assuming that is duplicate key.
		return entities.ErrIdempotencyKey
	}

	err = tx.Commit(ctx)
	if err != nil {
		return err
	}

	return nil
}
