package ledger

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/stone-co/the-amazing-ledger/pkg/command-handler/domain/ledger/entities"
)

type TransactionsMock struct {
	OnCreateTransaction func(ctx context.Context, id uuid.UUID, createdAt time.Time, entries []entities.Entry) error
}

func (m TransactionsMock) CreateTransaction(ctx context.Context, id uuid.UUID, createdAt time.Time, entries []entities.Entry) error {
	return m.OnCreateTransaction(ctx, id, createdAt, entries)
}
