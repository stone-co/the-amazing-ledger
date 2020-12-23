package mocks

import (
	"context"

	"github.com/stone-co/the-amazing-ledger/app/domain"
	"github.com/stone-co/the-amazing-ledger/app/domain/entities"
	"github.com/stone-co/the-amazing-ledger/app/domain/vos"
)

var _ domain.Repository = &Repository{}

type Repository struct {
	OnCreateTransaction    func(context.Context, *entities.Transaction) error
	OnLoadObjectsIntoCache func(ctx context.Context, cachedAccounts *entities.CachedAccounts) (vos.Version, error)
	OnGetAccountBalance    func(ctx context.Context, accountName vos.AccountName) (*vos.AccountBalance, error)
	OnGetAnalyticalData    func(ctx context.Context, path vos.AccountPath) ([]vos.Statement, error)
}

func (s Repository) CreateTransaction(ctx context.Context, transaction *entities.Transaction) error {
	return s.OnCreateTransaction(ctx, transaction)
}

func (s Repository) LoadObjectsIntoCache(ctx context.Context, cachedAccounts *entities.CachedAccounts) (vos.Version, error) {
	return s.OnLoadObjectsIntoCache(ctx, cachedAccounts)
}

func (s Repository) GetAccountBalance(ctx context.Context, accountName vos.AccountName) (*vos.AccountBalance, error) {
	return s.OnGetAccountBalance(ctx, accountName)
}

func (s Repository) GetAnalyticalData(ctx context.Context, path vos.AccountPath) ([]vos.Statement, error) {
	return s.OnGetAnalyticalData(ctx, path)
}
