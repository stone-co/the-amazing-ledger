package domain

import (
	"context"
	"time"

	"github.com/stone-co/the-amazing-ledger/app/domain/entities"
	"github.com/stone-co/the-amazing-ledger/app/domain/vos"
)

type UseCase interface {
	CreateTransaction(context.Context, entities.Transaction) error
	GetAccountBalance(context.Context, vos.AccountPath) (vos.AccountBalance, error)
	QueryAggregatedBalance(context.Context, vos.AccountQuery) (vos.QueryBalance, error)
	GetAnalyticalData(context.Context, vos.AccountQuery, func(vos.Statement) error) error
	GetSyntheticReport(context.Context, vos.AccountQuery, int, time.Time, time.Time) (*vos.SyntheticReport, error)
	ListAccountEntries(context.Context, vos.AccountEntryRequest) (vos.AccountEntryResponse, error)
}
