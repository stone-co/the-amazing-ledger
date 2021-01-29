package usecases

import (
	"context"

	"github.com/stone-co/the-amazing-ledger/app/domain/vos"
)

func (l *LedgerUseCase) GetAccountHistory(ctx context.Context, accountName vos.AccountName, fn func(vos.EntryHistory) error) error {
	return l.repository.GetAccountHistory(ctx, accountName, fn)
}
