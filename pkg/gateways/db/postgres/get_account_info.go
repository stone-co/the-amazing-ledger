package postgres

import (
	"context"

	"github.com/stone-co/the-amazing-ledger/pkg/command-handler/domain/ledger/entities"
)

func (r *LedgerRepository) GetAccountInfo(ctx context.Context, accountID string) (*entities.AccountInfo, error) {
	query := `
		SELECT
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
		WHERE account_id = $3
		GROUP BY account_id
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

	row := tx.QueryRow(context.Background(), query, creditOperation, debitOperation, accountID)

	var version uint64
	var totalCredit int
	var totalDebit int

	err = row.Scan(
		&accountID,
		&version,
		&totalCredit,
		&totalDebit,
	)
	if err != nil {
		return nil, err
	}

	accountInfo := entities.NewAccountInfo(accountID, entities.Version(version), totalCredit, totalDebit)
	return accountInfo, nil

}
