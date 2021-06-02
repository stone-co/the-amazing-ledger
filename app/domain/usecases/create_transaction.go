package usecases

import (
	"context"

	"github.com/stone-co/the-amazing-ledger/app/domain/entities"
)

func (l *LedgerUseCase) CreateTransaction(ctx context.Context, transaction entities.Transaction) error {
	return l.repository.CreateTransaction(ctx, transaction)
}
