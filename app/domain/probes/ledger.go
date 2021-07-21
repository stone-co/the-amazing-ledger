package probes

import (
	"context"

	"github.com/newrelic/go-agent/v3/newrelic"
	"github.com/sirupsen/logrus"
	"github.com/stone-co/the-amazing-ledger/app/domain"
)

var _ domain.Probe = &LedgerProbe{}

type LedgerProbe struct {
	logger   *logrus.Logger
	newrelic *newrelic.Application
}

func NewLedgerProbe(l *logrus.Logger, nr *newrelic.Application) *LedgerProbe {
	return &LedgerProbe{
		logger:   l,
		newrelic: nr,
	}
}

func (lp LedgerProbe) Log(ctx context.Context, value string) {
	lp.logger.Infof(value)
}

func (lp LedgerProbe) MonitorDataSegment(ctx context.Context, collection, operation, query string) *newrelic.DatastoreSegment {
	txn := newrelic.FromContext(ctx)
	seg := &newrelic.DatastoreSegment{
		Product:            newrelic.DatastorePostgres,
		Collection:         collection,
		Operation:          operation,
		ParameterizedQuery: query,
	}
	seg.StartTime = txn.StartSegmentNow()
	return seg
}
