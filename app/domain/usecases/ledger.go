package usecases

import (
	"github.com/stone-co/the-amazing-ledger/app/domain"
	"github.com/stone-co/the-amazing-ledger/app/domain/probes"
)

var _ domain.UseCase = &LedgerUseCase{}

type LedgerUseCase struct {
	probe      *probes.LedgerProbe
	repository domain.Repository
}

func NewLedgerUseCase(repository domain.Repository, probe *probes.LedgerProbe) *LedgerUseCase {
	return &LedgerUseCase{
		repository: repository,
		probe:      probe,
	}
}
