package usecases

import (
	"context"
	"fmt"

	"github.com/stone-co/the-amazing-ledger/app/domain/vos"
)

func (l *LedgerUseCase) GetAccountHistory(ctx context.Context, account vos.AccountPath, fn func(vos.EntryHistory) error) error {
	err := l.repository.GetAccountHistory(ctx, account, fn)
	if err != nil {
		return fmt.Errorf("failed to get account history: %w", err)
	}

	return nil
}
