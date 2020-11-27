package usecase

import (
	"github.com/sirupsen/logrus"
	"github.com/stone-co/the-amazing-ledger/pkg/command-handler/domain/ledger"
	"github.com/stone-co/the-amazing-ledger/pkg/command-handler/domain/ledger/entities"
)

var _ ledger.UseCase = &LedgerUseCase{}

type LedgerUseCase struct {
	log            *logrus.Logger
	repository     ledger.Repository
	cachedAccounts *entities.CachedAccounts
	lastVersion    entities.Version
}

func NewLedgerUseCase(log *logrus.Logger, repository ledger.Repository) *LedgerUseCase {
	return &LedgerUseCase{
		log:            log,
		repository:     repository,
		cachedAccounts: entities.NewCachedAccounts(),
		lastVersion:    1,
	}
}
