package domain

import (
	"context"

	"github.com/stone-co/the-amazing-ledger/app/domain/entities"
	"github.com/stone-co/the-amazing-ledger/app/domain/vo"
)

type Repository interface {
	CreateTransaction(context.Context, *entities.Transaction) error
	LoadObjectsIntoCache(ctx context.Context, objects *entities.CachedAccounts) (vo.Version, error)
	GetAccountBalance(ctx context.Context, accountName vo.AccountName) (*vo.AccountBalance, error)
}
