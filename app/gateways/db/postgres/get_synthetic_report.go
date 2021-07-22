package postgres

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"time"

	"github.com/jackc/pgx/v4"

	"github.com/stone-co/the-amazing-ledger/app/domain/vos"
)

const syntheticReportQuery = `
select 
	subpath(account, 0, $1),
	coalesce(SUM(CASE operation WHEN %d THEN amount ELSE 0::bigint END),0::bigint) AS creditSum, 
	coalesce(SUM(CASE operation WHEN %d THEN amount ELSE 0::bigint END),0::bigint) AS debitSum 
from 
	entry 
where 
	account ~ $2
`

const timeParamsQueryStart = " and created_at >= $3 "
const timeParamsQueryEnd = " and created_at < $4 "

const groupByQuery = "group by 1;"

func (r *LedgerRepository) GetSyntheticReport(ctx context.Context, query vos.Account, level int, startTime time.Time, endTime time.Time) (*vos.SyntheticReport, error) {
	const operation = "Repository.GetSyntheticReport"

	sqlQuery, params := buildQueryAndParams(query, level, startTime, endTime)

	paramsPgx := make([]interface{}, len(params))
	for i, s := range params {
		paramsPgx[i] = s
	}

	defer r.pb.MonitorDataSegment(ctx, collection, operation, sqlQuery).End()
	rows, errQuery := r.db.Query(
		ctx,
		sqlQuery,
		paramsPgx...,
	)

	if errQuery != nil {
		if errors.Is(errQuery, pgx.ErrNoRows) {
			return &vos.SyntheticReport{}, nil
		}
		return nil, errQuery
	}

	defer rows.Close()

	pathsReport := []vos.Path{}
	var totalCredit int64
	var totalDebit int64

	for rows.Next() {
		var accStr string
		var credit int64
		var debit int64

		err := rows.Scan(
			&accStr,
			&credit,
			&debit,
		)

		if err != nil {
			return nil, err
		}

		account, err := vos.NewAnalyticalAccount(accStr)
		if err != nil {
			return nil, err
		}

		path := vos.Path{
			Account: account,
			Credit:  credit,
			Debit:   debit,
		}

		pathsReport = append(pathsReport, path)

		totalCredit = totalCredit + credit
		totalDebit = totalDebit + debit
	}

	errNext := rows.Err()
	if errNext != nil {
		return nil, errNext
	}

	if pathsReport == nil || len(pathsReport) < 1 {
		return &vos.SyntheticReport{}, nil
	}

	syntheticReport, errEntity := vos.NewSyntheticReport(totalCredit, totalDebit, pathsReport)
	if errEntity != nil {
		return nil, errEntity
	}

	return syntheticReport, nil
}

func buildQueryAndParams(query vos.Account, level int, startTime time.Time, endTime time.Time) (string, []string) {
	sqlQuery := syntheticReportQuery
	sqlQuery = fmt.Sprintf(sqlQuery, vos.CreditOperation, vos.DebitOperation)

	params := []string{
		strconv.Itoa(level),
		query.Value(),
	}

	if !startTime.IsZero() {
		params = append(params, startTime.Format(time.RFC3339))
		sqlQuery += timeParamsQueryStart
	}

	if !endTime.IsZero() {
		params = append(params, endTime.Format(time.RFC3339))
		sqlQuery += timeParamsQueryEnd
	}

	sqlQuery += groupByQuery

	return sqlQuery, params
}
