package ledger

import (
	"context"

	"github.com/google/uuid"
	"github.com/stone-co/the-amazing-ledger/pkg/command-handler/domain/ledger/entities"
)

type TransactionsUseCase interface {
	CreateTransaction(ctx context.Context, id uuid.UUID, entries []entities.Entry) error
	LoadObjectsIntoCache(ctx context.Context) error
	GetAccountInfo(ctx context.Context, accountID string) (*entities.AccountInfo, error)
}
