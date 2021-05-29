package usecases

import (
	"context"

	"github.com/stone-co/the-amazing-ledger/app/domain/vos"
)

func (l *LedgerUseCase) GetAccountBalance(ctx context.Context, account vos.AccountPath) (vos.AccountBalance, error) {
	return l.getAccountBalance(ctx, account)
}

func (l *LedgerUseCase) getAccountBalance(ctx context.Context, account vos.AccountPath) (vos.AccountBalance, error) {
	accountBalance, err := l.repository.GetAccountBalance(ctx, account)
	if err != nil {
		return vos.AccountBalance{}, err
	}

	return accountBalance, nil
}
