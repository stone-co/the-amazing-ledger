package usecases

import (
	"context"
	"testing"
	"time"

	"github.com/stone-co/the-amazing-ledger/app/domain/vos"
	"github.com/stretchr/testify/assert"
)

func TestLedgerUseCase_GetAccountHistory(t *testing.T) {
	res := []vos.EntryHistory{}
	fn := func(et vos.EntryHistory) error {
		res = append(res, et)
		return nil
	}

	t.Run("The account history can be empty", func(t *testing.T) {
		accountName, err := vos.NewAccountName("liability:stone:clients:user-1")
		assert.Nil(t, err)

		entries := []vos.EntryHistory{}

		// clear res
		res = []vos.EntryHistory{}
		useCase := newFakeGetAccountHistory(entries, nil)
		err = useCase.GetAccountHistory(context.Background(), *accountName, fn)

		assert.Nil(t, err)
		assert.Equal(t, entries, res)
	})

	t.Run("The account history should don't modify any Entry History", func(t *testing.T) {
		accountName, err := vos.NewAccountName("liability:stone:clients:user-2")
		assert.Nil(t, err)

		entries := make([]vos.EntryHistory, 4)
		for i := range entries {
			entryHistory, _ := vos.NewEntryHistory(vos.CreditOperation, 150, time.Now())
			entries[i] = *entryHistory
		}

		// clear res
		res = []vos.EntryHistory{}
		useCase := newFakeGetAccountHistory(entries, nil)
		err = useCase.GetAccountHistory(context.Background(), *accountName, fn)

		assert.Nil(t, err)
		assert.Equal(t, entries, res)
	})
}
