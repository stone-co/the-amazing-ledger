package usecases

import (
	"context"

	"github.com/sirupsen/logrus"

	"github.com/stone-co/the-amazing-ledger/app/domain/mocks"
	"github.com/stone-co/the-amazing-ledger/app/domain/vos"
)

func newFakeGetAnalyticalData(entries []vos.Statement, result error) *LedgerUseCase {
	log := logrus.New()

	mockRepository := &mocks.Repository{
		OnGetAnalyticalData: func(ctx context.Context, query vos.AccountQuery, fn func(vos.Statement) error) error {
			for _, entry := range entries {
				if err := fn(entry); err != nil {
					return err
				}
			}
			return result
		},
	}

	return NewLedgerUseCase(log, mockRepository)
}
