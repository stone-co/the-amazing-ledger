package usecases

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/stone-co/the-amazing-ledger/app/domain/vos"
)

func TestLedgerUseCase_QueryAggregatedBalance(t *testing.T) {
	t.Run("The Balance is totalCredit subtract by totalDebit", func(t *testing.T) {
		totalCredit := 150
		totalDebit := 130
		expectedBalance := totalCredit - totalDebit

		query, err := vos.NewAccountQuery("liability.stone.clients.*")
		assert.Nil(t, err)

		queryBalance := vos.NewQueryBalance(query, 20)

		useCase := newFakeQueryAggregatedBalance(queryBalance, nil)
		a, err := useCase.QueryAggregatedBalance(context.Background(), query)
		assert.Nil(t, err)
		assert.Equal(t, expectedBalance, a.Balance)
	})
}
