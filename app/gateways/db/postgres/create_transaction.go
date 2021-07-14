package postgres

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/jackc/pgconn"
	"github.com/jackc/pgerrcode"

	"github.com/stone-co/the-amazing-ledger/app"
	"github.com/stone-co/the-amazing-ledger/app/domain/entities"
	"github.com/stone-co/the-amazing-ledger/app/instrumentation/newrelic"
)

const (
	queryArgsLength   = 10
	maxQueriesDefault = 5
)

const createTransactionQuery = `
insert into entry (id, tx_id, event, operation, version, amount, competence_date, account, company, metadata)
values %s;`

var createTransactionQueryMap map[int]string

func (r LedgerRepository) CreateTransaction(ctx context.Context, transaction entities.Transaction) error {
	const operation = "Repository.CreateTransaction"

	query := getQuery(len(transaction.Entries))

	args := make([]interface{}, 0)

	for _, entry := range transaction.Entries {
		args = append(
			args,
			entry.ID,
			transaction.ID,
			transaction.Event,
			entry.Operation,
			entry.Version,
			entry.Amount,
			transaction.CompetenceDate,
			entry.Account.Name(),
			transaction.Company,
			entry.Metadata,
		)
	}

	defer newrelic.NewDatastoreSegment(ctx, collection, operation, query).End()

	_, err := r.db.Exec(ctx, query, args...)
	if err != nil {
		var pgErr *pgconn.PgError
		if ok := errors.As(err, &pgErr); !ok {
			return err
		}

		if pgErr.Code == pgerrcode.RaiseException {
			return app.ErrInvalidVersion
		}

		if pgErr.Code == pgerrcode.UniqueViolation {
			return app.ErrIdempotencyKeyViolation
		}

		return err
	}

	return nil
}

func getQuery(entriesSize int) string {
	query, ok := createTransactionQueryMap[entriesSize]
	if ok {
		return query
	}

	query = buildQuery(entriesSize)
	createTransactionQueryMap[entriesSize] = query

	return query
}

func buildQuery(entriesSize int) string {
	var sb strings.Builder

	for i := 0; i < entriesSize; i++ {
		n := i * queryArgsLength

		sb.WriteString("(")

		for j := 0; j < queryArgsLength; j++ {
			sb.WriteString("$")
			sb.WriteString(strconv.Itoa(n + j + 1))

			if j != queryArgsLength-1 {
				sb.WriteString(", ")
			}
		}

		if i != entriesSize-1 {
			sb.WriteString("),\n")
		}
	}

	sb.WriteString(")")

	return fmt.Sprintf(createTransactionQuery, sb.String())
}
