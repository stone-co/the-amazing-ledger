package usecases

import (
	"context"
	"fmt"

	"github.com/stone-co/the-amazing-ledger/app/domain/vos"
)

func (l *LedgerUseCase) GetAnalyticalData(ctx context.Context, query vos.AccountQuery, fn func(vos.Statement) error) error {
	err := l.repository.GetAnalyticalData(ctx, query, fn)
	if err != nil {
		return fmt.Errorf("failed to get analytical data: %w", err)
	}

	return nil
}
