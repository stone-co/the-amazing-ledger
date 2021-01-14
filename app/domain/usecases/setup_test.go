package usecases

import (
	"context"

	"github.com/sirupsen/logrus"

	"github.com/stone-co/the-amazing-ledger/app/domain/entities"
	"github.com/stone-co/the-amazing-ledger/app/domain/mocks"
	"github.com/stone-co/the-amazing-ledger/app/domain/vos"
)

func newFakeCreateTransactionUseCase(result error) *LedgerUseCase {
	log := logrus.New()

	mockRepository := &mocks.Repository{}
	mockRepository.OnCreateTransaction = func(context.Context, *entities.Transaction) error {
		return result
	}

	return NewLedgerUseCase(log, mockRepository)
}

func newFakeLoadObjectsIntoCacheUseCase(maxVersion vos.Version, result error) *LedgerUseCase {
	log := logrus.New()

	mockRepository := &mocks.Repository{
		OnLoadObjectsIntoCache: func(ctx context.Context, cachedAccounts *entities.CachedAccounts) (vos.Version, error) {
			return maxVersion, result
		},
	}

	return NewLedgerUseCase(log, mockRepository)
}

func newFakeGetAccountBalance(accountBalance *vos.AccountBalance, result error) *LedgerUseCase {
	log := logrus.New()

	mockRepository := &mocks.Repository{
		OnGetAccountBalance: func(ctx context.Context, accountName vos.AccountName) (*vos.AccountBalance, error) {
			return accountBalance, result
		},
	}

	return NewLedgerUseCase(log, mockRepository)
}

func newFakeGetAnalyticalData(entries []vos.Statement, result error) *LedgerUseCase {
	log := logrus.New()

	mockRepository := &mocks.Repository{
		OnGetAnalyticalData: func(ctx context.Context, path vos.AccountPath, fn func(vos.Statement) error) error {
			for _, entry := range entries {
				if err := fn(entry); err != nil {
					return err
				}
			}
			return result
		},
	}

	return NewLedgerUseCase(log, mockRepository)
}

func newFakeGetAccountHistory(accountHistory *vos.AccountHistory, result error) *LedgerUseCase {
	log := logrus.New()

	mockRepository := &mocks.Repository{
		OnGetAccountHistory: func(ctx context.Context, accountName vos.AccountName) (*vos.AccountHistory, error) {
			return accountHistory, result
		},
	}

	return NewLedgerUseCase(log, mockRepository)
}
