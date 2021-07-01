package postgres

import (
	"context"
	"time"

	"github.com/stone-co/the-amazing-ledger/app/domain/vos"
	"github.com/stone-co/the-amazing-ledger/app/shared/instrumentation/newrelic"
)

const accountHistoryQuery = `
select
    operation,
    amount,
    created_at
from
	entry
where
	account = $1
order by
	version;
`

func (r LedgerRepository) GetAccountHistory(ctx context.Context, accountPath vos.AccountPath, fn func(vos.EntryHistory) error) error {
	const operation = "Repository.GetAccountHistory"

	defer newrelic.NewDatastoreSegment(ctx, collection, operation, accountHistoryQuery).End()

	rows, err := r.db.Query(
		context.Background(),
		accountHistoryQuery,
		accountPath.Name(),
	)
	if err != nil {
		return err
	}
	defer rows.Close()

	for rows.Next() {
		var amount int
		var operation int8
		var createdAt time.Time

		if err = rows.Scan(
			&amount,
			&operation,
			&createdAt,
		); err != nil {
			return err
		}

		err = fn(vos.EntryHistory{
			Amount:    amount,
			Operation: vos.OperationType(operation),
			CreatedAt: createdAt,
		})

		if err != nil {
			return err
		}
	}

	if err = rows.Err(); err != nil {
		return err
	}

	return nil
}
