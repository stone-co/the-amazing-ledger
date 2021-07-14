package postgres

import (
	"context"
	"fmt"
	"time"

	"github.com/stone-co/the-amazing-ledger/app/domain/vos"
	"github.com/stone-co/the-amazing-ledger/app/instrumentation/newrelic"
	pag "github.com/stone-co/the-amazing-ledger/app/pagination"
)

const (
	_accountEntriesQueryPrefix = `
select
	id,
	version,
    operation,
    amount,
	event,
	competence_date,
	metadata
from
	entry
where
	account = $1
	and competence_date >= $2
	and competence_date < $3
`

	_accountEntriesQueryPagination = `
	and (competence_date, version) <= ($5, $6)
`

	_accountEntriesQuerySuffix = `
order by
	competence_date desc,
	version desc
limit $4;
`
)

type listAccountEntriesCursor struct {
	CompetenceDate time.Time `json:"competence_date"`
	Version        int64     `json:"version"`
}

func (r LedgerRepository) ListAccountEntries(ctx context.Context, req vos.AccountEntryRequest) ([]vos.AccountEntry, pag.Cursor, error) {
	const op = "Repository.ListAccountEntries"

	query, args, err := generateListAccountEntriesQuery(req)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to generate %s query: %w", op, err)
	}

	defer newrelic.NewDatastoreSegment(ctx, collection, op, query).End()

	rows, err := r.db.Query(ctx, query, args...)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to execute query: %w", err)
	}

	defer rows.Close()

	entries := make([]vos.AccountEntry, 0)

	for rows.Next() {
		var entry vos.AccountEntry
		var operation int8

		if err = rows.Scan(
			&entry.ID,
			&entry.Version,
			&operation,
			&entry.Amount,
			&entry.Event,
			&entry.CompetenceDate,
			&entry.Metadata,
		); err != nil {
			return nil, nil, fmt.Errorf("failed to scan row: %w", err)
		}

		entry.Operation = vos.OperationType(operation)

		entries = append(entries, entry)
	}

	if err = rows.Err(); err != nil {
		return nil, nil, fmt.Errorf("%s rows have error: %w", op, err)
	}

	if len(entries) <= req.Page.Size {
		return entries, nil, nil
	}

	lastEntry := entries[len(entries)-1]
	entries = entries[:len(entries)-1]

	cursor, err := pag.NewCursor(listAccountEntriesCursor{
		CompetenceDate: lastEntry.CompetenceDate,
		Version:        lastEntry.Version.AsInt64(),
	})
	if err != nil {
		return nil, nil, fmt.Errorf("failed to generate next page token: %w", err)
	}

	return entries, cursor, nil
}

func generateListAccountEntriesQuery(req vos.AccountEntryRequest) (string, []interface{}, error) {
	query := _accountEntriesQueryPrefix
	args := []interface{}{req.Account.Name(), req.StartDate, req.EndDate, req.Page.Size + 1}

	if req.Page.Cursor != nil {
		var cursor listAccountEntriesCursor
		err := req.Page.Extract(&cursor)
		if err != nil {
			return "", nil, err
		}

		query += _accountEntriesQueryPagination
		args = append(args, cursor.CompetenceDate, cursor.Version)
	}
	query += _accountEntriesQuerySuffix

	return query, args, nil
}
