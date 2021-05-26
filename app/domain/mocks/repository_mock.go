package mocks

import (
	"context"

	"github.com/stone-co/the-amazing-ledger/app/domain"
	"github.com/stone-co/the-amazing-ledger/app/domain/entities"
	"github.com/stone-co/the-amazing-ledger/app/domain/vos"
)

var _ domain.Repository = &Repository{}

type Repository struct {
	OnCreateTransaction           func(context.Context, entities.Transaction) error
	OnLoadObjectsIntoCache        func(ctx context.Context, cachedAccounts *entities.CachedAccounts) (vos.Version, error)
	OnGetAccountBalance           func(ctx context.Context, account vos.AccountPath) (*vos.AccountBalance, error)
	OnGetAccountBalanceAggregated func(ctx context.Context, account vos.AccountPath) (*vos.AccountBalance, error)
	OnGetAnalyticalData           func(ctx context.Context, query vos.AccountQuery, fn func(vos.Statement) error) error
	OnGetAccountHistory           func(ctx context.Context, account vos.AccountPath, fn func(vos.EntryHistory) error) error
}

func (s Repository) CreateTransaction(ctx context.Context, transaction entities.Transaction) error {
	return s.OnCreateTransaction(ctx, transaction)
}

func (s Repository) LoadObjectsIntoCache(ctx context.Context, cachedAccounts *entities.CachedAccounts) (vos.Version, error) {
	return s.OnLoadObjectsIntoCache(ctx, cachedAccounts)
}

func (s Repository) GetAccountBalance(ctx context.Context, account vos.AccountPath) (*vos.AccountBalance, error) {
	return s.OnGetAccountBalance(ctx, account)
}

func (s Repository) GetAccountBalanceAggregated(ctx context.Context, account vos.AccountPath) (*vos.AccountBalance, error) {
	return s.OnGetAccountBalanceAggregated(ctx, account)
}

func (s Repository) GetAnalyticalData(ctx context.Context, query vos.AccountQuery, fn func(vos.Statement) error) error {
	return s.OnGetAnalyticalData(ctx, query, fn)
}

func (s Repository) GetAccountHistory(ctx context.Context, account vos.AccountPath, fn func(vos.EntryHistory) error) error {
	return s.OnGetAccountHistory(ctx, account, fn)
}
