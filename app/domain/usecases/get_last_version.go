package usecases

import "github.com/stone-co/the-amazing-ledger/app/domain/vos"

func (l LedgerUseCase) GetLastVersion() vos.Version {
	return l.lastVersion
}
