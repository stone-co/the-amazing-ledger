package postgres

import (
	"context"
	"log"
	"time"

	"github.com/stone-co/the-amazing-ledger/app"
	"github.com/stone-co/the-amazing-ledger/app/domain/vos"
)

const accountReportQuery = `
select 
	account,
	coalesce(MAX(version),0) as current_version,
	coalesce(SUM(CASE operation  WHEN $1 THEN amount  ELSE 0 END),0) AS credit,
	coalesce(SUM(CASE operation WHEN $2 THEN amount ELSE 0 END),0) AS debit
from
	entry
where
	account ~ $3
	and created_at >= $4
	and created_at < $5
group by 1;
`

func (r *LedgerRepository) GetSyntheticReport(ctx context.Context, account vos.AccountPath, startTime time.Time, endTime time.Time) (*vos.SyntheticReport, error) {
	creditOperation := vos.CreditOperation.String()
	debitOperation := vos.DebitOperation.String()

	rows, errQuery := r.db.Query(
		context.Background(),
		accountReportQuery,
		creditOperation,
		debitOperation,
		account,
		startTime,
		endTime,
	)

	if errQuery != nil {
		log.Printf("> err query: %v\n", errQuery)
		return nil, errQuery
	}

	log.Println("> query run ok!")

	defer rows.Close()

	pathsReport := []vos.Path{}
	var currentVersion uint64
	var totalCredit int
	var totalDebit int

	log.Println("> rows.next()")

	for rows.Next() {
		var credit int
		var debit int

		log.Println("> scan results!")

		err := rows.Scan(
			nil,
			nil,
			nil,
			nil,
			&currentVersion,
			&credit,
			&debit,
		)

		if err != nil {
			log.Printf("> error on Scan: %v", err)
			return nil, err
		}

		path := vos.Path{
			Account: account,
			Credit:  credit,
			Debit:   debit,
		}

		log.Printf("> path: %v", path)

		pathsReport = append(pathsReport, path)

		totalCredit = totalCredit + credit
		totalDebit = totalDebit + debit
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
