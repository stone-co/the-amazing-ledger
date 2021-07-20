package probes

import (
	"context"
)

func (lp LedgerProbe) Log(ctx context.Context, value string) {
	lp.logger.Infoln(value)
}
