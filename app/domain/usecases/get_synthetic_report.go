package usecases

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/stone-co/the-amazing-ledger/app/domain/vos"
)

func (l *LedgerUseCase) GetSyntheticReport(ctx context.Context, query vos.Account, level int, startTime time.Time, endTime time.Time) (*vos.SyntheticReport, error) {
	if level < 1 {
		level = len(strings.Split(query.Value(), "."))
	}

	syntheticReport, err := l.repository.GetSyntheticReport(ctx, query, level, startTime, endTime)
	if err != nil {
		return nil, fmt.Errorf("failed to get synthetic report: %w", err)
	}

	return syntheticReport, nil
}
