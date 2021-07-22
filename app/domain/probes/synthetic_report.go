package probes

import (
	"context"
	"fmt"
	"time"

	"github.com/stone-co/the-amazing-ledger/app/domain/vos"
)

func (lp *LedgerProbe) GettingSyntheticReport(ctx context.Context, account vos.Account, startTime time.Time, endTime time.Time) {
	lp.Log(ctx, fmt.Sprintf("%s/%s-%s", account.Value(), startTime.String(), endTime.String()))
}

func (lp *LedgerProbe) GotSyntheticReport(ctx context.Context, value string, report vos.SyntheticReport) {
	lp.Log(ctx, fmt.Sprintf("%s, tc = %d, td = %d", value, report.TotalCredit, report.TotalDebit))
}