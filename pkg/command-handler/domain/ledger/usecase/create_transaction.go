package usecase

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/stone-co/the-amazing-ledger/pkg/command-handler/domain/ledger/entities"
)

func (l LedgerUseCase) CreateTransaction(ctx context.Context, id uuid.UUID, createdAt time.Time, entries []entities.Entry) error {
	transaction, err := entities.NewTransaction(id, createdAt, entries...)
	if err != nil {
		return err
	}

	if err := l.repository.CreateTransaction(ctx, transaction); err != nil {
		return fmt.Errorf("can't create transaction: %s", err.Error())
	}

	return nil
}
