package domain

import (
	"context"

	"github.com/stone-co/the-amazing-ledger/app/domain/entities"
	"github.com/stone-co/the-amazing-ledger/app/domain/vos"
)

type UseCase interface {
	CreateTransaction(context.Context, entities.Transaction) error
	LoadObjectsIntoCache(context.Context) error
	GetAccountBalance(context.Context, vos.AccountPath) (*vos.AccountBalance, error)
	GetAnalyticalData(context.Context, vos.AccountQuery, func(vos.Statement) error) error
	GetAccountHistory(context.Context, vos.AccountPath, func(vos.EntryHistory) error) error
}
