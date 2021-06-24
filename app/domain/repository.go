package domain

import (
	"context"

	"github.com/stone-co/the-amazing-ledger/app/domain/entities"
	"github.com/stone-co/the-amazing-ledger/app/domain/vos"
)

type Repository interface {
	CreateTransaction(context.Context, entities.Transaction) error
	GetAccountBalance(ctx context.Context, account vos.AccountPath) (vos.AccountBalance, error)
	QueryAggregatedBalance(ctx context.Context, account vos.AccountQuery) (vos.QueryBalance, error)
	GetAnalyticalData(ctx context.Context, query vos.AccountQuery, fn func(vos.Statement) error) error
	GetAccountHistory(ctxt context.Context, account vos.AccountPath, fn func(vos.EntryHistory) error) error
}
