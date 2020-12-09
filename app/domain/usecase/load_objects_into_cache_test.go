package usecase

import (
	"context"
	"testing"

	"github.com/stone-co/the-amazing-ledger/app/domain/vo"
	"github.com/stretchr/testify/assert"
)

func TestLedgerUseCase_LoadObjectsIntoCache(t *testing.T) {
	t.Run("The global version must be 1 when the entries table is empty", func(t *testing.T) {
		maxVersion := vo.Version(0)
		useCase := newFakeLoadObjectsIntoCacheUseCase(maxVersion, nil)
		err := useCase.LoadObjectsIntoCache(context.Background())
		assert.Nil(t, err)
		assert.Equal(t, vo.Version(1), useCase.GetLastVersion())
	})

	t.Run("The global version must be the last found in the entries table", func(t *testing.T) {
		maxVersion := vo.Version(12345)
		useCase := newFakeLoadObjectsIntoCacheUseCase(maxVersion, nil)
		err := useCase.LoadObjectsIntoCache(context.Background())
		assert.Nil(t, err)
		assert.Equal(t, maxVersion, useCase.GetLastVersion())
	})
}
