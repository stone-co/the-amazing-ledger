package postgres

import (
	"context"

	"github.com/stone-co/the-amazing-ledger/pkg/command-handler/domain/ledger/entities"
)

func (r *LedgerRepository) GetAccountInfo(ctx context.Context, accountName *entities.AccountName) (*entities.AccountInfo, error) {
	query := `
		SELECT
			account_class,
			account_group,
			account_subgroup,
			account_id,
		MAX(version) as version,
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

	tx, err := r.db.Begin(ctx)
	if err != nil {
		return nil, err
	}
	defer func() {
		_ = tx.Rollback(ctx)
	}()

	row := tx.QueryRow(
		context.Background(),
		query,
		creditOperation,
		debitOperation,
		accountName.Class.String(),
		accountName.Group,
		accountName.Subgroup,
		accountName.ID,
	)
	var empty string
	var version uint64
	var totalCredit int
	var totalDebit int

	err = row.Scan(
		&empty,
		&empty,
		&empty,
		&empty,
		&version,
		&totalCredit,
		&totalDebit,
	)
	if err != nil {
		return nil, err
	}

	accountPath := accountName.Name()
	accountInfo := entities.NewAccountInfo(accountPath, entities.Version(version), totalCredit, totalDebit)
	return accountInfo, nil

}
