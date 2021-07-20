package probes

import (
	"context"
	"fmt"
	"time"

	"github.com/stone-co/the-amazing-ledger/app/domain/vos"
)

const tag = "SyntheticReport[usecase]"

func (i *LedgerProbe) GettingSyntheticReport(ctx context.Context, account vos.AccountQuery, startTime time.Time, endTime time.Time) {
	i.Log(ctx, fmt.Sprintf("%s:%s/%d-%d", tag, account.Value(), startTime, endTime))
}

func (i *LedgerProbe) GotSyntheticReport(ctx context.Context, value string, report vos.SyntheticReport) {
	i.Log(ctx, fmt.Sprintf("%s:%s, tc = %d, td = %d", tag, value, report.TotalCredit, report.TotalDebit))
}
