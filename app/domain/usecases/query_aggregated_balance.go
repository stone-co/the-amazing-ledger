package usecases

import (
	"context"
	"fmt"

	"github.com/stone-co/the-amazing-ledger/app/domain/vos"
)

func (l *LedgerUseCase) QueryAggregatedBalance(ctx context.Context, query vos.Account) (vos.QueryBalance, error) {
	queryBalance, err := l.repository.QueryAggregatedBalance(ctx, query)
	if err != nil {
		return vos.QueryBalance{}, fmt.Errorf("failed to query aggregated balance: %w", err)
	}

	return queryBalance, nil
}
