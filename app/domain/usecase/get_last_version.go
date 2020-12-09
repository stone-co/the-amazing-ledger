package usecase

import "github.com/stone-co/the-amazing-ledger/app/domain/vo"

func (l LedgerUseCase) GetLastVersion() vo.Version {
	return l.lastVersion
}
