package usecases

import (
	"github.com/sirupsen/logrus"

	"github.com/stone-co/the-amazing-ledger/app/domain"
	"github.com/stone-co/the-amazing-ledger/app/domain/probes"
)

var _ domain.UseCase = &LedgerUseCase{}

type LedgerUseCase struct {
	log        *logrus.Logger
	probe      *probes.LedgerProbe
	repository domain.Repository
}

func NewLedgerUseCase(log *logrus.Logger, repository domain.Repository, probe *probes.LedgerProbe) *LedgerUseCase {
	return &LedgerUseCase{
		log:        log,
		repository: repository,
		probe:      probe,
	}
}
