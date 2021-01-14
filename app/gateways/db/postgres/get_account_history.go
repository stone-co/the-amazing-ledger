package postgres

import (
	"context"
	"time"

	"github.com/stone-co/the-amazing-ledger/app/domain/vos"
)

func (r *LedgerRepository) GetAccountHistory(ctx context.Context, accountName vos.AccountName) (*vos.AccountHistory, error) {
	query := `
		SELECT
			amount,
			operation,
			created_at
		FROM entries
		WHERE account_class = $1 AND account_group = $2 AND account_subgroup = $3 AND account_id = $4
		ORDER BY version;
	`

	rows, err := r.db.Query(
		context.Background(),
		query,
		accountName.Class.String(),
		accountName.Group,
		accountName.Subgroup,
		accountName.ID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var entriesHistory []vos.EntryHistory
	for rows.Next() {
		var amount int
		var operation string
		var createdAt time.Time

		if err = rows.Scan(
			&amount,
			&operation,
			&createdAt,
		); err != nil {
			return nil, err
		}

		operationType := vos.OperationTypeFromString(operation)
		entryHistory, errH := vos.NewEntryHistory(operationType, amount, createdAt)

		if errH != nil {
			return nil, errH
		}

		entriesHistory = append(entriesHistory, *entryHistory)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}

	accountHistory, err := vos.NewAccountHistory(accountName, entriesHistory...)
	if err != nil {
		return nil, err
	}

	return &accountHistory, nil
}
