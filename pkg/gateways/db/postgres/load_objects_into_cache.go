package postgres

import (
	"context"

	"github.com/stone-co/the-amazing-ledger/pkg/command-handler/domain/ledger/entities"
)

func (r *LedgerRepository) LoadObjectsIntoCache(ctx context.Context, cachedAccounts *entities.CachedAccounts) (entities.Version, error) {
	query := `
		SELECT account_id, MAX(version) As version
			FROM entries
			GROUP BY account_id
			ORDER BY version desc
	`

	rows, err := r.db.Query(ctx, query)
	if err != nil {
		return 0, err
	}
	defer rows.Close()

	var maxVersion entities.Version

	for rows.Next() {
		var accountID string
		var version entities.Version

		if err := rows.Scan(
			&accountID,
			&version,
		); err != nil {
			return 0, err
		}

		// TODO: check for duplicated?
		cachedAccounts.Store(accountID, version)

		if version > maxVersion {
			maxVersion = version
		}
	}

	if err := rows.Err(); err != nil {
		return 0, err
	}

	return maxVersion, nil
}
