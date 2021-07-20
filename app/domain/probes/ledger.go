package probes

import (
	"github.com/sirupsen/logrus"
	"github.com/stone-co/the-amazing-ledger/app/domain"
)

var _ domain.Instrumentation = &LedgerProbe{}

type LedgerProbe struct {
	logger *logrus.Logger
}

func NewLedgerProbe(logger *logrus.Logger) *LedgerProbe {
	return &LedgerProbe{
		logger: logger,
	}
}
