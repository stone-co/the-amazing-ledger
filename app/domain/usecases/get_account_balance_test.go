package usecases

import (
	"context"
	"testing"

	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"

	"github.com/stone-co/the-amazing-ledger/app"
	"github.com/stone-co/the-amazing-ledger/app/domain/probes"
	"github.com/stone-co/the-amazing-ledger/app/domain/vos"
	"github.com/stone-co/the-amazing-ledger/app/tests/mocks"
	"github.com/stone-co/the-amazing-ledger/app/tests/testdata"
)

func TestLedgerUseCase_GetAccountBalance(t *testing.T) {
	t.Run("should return account balance successfully", func(t *testing.T) {
		accountPath, err := vos.NewAccountPath(testdata.GenerateAccountPath())
		assert.NoError(t, err)

		accountBalance := vos.NewAccountBalance(accountPath, vos.Version(1), 150, 130)
		mockedRepository := &mocks.RepositoryMock{
			GetAccountBalanceFunc: func(ctx context.Context, account vos.AccountPath) (vos.AccountBalance, error) {
				return accountBalance, nil
			},
		}
		usecase := NewLedgerUseCase(mockedRepository, probes.NewLedgerProbe(logrus.New()))

		got, err := usecase.GetAccountBalance(context.Background(), accountPath)
		assert.NoError(t, err)

		assert.Equal(t, accountBalance.Account, got.Account)
		assert.Equal(t, accountBalance.CurrentVersion, got.CurrentVersion)
		assert.Equal(t, accountBalance.TotalCredit, got.TotalCredit)
		assert.Equal(t, accountBalance.TotalDebit, got.TotalDebit)
		assert.Equal(t, accountBalance.Balance(), got.Balance())
	})

	t.Run("should return an error if account does not exist", func(t *testing.T) {
		accountPath, err := vos.NewAccountPath(testdata.GenerateAccountPath())
		assert.NoError(t, err)

		mockedRepository := &mocks.RepositoryMock{
			GetAccountBalanceFunc: func(ctx context.Context, account vos.AccountPath) (vos.AccountBalance, error) {
				return vos.AccountBalance{}, app.ErrAccountNotFound
			},
		}
		usecase := NewLedgerUseCase(mockedRepository, probes.NewLedgerProbe(logrus.New()))

		got, err := usecase.GetAccountBalance(context.Background(), accountPath)
		assert.Empty(t, got)
		assert.ErrorIs(t, err, app.ErrAccountNotFound)
	})
}
