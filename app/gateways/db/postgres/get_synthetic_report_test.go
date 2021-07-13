package postgres

import (
	"context"
	"encoding/json"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"

	"github.com/stone-co/the-amazing-ledger/app/domain/entities"
	"github.com/stone-co/the-amazing-ledger/app/domain/vos"
)

func TestLedgerRepository_GetSyntheticReport(t *testing.T) {
	event := uint32(4)
	company := "abc"
	competenceDate := time.Now().UTC()
	metadata := json.RawMessage(`{}`)

	r := NewLedgerRepository(pgDocker.DB, logrus.New())
	ctx := context.Background()

	_, err := pgDocker.DB.Exec(ctx, `insert into event (id, name) values (4, 'default4');`)
	assert.NoError(t, err)

	accountBase := "liability.assets"
	accountBaseEmpty := "liability.income"

	acc1Level3, err := vos.NewAccountPath(accountBase + vos.DepthSeparator + "account11")
	assert.NoError(t, err)

	acc2Level3, err := vos.NewAccountPath(accountBase + vos.DepthSeparator + "account22")
	assert.NoError(t, err)

	e1, _ := entities.NewEntry(
		uuid.New(),
		vos.DebitOperation,
		acc1Level3.Name(),
		vos.IgnoreAccountVersion,
		100,
		metadata)

	e2, _ := entities.NewEntry(
		uuid.New(),
		vos.CreditOperation,
		acc2Level3.Name(),
		vos.IgnoreAccountVersion,
		100,
		metadata,
	)

	tx, err := entities.NewTransaction(uuid.New(), event, company, competenceDate, e1, e2)
	assert.NoError(t, err)

	testCases := []struct {
		name        string
		query       string
		level       int
		startTime   time.Time
		endTime     time.Time
		transaction entities.Transaction
		report      vos.SyntheticReport
		err         error
	}{
		{
			name:        "should not get a result because there was no data inserted",
			query:       accountBaseEmpty + vos.DepthSeparator + "*",
			level:       3,
			startTime:   time.Now().UTC(),
			endTime:     time.Now().UTC(),
			transaction: entities.Transaction{},
			report:      vos.SyntheticReport{},
			err:         nil,
		},
		{
			name:        "should get a result",
			query:       accountBase + vos.DepthSeparator + "*",
			level:       3,
			startTime:   time.Now().UTC(),
			endTime:     time.Now().UTC().Add(time.Hour * 1),
			transaction: tx,
			report: vos.SyntheticReport{
				TotalCredit: 100,
				TotalDebit:  100,
			},
			err: nil,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			query, err := vos.NewAccountQuery(tc.query)
			assert.NoError(t, err)

			if tc.transaction.Company != "" {
				err = r.CreateTransaction(ctx, tx)
				assert.NoError(t, err)
			}

			got, err := r.GetSyntheticReport(ctx, query, tc.level, tc.startTime, tc.endTime)
			assert.NoError(t, err)

			if tc.transaction.Company != "" {
				assert.NotEmpty(t, got)
				assert.Equal(t, tc.report.TotalCredit, got.TotalCredit)
				assert.Equal(t, tc.report.TotalDebit, got.TotalDebit)
			} else {
				assert.Empty(t, got)
			}
		})
	}
}
