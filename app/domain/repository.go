package domain

import (
	"context"

	"github.com/stone-co/the-amazing-ledger/app/domain/entities"
	"github.com/stone-co/the-amazing-ledger/app/domain/vos"
)

type Repository interface {
	CreateTransaction(context.Context, *entities.Transaction) error
	LoadObjectsIntoCache(ctx context.Context, objects *entities.CachedAccounts) (vos.Version, error)
	GetAccountBalance(ctx context.Context, accountName vos.AccountName) (*vos.AccountBalance, error)
	GetAnalyticalData(ctx context.Context, path vos.AccountPath) ([]vos.Statement, error)
}
