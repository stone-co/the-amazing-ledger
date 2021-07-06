package mocks

import (
	"context"
	"time"

	"github.com/stone-co/the-amazing-ledger/app/domain"
	"github.com/stone-co/the-amazing-ledger/app/domain/entities"
	"github.com/stone-co/the-amazing-ledger/app/domain/vos"
)

var _ domain.Repository = &Repository{}

type Repository struct {
	OnCreateTransaction  func(context.Context, entities.Transaction) error
	OnGetAccountBalance  func(ctx context.Context, account vos.AccountPath) (vos.AccountBalance, error)
	OnGetAnalyticalData  func(ctx context.Context, query vos.AccountQuery, fn func(vos.Statement) error) error
	OnGetAccountHistory  func(ctx context.Context, account vos.AccountPath, fn func(vos.EntryHistory) error) error
	OnGetSyntheticReport func(ctx context.Context, accountPath vos.AccountPath, level int, startTime time.Time, endTime time.Time) (*vos.SyntheticReport, error)
}

func (s Repository) CreateTransaction(ctx context.Context, transaction entities.Transaction) error {
	return s.OnCreateTransaction(ctx, transaction)
}

func (s Repository) GetAccountBalance(ctx context.Context, account vos.AccountPath) (vos.AccountBalance, error) {
	return s.OnGetAccountBalance(ctx, account)
}

func (s Repository) GetAnalyticalData(ctx context.Context, query vos.AccountQuery, fn func(vos.Statement) error) error {
	return s.OnGetAnalyticalData(ctx, query, fn)
}

func (s Repository) GetAccountHistory(ctx context.Context, account vos.AccountPath, fn func(vos.EntryHistory) error) error {
	return s.OnGetAccountHistory(ctx, account, fn)
}

func (s Repository) GetSyntheticReport(ctx context.Context, accountPath vos.AccountPath, level int, startTime time.Time, endTime time.Time) (*vos.SyntheticReport, error) {
	return s.OnGetSyntheticReport(ctx, accountPath, level, startTime, endTime)
}
