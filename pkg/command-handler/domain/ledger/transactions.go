package ledger

import (
	"context"

	"github.com/google/uuid"
	"github.com/stone-co/the-amazing-ledger/pkg/command-handler/domain/ledger/entities"
)

type TransactionsUseCase interface {
	CreateTransaction(ctx context.Context, id uuid.UUID, entries []entities.Entry) error
	LoadObjectsIntoCache(ctx context.Context) error
	GetAccountBalance(ctx context.Context, accountName entities.AccountName) (*entities.AccountBalance, error)
}
