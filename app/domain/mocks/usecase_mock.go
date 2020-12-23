package mocks

import (
	"context"

	"github.com/google/uuid"
	"github.com/stone-co/the-amazing-ledger/app/domain"
	"github.com/stone-co/the-amazing-ledger/app/domain/entities"
	"github.com/stone-co/the-amazing-ledger/app/domain/vos"
)

var _ domain.UseCase = &UseCase{}

type UseCase struct {
	OnCreateTransaction    func(ctx context.Context, id uuid.UUID, entries []entities.Entry) error
	OnLoadObjectsIntoCache func(ctx context.Context) error
	OnGetAccountBalance    func(ctx context.Context, accountName vos.AccountName) (*vos.AccountBalance, error)
	OnGetAnalyticalData    func(ctx context.Context, path vos.AccountPath) ([]vos.Statement, error)
}

func (m UseCase) CreateTransaction(ctx context.Context, id uuid.UUID, entries []entities.Entry) error {
	return m.OnCreateTransaction(ctx, id, entries)
}

func (m UseCase) LoadObjectsIntoCache(ctx context.Context) error {
	return m.OnLoadObjectsIntoCache(ctx)
}

func (m UseCase) GetAccountBalance(ctx context.Context, accountName vos.AccountName) (*vos.AccountBalance, error) {
	return m.OnGetAccountBalance(ctx, accountName)
}

func (m UseCase) GetAnalyticalData(ctx context.Context, path vos.AccountPath) ([]vos.Statement, error) {
	return m.OnGetAnalyticalData(ctx, path)
}

func SuccessfulTransactionMock() UseCase {
	return UseCase{
		OnCreateTransaction: func(ctx context.Context, id uuid.UUID, entries []entities.Entry) error {
			return nil
		},
		OnLoadObjectsIntoCache: func(ctx context.Context) error {
			return nil
		},
		OnGetAccountBalance: func(ctx context.Context, accountName vos.AccountName) (*vos.AccountBalance, error) {
			return &vos.AccountBalance{
				AccountName:    vos.AccountName{},
				CurrentVersion: 0,
				TotalCredit:    0,
				TotalDebit:     0,
			}, nil
		},
		OnGetAnalyticalData: func(ctx context.Context, accountName vos.AccountPath) ([]vos.Statement, error) {
			return []vos.Statement{}, nil
		},
	}
}
