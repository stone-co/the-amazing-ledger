package usecases

import (
	"context"

	"github.com/stone-co/the-amazing-ledger/app/domain/vos"
)

func (l *LedgerUseCase) GetAnalyticalData(ctx context.Context, path vos.AccountPath, fn func(vos.Statement) error) error {

	return l.repository.GetAnalyticalData(ctx, path, fn)
}
