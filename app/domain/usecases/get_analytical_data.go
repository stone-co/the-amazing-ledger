package usecases

import (
	"context"

	"github.com/stone-co/the-amazing-ledger/app/domain/vos"
)

func (l *LedgerUseCase) GetAnalyticalData(ctx context.Context, path vos.AccountPath) ([]vos.Entry, error) {

	accountBalance, err := l.repository.GetAnalyticalData(ctx, path)
	if err != nil {
		return nil, err
	}

	return accountBalance, nil
}
