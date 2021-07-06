package postgres

import (
	"context"
	"log"
	"time"

	"github.com/jackc/pgx/v4"
	"github.com/stone-co/the-amazing-ledger/app"
	"github.com/stone-co/the-amazing-ledger/app/domain/vos"
	"github.com/stone-co/the-amazing-ledger/app/shared/instrumentation/newrelic"
)

const accountReportQuery = `
select 
	account, 
	coalesce(SUM(CASE operation WHEN $1 THEN amount ELSE 0::bigint END),0::bigint) AS creditSum, 
	coalesce(SUM(CASE operation WHEN $2 THEN amount ELSE 0::bigint END),0::bigint) AS debitSum 
from 
	entry 
where 
	account ~ $3
	and created_at >= $4 
	and created_at < $5 
group by 1;
`

//TODO rename param 'account' to 'base_account'?
func (r *LedgerRepository) GetSyntheticReport(ctx context.Context, account vos.AccountPath, startTime time.Time, endTime time.Time) (*vos.SyntheticReport, error) {
	const operation = "Repository.GetSyntheticReport"

	defer newrelic.NewDatastoreSegment(ctx, collection, operation, accountReportQuery).End()

	rows, errQuery := r.db.Query(
		context.Background(),
		accountReportQuery,
		vos.CreditOperation,
		vos.DebitOperation,
		account.Name()+".*",
		startTime,
		endTime,
	)

	if errQuery != nil {
		log.Printf("> err query: %v\n", errQuery)
		if errQuery == pgx.ErrNoRows {
			return nil, app.ErrAccountNotFound
		}

		return nil, errQuery
	}

	log.Println("> query run ok!")

	defer rows.Close()

	pathsReport := []vos.Path{}
	var totalCredit int64
	var totalDebit int64

	for rows.Next() {
		var accStr string
		var creditX int64
		var debitX int64

		log.Println("> scan results!")

		err := rows.Scan(
			&accStr,
			&creditX,
			&debitX,
		)

		if err != nil {
			log.Printf("> error on Scan: %v", err)
			return nil, err
		}

		path := vos.Path{
			Account: account,
			Credit:  creditX,
			Debit:   debitX,
		}

		log.Printf("> path: %v", path)

		pathsReport = append(pathsReport, path)

		totalCredit = totalCredit + creditX
		totalDebit = totalDebit + debitX
	}

	errNext := rows.Err()
	if errNext != nil {
		log.Printf("> error on Next: %v", errNext)
		return nil, errNext
	}

	if pathsReport == nil || len(pathsReport) < 1 {
		return nil, app.ErrAccountNotFound
	}

	log.Printf("> pathsReport: %v", pathsReport)

	syntheticReport, errEntity := vos.NewSyntheticReport(totalCredit, totalDebit, pathsReport)
	if errEntity != nil {
		log.Printf("> error on entity: %v", errEntity)
		return nil, errEntity
	}

	return syntheticReport, nil

}
