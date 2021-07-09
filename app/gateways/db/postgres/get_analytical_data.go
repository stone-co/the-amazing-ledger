package postgres

import (
	"context"
	"fmt"
	"time"

	"github.com/stone-co/the-amazing-ledger/app/domain/vos"
	"github.com/stone-co/the-amazing-ledger/app/shared/instrumentation/newrelic"
)

const analyticalDataQuery = `
select
	account,
    operation,
    amount,
    created_at
from
	entry
where
	account ~ $1
order by
	version;
`

func (r *LedgerRepository) GetAnalyticalData(ctx context.Context, query vos.AccountQuery, fn func(vos.Statement) error) error {
	const operation = "Repository.GetAnalyticalData"

	defer newrelic.NewDatastoreSegment(ctx, collection, operation, analyticalDataQuery).End()

	rows, err := r.db.Query(ctx, analyticalDataQuery, query.Value())
	if err != nil {
		return fmt.Errorf("failed to execute query: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var account string
		var operation string
		var amount int
		var createdAt time.Time

		if err = rows.Scan(
			&account,
			&amount,
			&createdAt,
		); err != nil {
			return fmt.Errorf("failed to scan row: %w", err)
		}

		err = fn(vos.Statement{
			Account:   account,
			Operation: vos.OperationTypeFromString(operation),
			Amount:    amount,
		})
		if err != nil {
			return err
		}
	}

	if err := rows.Err(); err != nil {
		return fmt.Errorf("%s rows have error: %w", operation, err)
	}

	return nil
}
