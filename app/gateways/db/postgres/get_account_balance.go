package postgres

import (
	"context"
	"errors"
	"time"

	"github.com/jackc/pgconn"
	"github.com/jackc/pgerrcode"

	"github.com/stone-co/the-amazing-ledger/app"
	"github.com/stone-co/the-amazing-ledger/app/domain/vos"
	"github.com/stone-co/the-amazing-ledger/app/shared/instrumentation/newrelic"
)

const getAccountBalanceQuery = `
select
    credit_balance  as credit,
    debit_balance   as debit,
    dt              as tx_date,
    version         as tx_version
from
    get_account_balance($1)
;
`

func (r LedgerRepository) GetAccountBalance(ctx context.Context, account vos.AccountPath) (vos.AccountBalance, error) {
	const operation = "Repository.GetAccountBalance"

	defer newrelic.NewDatastoreSegment(ctx, collection, operation, getAccountBalanceQuery).End()

	var totalCredit int
	var totalDebit int
	var lastTransactionTime time.Time
	var currentVersion int64

	err := r.db.QueryRow(ctx, getAccountBalanceQuery, account.Name()).Scan(
		&totalCredit,
		&totalDebit,
		&lastTransactionTime,
		&currentVersion,
	)

	if err == nil {
		return vos.NewAccountBalance(
			account,
			vos.Version(currentVersion),
			totalCredit,
			totalDebit,
			lastTransactionTime,
		), nil
	}

	var pgErr *pgconn.PgError
	if !errors.As(err, &pgErr) {
		return vos.AccountBalance{}, err
	}

	if pgErr.Code == pgerrcode.NoDataFound {
		return vos.AccountBalance{}, app.ErrAccountNotFound
	}

	return vos.AccountBalance{}, err
}
