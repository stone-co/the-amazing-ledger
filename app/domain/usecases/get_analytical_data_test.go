package usecases

import (
	"context"
	"testing"

	"github.com/stone-co/the-amazing-ledger/app/domain/vos"
	"github.com/stretchr/testify/assert"
)

func TestLedgerUseCase_GetAnalyticalData(t *testing.T) {
	t.Run("Account can be empty", func(t *testing.T) {
		accountPath, err := vos.NewAccountPath("liability:stone:clients")
		assert.Nil(t, err)

		entries := []vos.Entry{}

		useCase := newFakeGetAnalyticalData(entries, nil)
		res, err := useCase.GetAnalyticalData(context.Background(), *accountPath)
		assert.Nil(t, err)
		assert.Equal(t, entries, res)
	})

	t.Run("Account with 2 entries", func(t *testing.T) {
		accountPath, err := vos.NewAccountPath("liability:stone:clients")
		assert.Nil(t, err)

		accountName1 := "liability:stone:clients:user-1"
		accountName2 := "liability:stone:clients:user-2"

		entries := []vos.Entry{
			{
				Account:   accountName1,
				Operation: vos.CreditOperation,
				Amount:    15000,
			},
			{
				Account:   accountName2,
				Operation: vos.CreditOperation,
				Amount:    25000,
			},
		}

		useCase := newFakeGetAnalyticalData(entries, nil)
		res, err := useCase.GetAnalyticalData(context.Background(), *accountPath)
		assert.Nil(t, err)
		assert.Equal(t, entries, res)
	})
}
