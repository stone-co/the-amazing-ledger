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
    max(version) as current_version,
    SUM(
        CASE operation
        WHEN 'credit' THEN amount
        ELSE 0
        END
    ) AS total_credit,
    SUM(
        CASE operation
        WHEN 'debit' THEN amount
        ELSE 0
        END
    ) AS total_debit
from
     entry
where account = $1
group by account;
`

func (r LedgerRepository) GetAccountBalance(ctx context.Context, account vos.AccountPath) (*vos.AccountBalance, error) {
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
			return nil, app.ErrAccountNotFound
		}
		return nil, err
	}

	accountBalance := vos.NewAccountBalance(account, vos.Version(currentVersion), totalCredit, totalDebit)
	return accountBalance, nil

}

const getAccountBalanceAggregatedQuery = `
select
    max(version) as current_version,
    SUM(
        CASE operation
        WHEN 'credit' THEN amount
        ELSE 0
        END
    ) AS total_credit,
    SUM(
        CASE operation
        WHEN 'debit' THEN amount
        ELSE 0
        END
    ) AS total_debit
from
     entry
where account <@ $1
group by account;
`

func (r LedgerRepository) GetAccountBalanceAggregated(ctx context.Context, account vos.AccountPath) (*vos.AccountBalance, error) {
	const operation = "Repository.GetAccountBalanceAggregated"

	defer newrelic.NewDatastoreSegment(ctx, collection, operation, getAccountBalanceAggregatedQuery).End()

	row := r.db.QueryRow(
		context.Background(),
		getAccountBalanceAggregatedQuery,
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
			return nil, app.ErrAccountNotFound
		}
		return nil, err
	}

	accountBalance := vos.NewAccountBalance(account, vos.Version(currentVersion), totalCredit, totalDebit)
	return accountBalance, nil
}
