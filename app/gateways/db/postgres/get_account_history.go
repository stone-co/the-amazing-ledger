package postgres

import (
	"context"
	"time"

	"github.com/stone-co/the-amazing-ledger/app/domain/vos"
)

func (r *LedgerRepository) GetAccountHistory(ctx context.Context, accountName vos.AccountName, fn func(vos.EntryHistory) error) error {
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
		return err
	}
	defer rows.Close()

	for rows.Next() {
		var amount int
		var operation string
		var createdAt time.Time

		if err = rows.Scan(
			&amount,
			&operation,
			&createdAt,
		); err != nil {
			return err
		}

		err = fn(vos.EntryHistory{
			Amount:    amount,
			Operation: vos.OperationTypeFromString(operation),
			CreatedAt: createdAt,
		})

		if err != nil {
			return err
		}
	}

	if err = rows.Err(); err != nil {
		return err
	}

	return nil
}
