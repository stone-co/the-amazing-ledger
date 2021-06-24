package usecases

import (
	"context"

	"github.com/stone-co/the-amazing-ledger/app/domain/vos"
)

func (l *LedgerUseCase) QueryAggregatedBalance(ctx context.Context, query vos.AccountQuery) (vos.QueryBalance, error) {
	queryBalance, err := l.repository.QueryAggregatedBalance(ctx, query)
	if err != nil {
		return vos.QueryBalance{}, err
	}

	return queryBalance, nil
}
