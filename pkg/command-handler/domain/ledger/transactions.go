package ledger

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/stone-co/the-amazing-ledger/pkg/command-handler/domain/ledger/entities"
)

type TransactionsUseCase interface {
	CreateTransaction(ctx context.Context, id uuid.UUID, createdAt time.Time, entries []entities.Entry) error
}
