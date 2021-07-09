package postgres

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"

	"github.com/stone-co/the-amazing-ledger/app/domain/entities"
	"github.com/stone-co/the-amazing-ledger/app/domain/vos"
)

func TestLedgerRepository_GetSyntheticReport(t *testing.T) {
	event := uint32(1)
	company := "abc"
	competenceDate := time.Now().UTC()

	r := NewLedgerRepository(pgDocker.DB, logrus.New())
	ctx := context.Background()

	//
	// preparation
	//

	// dependencies
	_, err := pgDocker.DB.Exec(ctx, `insert into event (name) values ('defaults');`)
	assert.NoError(t, err)

	// data
	accountBase := "liability.assets"
	accountBaseEmpty := "liability.income"

	acc1Level3, err := vos.NewAccountPath(accountBase + vos.DepthSeparator + "account11")
	assert.NoError(t, err)

	acc2Level3, err := vos.NewAccountPath(accountBase + vos.DepthSeparator + "account22")
	assert.NoError(t, err)

	//
	// exec
	//

	// TEST_01.A: all params & no results because it is not inserted
	fmt.Println("========== 01.A ==========>")
	query, err := vos.NewAccountQuery(accountBaseEmpty + vos.DepthSeparator + "*")
	assert.NoError(t, err)

	level := 3
	startTime := time.Now().UTC()
	endTime := time.Now().UTC()

	report, err := r.GetSyntheticReport(ctx, query, level, startTime, endTime)
	assert.Empty(t, report)
	assert.True(t, err == nil)

	// TEST_01.B: all params & no results because wrong time range
	fmt.Println("========== 01.B ==========>")
	query, err = vos.NewAccountQuery(accountBase + vos.DepthSeparator + "*")
	assert.NoError(t, err)

	// insert transaction 01
	e1, _ := entities.NewEntry(
		uuid.New(),
		vos.DebitOperation,
		acc1Level3.Name(),
		vos.NextAccountVersion,
		100,
	)
	e2, _ := entities.NewEntry(
		uuid.New(),
		vos.CreditOperation,
		acc2Level3.Name(),
		vos.IgnoreAccountVersion,
		100,
	)

	tx, err := entities.NewTransaction(uuid.New(), event, company, competenceDate, e1, e2)
	assert.NoError(t, err)

	err = r.CreateTransaction(ctx, tx)
	assert.NoError(t, err)
	// end insert transaction 01

	report, err = r.GetSyntheticReport(ctx, query, level, startTime, endTime)
	assert.Empty(t, report)
	assert.True(t, err == nil)

	// TEST_02.A: all params & results
	fmt.Println("========== 02.A ==========>")
	endTime = time.Now().UTC()
	fmt.Println("---")
	fmt.Println(startTime)
	fmt.Println(".")
	fmt.Println(endTime)
	fmt.Println("---")
	fmt.Println()
	report, err = r.GetSyntheticReport(ctx, query, level, startTime, endTime)
	assert.NotEmpty(t, report)
	assert.NoError(t, err)
	assert.Equal(t, int64(100), report.TotalCredit)
	assert.Equal(t, int64(100), report.TotalDebit)
}
