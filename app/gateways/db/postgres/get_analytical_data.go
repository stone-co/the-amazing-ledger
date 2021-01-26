package postgres

import (
	"context"

	"github.com/stone-co/the-amazing-ledger/app/domain/vos"
	"github.com/stone-co/the-amazing-ledger/app/shared/instrumentation/newrelic"
)

func (r *LedgerRepository) GetAnalyticalData(ctx context.Context, path vos.AccountPath, fn func(vos.Statement) error) error {
	operation := "Repository.GetAnalyticalData"
	query := `
	SELECT
		account_class,
		account_group,
		account_subgroup,
		account_id,
		operation,
		amount
	FROM
		entries
	`

	args := []interface{}{}

	if path.TotalLevels >= 1 {
		query += " WHERE account_class = $1"
		args = append(args, path.Class.String())

		if path.TotalLevels >= 2 {
			query += " AND account_group = $2"
			args = append(args, path.Group)

			if path.TotalLevels >= 3 {
				query += " AND account_subgroup = $3"
				args = append(args, path.Subgroup)
			}
		}
	}

	query += " ORDER BY version"

	defer newrelic.NewDatastoreSegment(ctx, collection, operation, query).End()

	rows, err := r.db.Query(ctx, query, args...)
	if err != nil {
		return err
	}
	defer rows.Close()

	for rows.Next() {
		var class string
		var group string
		var subgroup string
		var id string
		var op string
		var amount int

		if err = rows.Scan(
			&class,
			&group,
			&subgroup,
			&id,
			&op,
			&amount,
		); err != nil {
			return err
		}

		account := vos.FormatAccount(class, group, subgroup, id)

		err = fn(vos.Statement{
			Account:   account,
			Operation: vos.OperationTypeFromString(op),
			Amount:    amount,
		})
		if err != nil {
			return err
		}
	}

	if err := rows.Err(); err != nil {
		return err
	}

	return nil
}
