package postgres

import (
	"context"
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"

	"github.com/stone-co/the-amazing-ledger/app"
	"github.com/stone-co/the-amazing-ledger/app/domain/vos"
)

func (r *LedgerRepository) GetSyntheticReport(ctx context.Context, accountName string, startTime time.Time, endTime time.Time) (*vos.SyntheticReport, error) {
	var index uint16 = 2
	var params = []string{}

	creditOperation := vos.CreditOperation.String()
	debitOperation := vos.DebitOperation.String()

	params = append(params, creditOperation)
	params = append(params, debitOperation)

	var columns string = ` SELECT account_class `
	var query string = ` coalesce(MAX(version),0) as current_version, coalesce(SUM(CASE operation  WHEN $1 THEN amount  ELSE 0 END),0) AS credit, coalesce(SUM(CASE operation WHEN $2 THEN amount ELSE 0 END),0) AS debit FROM entries WHERE 1=1 `
	var groupBy string = " GROUP BY 1 "

	var paths []string

	if accountName != "" {
		paths = strings.Split(accountName, ":")
	}

	var class string
	var group string
	var subgroup string
	var id string

	if len(paths) >= 1 {
		class = paths[0]

		index++

		columns += ",account_group"
		query += fmt.Sprint(" AND account_class = $", index)
		query += " AND account_group is not null and account_group != '' "
		groupBy += ",2"

		params = append(params, class)

		if len(paths) >= 2 {
			group = paths[1]

			index++

			columns += ",account_subgroup"
			query += fmt.Sprint(` AND account_group = $`, index)
			query += " AND account_subgroup is not null and account_subgroup != '' "
			groupBy += ",3"

			params = append(params, group)

			if len(paths) >= 3 {
				subgroup = paths[2]

				index++

				columns += ",account_id,"
				query += fmt.Sprint(` AND account_subgroup = $`, index)
				query += " AND account_id is not null and account_id != '' "
				groupBy += ",4"

				params = append(params, subgroup)

				if len(paths) >= 4 {
					id = paths[3]
					index++
					query += fmt.Sprint(` AND account_id = $`, index)
					params = append(params, id)
				}
			}
		}
	}

	var dates string
	var startYear int
	var startMonth int
	var startDay int

	var endYear string
	var endMonth string
	var endDay string

	if !startTime.IsZero() {
		log.Printf("> init date: %v\n", startTime)

		index++
		dates += fmt.Sprintf(" AND date_part('year', created_at) >= $%v::integer and date_part('year', created_at)  <= coalesce($%v::integer, date_part('year', created_at)::integer) ", index, (index + 1))

		index++
		index++
		dates += fmt.Sprintf(" AND date_part('month', created_at) >= $%v::integer and date_part('month', created_at)  <= coalesce($%v::integer, date_part('month', created_at)::integer) ", index, (index + 1))

		index++
		index++
		dates += fmt.Sprintf(" AND date_part('day', created_at) >= $%v::integer and date_part('day', created_at)  <= coalesce($%v::integer, date_part('day', created_at)::integer) ", index, (index + 1))

		log.Printf("> st: %v\n", startTime)

		startYear = startTime.Year()
		startMonth = int(startTime.Month())
		startDay = startTime.Day()

		log.Printf("> sy: %v\n", startYear)
		log.Printf("> sm: %v\n", startMonth)
		log.Printf("> sd: %v\n", startDay)

		if !endTime.IsZero() {
			log.Printf("> end date: %v\n", endTime)

			endYear = strconv.Itoa(endTime.Year())
			endMonth = strconv.Itoa(int(endTime.Month()))
			endDay = strconv.Itoa(endTime.Day())

		} else {
			endYear = strconv.Itoa(startTime.Year())
			endMonth = strconv.Itoa(int(startTime.Month()))
			endDay = strconv.Itoa(startTime.Day())
		}

		params = append(params, strconv.Itoa(startYear))
		params = append(params, endYear)
		params = append(params, strconv.Itoa(startMonth))
		params = append(params, endMonth)
		params = append(params, strconv.Itoa(startDay))
		params = append(params, endDay)
	}

	finalQuery := columns + query + dates + groupBy

	log.Printf("> params: %v\n", params)

	// transforma []string para []interface
	paramss := make([]interface{}, len(params))
	for i, s := range params {
		paramss[i] = s
	}

	log.Printf("> number of params: %v\n", len(paramss))

	rows, errQuery := r.db.Query(
		context.Background(),
		finalQuery,
		paramss...,
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

		acc := class + ":" + group + ":" + subgroup + ":" + id

		path := vos.Path{
			Account: acc,
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

	syntheticReport, errEntity := vos.NewSyntheticReport(totalCredit, totalDebit, pathsReport, vos.Version(currentVersion))
	if errEntity != nil {
		log.Printf("> error on entity: %v", errEntity)
		return nil, errEntity
	}

	return syntheticReport, nil

}
