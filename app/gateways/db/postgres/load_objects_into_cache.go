package postgres

import (
	"context"

	"github.com/stone-co/the-amazing-ledger/app/domain/entities"
	"github.com/stone-co/the-amazing-ledger/app/domain/vos"
	"github.com/stone-co/the-amazing-ledger/app/shared/instrumentation/newrelic"
)

const loadCacheQuery = `
select account, max(version) as version 
from entries
group by account;
`

func (r *LedgerRepository) LoadObjectsIntoCache(ctx context.Context, cachedAccounts *entities.CachedAccounts) (vos.Version, error) {
	const operation = "Repository.LoadObjectsIntoCache"

	defer newrelic.NewDatastoreSegment(ctx, collection, operation, loadCacheQuery).End()

	rows, err := r.db.Query(ctx, loadCacheQuery)
	if err != nil {
		return 0, err
	}
	defer rows.Close()

	var maxVersion vos.Version

	for rows.Next() {
		var account string
		var version vos.Version

		if err := rows.Scan(
			&account,
			&version,
		); err != nil {
			return 0, err
		}

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
