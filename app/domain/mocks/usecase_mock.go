package mocks

import (
	"context"
	"time"

	"github.com/stone-co/the-amazing-ledger/app/domain"
	"github.com/stone-co/the-amazing-ledger/app/domain/entities"
	"github.com/stone-co/the-amazing-ledger/app/domain/vos"
)

var _ domain.UseCase = &UseCase{}

type UseCase struct {
	OnCreateTransaction  func(ctx context.Context, transaction entities.Transaction) error
	OnGetAccountBalance  func(ctx context.Context, account vos.AccountPath) (vos.AccountBalance, error)
	OnGetAnalyticalData  func(ctx context.Context, query vos.AccountQuery, fn func(vos.Statement) error) error
	OnGetAccountHistory  func(ctx context.Context, account vos.AccountPath, fn func(vos.EntryHistory) error) error
	OnGetSyntheticReport func(ctx context.Context, accountPath vos.AccountPath, level int, startTime time.Time, endTime time.Time) (*vos.SyntheticReport, error)
}

func (m UseCase) CreateTransaction(ctx context.Context, transaction entities.Transaction) error {
	return m.OnCreateTransaction(ctx, transaction)
}

func (m UseCase) GetAccountBalance(ctx context.Context, account vos.AccountPath) (vos.AccountBalance, error) {
	return m.OnGetAccountBalance(ctx, account)
}

func (m UseCase) GetAnalyticalData(ctx context.Context, query vos.AccountQuery, fn func(vos.Statement) error) error {
	return m.OnGetAnalyticalData(ctx, query, fn)
}

func (m UseCase) GetAccountHistory(ctx context.Context, account vos.AccountPath, fn func(vos.EntryHistory) error) error {
	return m.OnGetAccountHistory(ctx, account, fn)
}

func (m UseCase) GetSyntheticReport(ctx context.Context, accountPath vos.AccountPath, level int, startTime time.Time, endTime time.Time) (*vos.SyntheticReport, error) {
	return m.OnGetSyntheticReport(ctx, accountPath, level, startTime, endTime)
}

func SuccessfulTransactionMock() UseCase {
	return UseCase{
		OnCreateTransaction: func(ctx context.Context, transaction entities.Transaction) error {
			return nil
		},
		OnGetAccountBalance: func(ctx context.Context, account vos.AccountPath) (vos.AccountBalance, error) {
			return vos.AccountBalance{
				Account:        vos.AccountPath{},
				CurrentVersion: 0,
				TotalCredit:    0,
				TotalDebit:     0,
			}, nil
		},
		OnGetAnalyticalData: func(ctx context.Context, query vos.AccountQuery, fn func(vos.Statement) error) error {
			return nil
		},
		OnGetSyntheticReport: func(ctx context.Context, account vos.AccountPath, level int, startTime time.Time, endTime time.Time) (*vos.SyntheticReport, error) {
			return &vos.SyntheticReport{
				Paths:       []vos.Path{},
				TotalCredit: 0,
				TotalDebit:  0,
			}, nil
		},
	}
}
