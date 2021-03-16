package postgres

import (
	"context"

	"github.com/stone-co/the-amazing-ledger/app/domain/vos"
	"github.com/stone-co/the-amazing-ledger/app/shared/instrumentation/newrelic"
)

func (r *LedgerRepository) GetAccountBalance(ctx context.Context, accountName vos.AccountName) (*vos.AccountBalance, error) {
	operation := "Repository.GetAccountBalance"
	query := `
		SELECT
			account_class,
			account_group,
			account_subgroup,
			account_id,
			account_suffix,
			MAX(version) as current_version,
			SUM(CASE operation
				WHEN $1 THEN amount
				ELSE 0
				END) AS total_credit,
			SUM(CASE operation
				WHEN $2 THEN amount
				ELSE 0
				END) AS total_debit
		FROM entries
		WHERE account_class = $3 AND account_group = $4 AND account_subgroup = $5 AND account_id = $6 AND account_suffix = $7
		GROUP BY account_class, account_group, account_subgroup, account_id, account_suffix
	`

	defer newrelic.NewDatastoreSegment(ctx, collection, operation, query).End()

	creditOperation := vos.CreditOperation.String()
	debitOperation := vos.DebitOperation.String()

	row := r.db.QueryRow(
		context.Background(),
		query,
		creditOperation,
		debitOperation,
		accountName.Class.String(),
		accountName.Group,
		accountName.Subgroup,
		accountName.ID,
		accountName.Suffix,
	)
	var currentVersion uint64
	var totalCredit int
	var totalDebit int

	err := row.Scan(
		nil,
		nil,
		nil,
		nil,
		&currentVersion,
		&totalCredit,
		&totalDebit,
	)
	if err != nil {
		return nil, err
	}

	accountBalance := vos.NewAccountBalance(accountName, vos.Version(currentVersion), totalCredit, totalDebit)
	return accountBalance, nil

}
