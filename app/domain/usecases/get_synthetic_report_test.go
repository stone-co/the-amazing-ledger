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
		totalCredit := int64(150)
		totalDebit := int64(130)

		accountPath, err := vos.NewAccountPath("liability.credit_card.invoice")
		assert.Nil(t, err)

		paths := []vos.Path{{
			Account: accountPath,
			Debit:   1000,
			Credit:  2000,
		}}

		fakeSyntheticReport, err := vos.NewSyntheticReport(totalCredit, totalDebit, paths)
		assert.Nil(t, err)

		level := 3
		date := time.Now()

		useCase := newFakeGetSyntheticReport(fakeSyntheticReport, date, nil)
		a, err := useCase.GetSyntheticReport(context.Background(), accountPath, level, date, date)

		assert.Nil(t, err)
		assert.Equal(t, fakeSyntheticReport.TotalDebit, a.TotalDebit)
	})
}
