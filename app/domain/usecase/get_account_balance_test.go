package usecase

import (
	"context"
	"testing"

	"github.com/stone-co/the-amazing-ledger/app/domain/vo"
	"github.com/stretchr/testify/assert"
)

func TestLedgerUseCase_GetAccountBalance(t *testing.T) {
	t.Run("The Balance is totalCredit subtract by totalDebit", func(t *testing.T) {
		totalCredit := 150
		totalDebit := 130
		expectedBalance := totalCredit - totalDebit

		accountName, err := vo.NewAccountName("liability:stone:clients:user-1")
		assert.Nil(t, err)

		accountBalance := vo.NewAccountBalance(*accountName, 3, totalCredit, totalDebit)

		useCase := newFakeGetAccountBalance(accountBalance, nil)
		a, err := useCase.GetAccountBalance(context.Background(), accountBalance.AccountName)
		assert.Nil(t, err)
		assert.Equal(t, accountBalance.TotalCredit, a.TotalCredit)
		assert.Equal(t, accountBalance.TotalDebit, a.TotalDebit)
		assert.Equal(t, expectedBalance, a.Balance())
	})

	t.Run("The max version for account path must be version in account balance", func(t *testing.T) {
		expectedVersion := vo.Version(5)

		accountName, err := vo.NewAccountName("liability:stone:clients:user-1")
		assert.Nil(t, err)

		accountBalance := vo.NewAccountBalance(*accountName, expectedVersion, 0, 0)

		useCase := newFakeGetAccountBalance(accountBalance, nil)
		a, err := useCase.GetAccountBalance(context.Background(), accountBalance.AccountName)

		assert.Nil(t, err)
		assert.Equal(t, expectedVersion, a.CurrentVersion)
	})
}
