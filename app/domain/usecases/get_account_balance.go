package usecases

import (
	"context"

	"github.com/stone-co/the-amazing-ledger/app/domain/vos"
)

func (l *LedgerUseCase) GetAccountBalance(ctx context.Context, accountName vos.AccountName) (*vos.AccountBalance, error) {
	if accountName.Suffix != "*" {
		return l.getAccountBalance(ctx, accountName)
	}

	return l.getAccountBalanceAggregated(ctx, accountName)
}

func (l *LedgerUseCase) getAccountBalance(ctx context.Context, accountName vos.AccountName) (*vos.AccountBalance, error) {
	accountBalance, err := l.repository.GetAccountBalance(ctx, accountName)
	if err != nil {
		return nil, err
	}

	return accountBalance, nil
}

func (l *LedgerUseCase) getAccountBalanceAggregated(ctx context.Context, accountName vos.AccountName) (*vos.AccountBalance, error) {
	accountBalance, err := l.repository.GetAccountBalanceAggregated(ctx, accountName)
	if err != nil {
		return nil, err
	}

	return accountBalance, nil
}
