package postgres

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"strconv"

	"github.com/jackc/pgconn"
	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v4"

	"github.com/stone-co/the-amazing-ledger/app"
	"github.com/stone-co/the-amazing-ledger/app/domain/entities"
	"github.com/stone-co/the-amazing-ledger/app/shared/instrumentation/newrelic"
)

const queryArgsLength = 9

var createTransactionQueryMap = map[int]string{
	// 2: `insert into entry (id, tx_id, event, operation, version, amount, competence_date, account, company) values ($1, $2, $3, $4, $5, $6, $7, $8, $9), ($10, $11, $12, $13, $14, $15, $16, $17, $18);`,
}

func (r LedgerRepository) CreateTransaction(ctx context.Context, transaction entities.Transaction) error {
	const operation = "Repository.CreateTransaction"

	query := getQuery1(len(transaction.Entries))
	// query := getQuery2(len(transaction.Entries))

	defer newrelic.NewDatastoreSegment(ctx, collection, operation, query).End()

	tx, err := r.db.Begin(ctx)
	if err != nil {
		return err
	}

	defer func(tx pgx.Tx, ctx context.Context) {
		_ = tx.Rollback(ctx)
	}(tx, ctx)

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
		)
	}

	_, err = tx.Exec(ctx, query, args...)
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

	_ = tx.Commit(ctx)

	return nil
}

func getQuery1(entriesSize int) string {
	query, ok := createTransactionQueryMap[entriesSize]
	if ok {
		return query
	}

	maxArgs := entriesSize * queryArgsLength
	query = `insert into entry (id, tx_id, event, operation, version, amount, competence_date, account, company) values ($1`

	for i := 2; i <= maxArgs; i++ {
		query += fmt.Sprintf(", $%d", i)

		if i%queryArgsLength == 0 {
			if i != maxArgs {
				i += 1
				query += fmt.Sprintf("), ($%d", i)
			} else {
				query += ");"
			}
		}
	}

	// update query map
	createTransactionQueryMap[entriesSize] = query

	return query
}

func getQuery2(entriesSize int) string {
	query, ok := createTransactionQueryMap[entriesSize]
	if ok {
		return query
	}

	query = `insert into entry (id, tx_id, event, operation, version, amount, competence_date, account, company) values %s`
	buffer := bytes.Buffer{}

	for i := 0; i < entriesSize; i++ {
		n := i * queryArgsLength

		buffer.WriteString("(")

		for j := 0; j < queryArgsLength; j++ {
			buffer.WriteString("$")
			buffer.WriteString(strconv.Itoa(n + j + 1))
			buffer.WriteString(", ")
		}

		buffer.Truncate(buffer.Len() - 2)
		buffer.WriteString("), ")
	}

	buffer.Truncate(buffer.Len() - 2)

	query = fmt.Sprintf(query, buffer.String())

	// update query map
	createTransactionQueryMap[entriesSize] = query

	return query
}
