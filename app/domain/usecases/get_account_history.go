package usecases

import (
	"context"

	"github.com/stone-co/the-amazing-ledger/app/domain/vos"
)

func (l *LedgerUseCase) GetAccountHistory(ctx context.Context, accountName vos.AccountName) (*vos.AccountHistory, error) {
	accountHistory, err := l.repository.GetAccountHistory(ctx, accountName)
	if err != nil {
		return nil, err
	}

	return accountHistory, nil
}
