package domain

import (
	"context"

	"github.com/google/uuid"
	"github.com/stone-co/the-amazing-ledger/app/domain/entities"
	"github.com/stone-co/the-amazing-ledger/app/domain/vos"
)

type UseCase interface {
	CreateTransaction(ctx context.Context, id uuid.UUID, entries []entities.Entry) error
	LoadObjectsIntoCache(ctx context.Context) error
	GetAccountBalance(ctx context.Context, accountName vos.AccountName) (*vos.AccountBalance, error)
	GetAnalyticalData(ctx context.Context, path vos.AccountPath) ([]vos.Entry, error)
}
