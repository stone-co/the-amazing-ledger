package domain

import (
	"context"

	"github.com/stone-co/the-amazing-ledger/app/domain/entities"
	"github.com/stone-co/the-amazing-ledger/app/domain/vos"
)

type Repository interface {
	CreateTransaction(context.Context, entities.Transaction) error
	LoadObjectsIntoCache(ctx context.Context, objects *entities.CachedAccounts) (vos.Version, error)
	GetAccountBalance(ctx context.Context, account vos.AccountPath) (vos.AccountBalance, error)
	GetAnalyticalData(ctx context.Context, query vos.AccountQuery, fn func(vos.Statement) error) error
	GetAccountHistory(ctxt context.Context, account vos.AccountPath, fn func(vos.EntryHistory) error) error
}
