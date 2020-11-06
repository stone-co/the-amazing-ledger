package usecase

import (
	"context"

	"github.com/sirupsen/logrus"

	"github.com/stone-co/the-amazing-ledger/pkg/command-handler/domain/ledger"
	"github.com/stone-co/the-amazing-ledger/pkg/command-handler/domain/ledger/entities"
)

func newFakeCreateTransactionUseCase(result error) *LedgerUseCase {
	log := logrus.New()

	mockRepository := &ledger.RepositoryMock{}
	mockRepository.OnCreateTransaction = func(context.Context, *entities.Transaction) error {
		return result
	}

	return NewLedgerUseCase(log, mockRepository)
}

func newFakeLoadObjectsIntoCacheUseCase(maxVersion entities.Version, result error) *LedgerUseCase {
	log := logrus.New()

	mockRepository := &ledger.RepositoryMock{
		OnLoadObjectsIntoCache: func(ctx context.Context, cachedAccounts *entities.CachedAccounts) (entities.Version, error) {
			return maxVersion, result
		},
	}

	return NewLedgerUseCase(log, mockRepository)
}

func newFakeGetAccountInfo(accountInfo *entities.AccountInfo, result error) *LedgerUseCase {
	log := logrus.New()

	mockRepository := &ledger.RepositoryMock{
		OnGetAccountInfo: func(ctx context.Context, accountID string) (*entities.AccountInfo, error) {
			return accountInfo, result
		},
	}

	return NewLedgerUseCase(log, mockRepository)
}
