package postgres

import (
	"context"

	"github.com/stone-co/the-amazing-ledger/app/domain/entities"
	"github.com/stone-co/the-amazing-ledger/app/domain/vos"
)

func (r *LedgerRepository) LoadObjectsIntoCache(ctx context.Context, cachedAccounts *entities.CachedAccounts) (vos.Version, error) {
	query := `
		SELECT account_class, account_group, account_subgroup, account_id, MAX(version) As version
		FROM entries
		GROUP BY account_class, account_group, account_subgroup, account_id
		ORDER BY version desc
	`

	rows, err := r.db.Query(ctx, query)
	if err != nil {
		return 0, err
	}
	defer rows.Close()

	var maxVersion vos.Version

	for rows.Next() {
		var accClass string
		var accGroup string
		var accSubgroup string
		var accID string
		var version vos.Version

		if err := rows.Scan(
			&accClass,
			&accGroup,
			&accSubgroup,
			&accID,
			&version,
		); err != nil {
			return 0, err
		}

		// TODO: check for duplicated?
		account := vos.FormatAccount(accClass, accGroup, accSubgroup, accID)
		cachedAccounts.Store(account, version)

		if version > maxVersion {
			maxVersion = version
		}
	}

	if err := rows.Err(); err != nil {
		return 0, err
	}

	return maxVersion, nil
}
