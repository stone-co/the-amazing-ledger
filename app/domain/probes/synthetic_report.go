package probes

import (
	"context"
	"fmt"
	"time"

	"github.com/stone-co/the-amazing-ledger/app/domain/vos"
)

const tag = "SyntheticReport[usecase]"

func (lp *LedgerProbe) GettingSyntheticReport(ctx context.Context, account vos.AccountQuery, startTime time.Time, endTime time.Time) {
	lp.Log(ctx, fmt.Sprintf("%s:%s/%s-%s", tag, account.Value(), startTime.String(), endTime.String()))
}

func (lp *LedgerProbe) GotSyntheticReport(ctx context.Context, value string, report vos.SyntheticReport) {
	lp.Log(ctx, fmt.Sprintf("%s:%s, tc = %d, td = %d", tag, value, report.TotalCredit, report.TotalDebit))
}
