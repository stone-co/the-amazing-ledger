package mocks

import (
	"context"

	"github.com/stone-co/the-amazing-ledger/app/domain"
	"github.com/stone-co/the-amazing-ledger/app/domain/entities"
	"github.com/stone-co/the-amazing-ledger/app/domain/vo"
)

var _ domain.Repository = &Repository{}

type Repository struct {
	OnCreateTransaction    func(context.Context, *entities.Transaction) error
	OnLoadObjectsIntoCache func(ctx context.Context, cachedAccounts *entities.CachedAccounts) (vo.Version, error)
	OnGetAccountBalance    func(ctx context.Context, accountName vo.AccountName) (*vo.AccountBalance, error)
}

func (s Repository) CreateTransaction(ctx context.Context, transaction *entities.Transaction) error {
	return s.OnCreateTransaction(ctx, transaction)
}

func (s Repository) LoadObjectsIntoCache(ctx context.Context, cachedAccounts *entities.CachedAccounts) (vo.Version, error) {
	return s.OnLoadObjectsIntoCache(ctx, cachedAccounts)
}

func (s Repository) GetAccountBalance(ctx context.Context, accountName vo.AccountName) (*vo.AccountBalance, error) {
	return s.OnGetAccountBalance(ctx, accountName)
}
