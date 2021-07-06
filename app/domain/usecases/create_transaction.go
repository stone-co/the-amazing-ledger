package usecases

import (
	"context"
	"fmt"

	"github.com/stone-co/the-amazing-ledger/app/domain/entities"
)

func (l *LedgerUseCase) CreateTransaction(ctx context.Context, transaction entities.Transaction) error {
	err := l.repository.CreateTransaction(ctx, transaction)
	if err != nil {
		return fmt.Errorf("failed to create transaction: %w", err)
	}

	return nil
}
