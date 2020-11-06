package usecase

import (
	"context"
	"testing"

	"github.com/stone-co/the-amazing-ledger/pkg/command-handler/domain/ledger/entities"
	"github.com/stretchr/testify/assert"
)

func TestLedgerUseCase_GetAccountInfo(t *testing.T) {
	t.Run("The Balance is totalCredit subtract by totalDebit", func(t *testing.T) {
		totalCredit := 150
		totalDebit := 130
		expectedBalance := totalCredit - totalDebit

		accountInfo := entities.NewAccountInfo("liability:stone:clients:user-1", 3, totalCredit, totalDebit)

		useCase := newFakeGetAccountInfo(accountInfo, nil)
		a, err := useCase.GetAccountInfo(context.Background(), accountInfo.AccountID)
		assert.Nil(t, err)
		assert.Equal(t, a.TotalCredit, accountInfo.TotalCredit)
		assert.Equal(t, a.TotalDebit, accountInfo.TotalDebit)
		assert.Equal(t, a.Balance(), expectedBalance)
	})

	t.Run("The max version for account id shoud be version in account info", func(t *testing.T) {
		expectedVersion := entities.Version(5)

		accountInfo := entities.NewAccountInfo("liability:stone:clients:user-1", expectedVersion, 0, 0)

		useCase := newFakeGetAccountInfo(accountInfo, nil)
		a, err := useCase.GetAccountInfo(context.Background(), accountInfo.AccountID)

		assert.Nil(t, err)
		assert.Equal(t, a.CurrentVersion, expectedVersion)
	})

}
