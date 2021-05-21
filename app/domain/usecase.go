package domain

import (
	"context"

	"github.com/stone-co/the-amazing-ledger/app/domain/entities"
	"github.com/stone-co/the-amazing-ledger/app/domain/vos"
)

type UseCase interface {
	CreateTransaction(ctx context.Context, transaction entities.Transaction) error
	LoadObjectsIntoCache(ctx context.Context) error
	GetAccountBalance(ctx context.Context, accountName vos.AccountName) (*vos.AccountBalance, error)
	GetAnalyticalData(ctx context.Context, path vos.AccountPath, fn func(vos.Statement) error) error
	GetAccountHistory(ctx context.Context, accountName vos.AccountName, fn func(vos.EntryHistory) error) error
}
