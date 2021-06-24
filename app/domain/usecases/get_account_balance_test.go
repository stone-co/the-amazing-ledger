package usecases

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/stone-co/the-amazing-ledger/app/domain/vos"
)

func TestLedgerUseCase_GetAccountBalance(t *testing.T) {
	t.Run("The Balance is totalCredit subtract by totalDebit", func(t *testing.T) {
		totalCredit := 150
		totalDebit := 130
		expectedBalance := totalCredit - totalDebit

		account, err := vos.NewAccountPath("liability.stone.clients.user_1")
		assert.Nil(t, err)

		accountBalance := vos.NewAccountBalance(account, 3, totalCredit, totalDebit)

		useCase := newFakeGetAccountBalance(accountBalance, nil)
		a, err := useCase.GetAccountBalance(context.Background(), accountBalance.Account)
		assert.Nil(t, err)
		assert.Equal(t, accountBalance.TotalCredit, a.TotalCredit)
		assert.Equal(t, accountBalance.TotalDebit, a.TotalDebit)
		assert.Equal(t, expectedBalance, a.Balance())
	})

	t.Run("The max version for account path must be version in account balance", func(t *testing.T) {
		expectedVersion := vos.Version(5)

		account, err := vos.NewAccountPath("liability.stone.clients.user_1")
		assert.Nil(t, err)

		accountBalance := vos.NewAccountBalance(account, expectedVersion, 0, 0)

		useCase := newFakeGetAccountBalance(accountBalance, nil)
		a, err := useCase.GetAccountBalance(context.Background(), accountBalance.Account)

		assert.Nil(t, err)
		assert.Equal(t, expectedVersion, a.CurrentVersion)
	})
}
