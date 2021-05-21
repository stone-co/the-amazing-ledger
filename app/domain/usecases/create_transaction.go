package usecases

import (
	"context"

	"github.com/stone-co/the-amazing-ledger/app"
	"github.com/stone-co/the-amazing-ledger/app/domain/entities"
	"github.com/stone-co/the-amazing-ledger/app/domain/vos"
)

func (l *LedgerUseCase) CreateTransaction(ctx context.Context, transaction entities.Transaction) error {
	accounts := make([]*entities.CachedAccountInfo, 0, len(transaction.Entries))

	for _, entry := range transaction.Entries {
		account, err := l.checkVersion(entry)
		if err != nil {
			return err
		}
		accounts = append(accounts, account)
	}

	for i := range transaction.Entries {
		transaction.Entries[i].Version = l.lastVersion.Next()
	}

	if err := l.repository.CreateTransaction(ctx, transaction); err != nil {
		return err
	}

	for i := range accounts {
		accounts[i].CurrentVersion = transaction.Entries[i].Version
	}

	return nil
}

func (l *LedgerUseCase) checkVersion(entry entities.Entry) (*entities.CachedAccountInfo, error) {
	account := l.cachedAccounts.LoadOrStore(entry.Account.Name())

	account.Lock()
	defer account.Unlock()

	if entry.Version == vos.AnyAccountVersion {
		return account, nil
	}

	if entry.Version != account.CurrentVersion {
		return nil, app.ErrInvalidVersion
	}

	return account, nil
}
