package postgres

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgconn"
	"github.com/jackc/pgerrcode"

	"github.com/stone-co/the-amazing-ledger/app"
	"github.com/stone-co/the-amazing-ledger/app/domain/vos"
	"github.com/stone-co/the-amazing-ledger/app/shared/instrumentation/newrelic"
)

const queryAggregatedBalanceQuery = `
select query_aggregated_account_balance($1);
`

func (r LedgerRepository) QueryAggregatedBalance(ctx context.Context, query vos.AccountQuery) (vos.QueryBalance, error) {
	const operation = "Repository.QueryAggregatedBalance"

	defer newrelic.NewDatastoreSegment(ctx, collection, operation, queryAggregatedBalanceQuery).End()

	var balance int

	err := r.db.QueryRow(ctx, queryAggregatedBalanceQuery, query.Value()).Scan(&balance)
	if err != nil {
		var pgErr *pgconn.PgError
		if !errors.As(err, &pgErr) {
			return vos.QueryBalance{}, fmt.Errorf("failed to query aggregated balance: %w", err)
		}

		if pgErr.Code == pgerrcode.NoDataFound {
			return vos.QueryBalance{}, app.ErrAccountNotFound
		}

		return vos.QueryBalance{}, fmt.Errorf("failed to query aggregated balance: %w", pgErr)
	}

	return vos.NewQueryBalance(query, balance), nil
}
