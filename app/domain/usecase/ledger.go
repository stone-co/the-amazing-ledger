package usecase

import (
	"github.com/sirupsen/logrus"
	"github.com/stone-co/the-amazing-ledger/app/domain"
	"github.com/stone-co/the-amazing-ledger/app/domain/entities"
	"github.com/stone-co/the-amazing-ledger/app/domain/vo"
)

var _ domain.UseCase = &LedgerUseCase{}

type LedgerUseCase struct {
	log            *logrus.Logger
	repository     domain.Repository
	cachedAccounts *entities.CachedAccounts
	lastVersion    vo.Version
}

func NewLedgerUseCase(log *logrus.Logger, repository domain.Repository) *LedgerUseCase {
	return &LedgerUseCase{
		log:            log,
		repository:     repository,
		cachedAccounts: entities.NewCachedAccounts(),
		lastVersion:    1,
	}
}
