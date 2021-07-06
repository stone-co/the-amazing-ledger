package usecases

import (
	"context"
	"strings"

	"time"

	"github.com/stone-co/the-amazing-ledger/app"
	"github.com/stone-co/the-amazing-ledger/app/domain/vos"
)

func (l *LedgerUseCase) GetSyntheticReport(ctx context.Context, accountPath vos.AccountPath, level int, startTime time.Time, endTime time.Time) (*vos.SyntheticReport, error) {

	if (vos.AccountPath{}) == accountPath || accountPath.Name() == "" {
		return &vos.SyntheticReport{}, app.ErrEmptyAccountPath
	}

	if level < 1 {
		level = len(strings.Split(accountPath.Name(), ".")) // TODO get separator from default config
	}

	syntheticReport, err := l.repository.GetSyntheticReport(ctx, accountPath, level, startTime, endTime)

	if err != nil {
		return nil, err
	}

	return syntheticReport, nil
}
