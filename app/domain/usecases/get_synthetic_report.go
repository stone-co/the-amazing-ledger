package usecases

import (
	"context"

	"time"

	"github.com/jackc/pgx/v4"

	"github.com/stone-co/the-amazing-ledger/app"
	"github.com/stone-co/the-amazing-ledger/app/domain/vos"
)

func (l *LedgerUseCase) GetSyntheticReport(ctx context.Context, accountName string, startTime time.Time, endTime time.Time) (*vos.SyntheticReport, error) {

	syntheticReport, err := l.repository.GetSyntheticReport(ctx, accountName, startTime, endTime)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, app.ErrAccountNotFound
		}
		return nil, err
	}

	return syntheticReport, nil
}
