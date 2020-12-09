package domain

import (
	"context"

	"github.com/google/uuid"
	"github.com/stone-co/the-amazing-ledger/app/domain/entities"
	"github.com/stone-co/the-amazing-ledger/app/domain/vo"
)

type UseCase interface {
	CreateTransaction(ctx context.Context, id uuid.UUID, entries []entities.Entry) error
	LoadObjectsIntoCache(ctx context.Context) error
	GetAccountBalance(ctx context.Context, accountName vo.AccountName) (*vo.AccountBalance, error)
}
