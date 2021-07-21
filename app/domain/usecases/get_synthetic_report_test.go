package usecases

import (
	"context"
	"testing"
	"time"

	"github.com/newrelic/go-agent/v3/newrelic"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"

	"github.com/stone-co/the-amazing-ledger/app/domain/probes"
	"github.com/stone-co/the-amazing-ledger/app/domain/vos"
	"github.com/stone-co/the-amazing-ledger/app/tests/mocks"
)

func TestLedgerUseCase_GetSyntheticReport(t *testing.T) {
	t.Run("GetSyntheticReport with startDate and endDate being equals must return the report for the year", func(t *testing.T) {
		query, err := vos.NewAccount("liability.credit_card.invoice.*")
		assert.NoError(t, err)

		totalCredit := int64(2000)
		totalDebit := int64(1000)

		accountPath, err := vos.NewAnalyticalAccount("liability.credit_card.invoice")
		assert.NoError(t, err)

		paths := []vos.Path{{
			Account: accountPath,
			Debit:   1000,
			Credit:  2000,
		}}

		level := 3
		date := time.Now()

		fakeSyntheticReport, err := vos.NewSyntheticReport(totalCredit, totalDebit, paths)
		assert.NoError(t, err)

		mockedRepository := mocks.RepositoryMock{
			GetSyntheticReportFunc: func(ctx context.Context, account vos.Account, level int, startTime, endTime time.Time) (*vos.SyntheticReport, error) {
				return fakeSyntheticReport, nil
			},
		}

		useCase := NewLedgerUseCase(&mockedRepository, probes.NewLedgerProbe(logrus.New(), &newrelic.Application{}))

		got, err := useCase.GetSyntheticReport(context.Background(), query, level, date, date)
		assert.NoError(t, err)
		assert.Equal(t, fakeSyntheticReport.TotalCredit, got.TotalCredit)
		assert.Equal(t, fakeSyntheticReport.TotalDebit, got.TotalDebit)
		assert.Equal(t, len(fakeSyntheticReport.Paths), len(got.Paths))
	})
}
