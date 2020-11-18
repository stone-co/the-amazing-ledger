package postgres

import (
	"context"

	"github.com/stone-co/the-amazing-ledger/pkg/command-handler/domain/ledger/entities"
)

func (r *LedgerRepository) GetAccountBalance(ctx context.Context, accountName entities.AccountName) (*entities.AccountBalance, error) {
	query := `
		SELECT
			account_class,
			account_group,
			account_subgroup,
			account_id,
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
		WHERE account_class = $3 AND account_group = $4 AND account_subgroup = $5 AND account_id = $6
		GROUP BY account_class, account_group, account_subgroup, account_id
	`
	creditOperation := entities.CreditOperation.String()
	debitOperation := entities.DebitOperation.String()

	row := r.db.QueryRow(
		context.Background(),
		query,
		creditOperation,
		debitOperation,
		accountName.Class.String(),
		accountName.Group,
		accountName.Subgroup,
		accountName.ID,
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

	accountBalance := entities.NewAccountBalance(accountName, entities.Version(currentVersion), totalCredit, totalDebit)
	return accountBalance, nil

}
