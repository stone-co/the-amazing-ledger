package usecases

import (
	"context"
	"fmt"

	"github.com/stone-co/the-amazing-ledger/app/domain/vos"
)

func (l *LedgerUseCase) GetAccountBalance(ctx context.Context, account vos.Account) (vos.AccountBalance, error) {
	accountBalance, err := l.repository.GetAccountBalance(ctx, account)
	if err != nil {
		return vos.AccountBalance{}, fmt.Errorf("failed to get account balance: %w", err)
	}

	return accountBalance, nil
}
