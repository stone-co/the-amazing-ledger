package usecases

import (
	"context"
)

func (l *LedgerUseCase) LoadObjectsIntoCache(ctx context.Context) error {
	var err error
	l.lastVersion, err = l.repository.LoadObjectsIntoCache(ctx, l.cachedAccounts)
	if err != nil {
		return err
	}

	// 0 and 1 have a special meaning. Setting to 1, the first version will be 2.
	if l.lastVersion == 0 {
		l.lastVersion = 1
	}

	l.log.Infof("Last version: %d", l.lastVersion)

	return nil
}
