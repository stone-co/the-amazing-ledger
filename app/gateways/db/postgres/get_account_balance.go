package postgres

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgconn"
	"github.com/jackc/pgerrcode"

	"github.com/stone-co/the-amazing-ledger/app"
	"github.com/stone-co/the-amazing-ledger/app/domain/vos"
	"github.com/stone-co/the-amazing-ledger/app/instrumentation/newrelic"
)

const getAccountBalanceQuery = `
select
    credit_balance  as credit,
    debit_balance   as debit,
    version         as tx_version
from
    get_account_balance($1)
;
`

func (r LedgerRepository) GetAccountBalance(ctx context.Context, account vos.Account) (vos.AccountBalance, error) {
	const operation = "Repository.GetAccountBalance"

	defer newrelic.NewDatastoreSegment(ctx, collection, operation, getAccountBalanceQuery).End()

	var totalCredit int
	var totalDebit int
	var currentVersion int64

	err := r.db.QueryRow(ctx, getAccountBalanceQuery, account.Value()).Scan(
		&totalCredit,
		&totalDebit,
		&currentVersion,
	)

	if err != nil {
		var pgErr *pgconn.PgError
		if !errors.As(err, &pgErr) {
			return vos.AccountBalance{}, fmt.Errorf("failed to get account balance: %w", err)
		}

		if pgErr.Code == pgerrcode.NoDataFound {
			return vos.AccountBalance{}, app.ErrAccountNotFound
		}

		return vos.AccountBalance{}, fmt.Errorf("failed to get account balance: %w", pgErr)
	}

	return vos.NewAccountBalance(
		account,
		vos.Version(currentVersion),
		totalCredit,
		totalDebit,
	), nil
}
