package postgres

import (
	"context"

	"github.com/jackc/pgx/v4"

	"github.com/stone-co/the-amazing-ledger/app"
	"github.com/stone-co/the-amazing-ledger/app/domain/vos"
	"github.com/stone-co/the-amazing-ledger/app/shared/instrumentation/newrelic"
)

const getAccountBalanceQuery = `
select
    credit_balance as credit,
    debit_balance as debit
from
    get_account_balance($1)
;
`

func (r LedgerRepository) GetAccountBalance(ctx context.Context, account vos.AccountPath) (vos.AccountBalance, error) {
	const operation = "Repository.GetAccountBalance"

	defer newrelic.NewDatastoreSegment(ctx, collection, operation, getAccountBalanceQuery).End()

	row := r.db.QueryRow(
		context.Background(),
		getAccountBalanceQuery,
		account.Name(),
	)
	var currentVersion uint64
	var totalCredit int
	var totalDebit int

	err := row.Scan(
		&currentVersion,
		&totalCredit,
		&totalDebit,
	)
	if err != nil {
		if err == pgx.ErrNoRows {
			return vos.AccountBalance{}, app.ErrAccountNotFound
		}
		return vos.AccountBalance{}, err
	}

	accountBalance := vos.NewAccountBalance(account, vos.Version(currentVersion), totalCredit, totalDebit)
	return accountBalance, nil

}
