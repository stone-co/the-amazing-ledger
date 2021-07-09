package usecases

import (
	"context"
	"strings"

	"time"

	"github.com/stone-co/the-amazing-ledger/app/domain/vos"
)

func (l *LedgerUseCase) GetSyntheticReport(ctx context.Context, query vos.AccountQuery, level int, startTime time.Time, endTime time.Time) (*vos.SyntheticReport, error) {
	if level < 1 {
		level = len(strings.Split(query.Value(), vos.DepthSeparator))
	}

	syntheticReport, err := l.repository.GetSyntheticReport(ctx, query, level, startTime, endTime)

	if err != nil {
		return nil, err
	}

	return syntheticReport, nil
}
