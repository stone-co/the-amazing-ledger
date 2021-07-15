package domain

import (
	"context"
	"time"

	"github.com/stone-co/the-amazing-ledger/app/domain/entities"
	"github.com/stone-co/the-amazing-ledger/app/domain/vos"
	"github.com/stone-co/the-amazing-ledger/app/pagination"
)

type Repository interface {
	CreateTransaction(context.Context, entities.Transaction) error
	GetAccountBalance(ctx context.Context, account vos.AccountPath) (vos.AccountBalance, error)
	QueryAggregatedBalance(ctx context.Context, account vos.AccountQuery) (vos.QueryBalance, error)
	GetSyntheticReport(ctx context.Context, query vos.AccountQuery, level int, startTime time.Time, endTime time.Time) (*vos.SyntheticReport, error)
	ListAccountEntries(context.Context, vos.AccountEntryRequest) ([]vos.AccountEntry, pagination.Cursor, error)
}
