package usecase

import "github.com/stone-co/the-amazing-ledger/pkg/command-handler/domain/ledger/entities"

func (l LedgerUseCase) GetLastVersion() entities.Version {
	return l.lastVersion
}
