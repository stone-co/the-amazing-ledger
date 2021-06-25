package usecases

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/stone-co/the-amazing-ledger/app/domain/vos"
)

func TestLedgerUseCase_GetSyntheticReport(t *testing.T) {
	t.Run("GetSyntheticReport with startDate and endDate being equals must return the report for the year", func(t *testing.T) {
		totalCredit := 150
		totalDebit := 130

		version := vos.Version(1)
		accountNameStr := "liability"

		paths := []vos.Path{{
			Account: accountNameStr,
			Debit:   1000,
			Credit:  2000,
		}}

		syntheticReport, err := vos.NewSyntheticReport(totalCredit, totalDebit, paths, version)
		assert.Nil(t, err)

		date := time.Now()

		useCase := newFakeGetSyntheticReport(syntheticReport, date, nil)
		a, err := useCase.GetSyntheticReport(context.Background(), *&accountNameStr, date, date)

		assert.Nil(t, err)
		assert.Equal(t, syntheticReport.TotalDebit, a.TotalDebit)
	})
}
