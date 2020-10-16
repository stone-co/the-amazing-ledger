package ledger

import (
	"context"

	"github.com/google/uuid"
	"github.com/stone-co/the-amazing-ledger/pkg/command-handler/domain/ledger/entities"
)

type TransactionsMock struct {
	OnCreateTransaction    func(ctx context.Context, id uuid.UUID, entries []entities.Entry) error
	OnLoadObjectsIntoCache func(ctx context.Context) error
}

func (m TransactionsMock) CreateTransaction(ctx context.Context, id uuid.UUID, entries []entities.Entry) error {
	return m.OnCreateTransaction(ctx, id, entries)
}

func (m TransactionsMock) LoadObjectsIntoCache(ctx context.Context) error {
	return m.OnLoadObjectsIntoCache(ctx)
}
