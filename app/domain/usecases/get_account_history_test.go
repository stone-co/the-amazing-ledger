package usecases

import (
	"context"
	"testing"
	"time"

	"github.com/stone-co/the-amazing-ledger/app/domain/vos"
	"github.com/stretchr/testify/assert"
)

func TestLedgerUseCase_GetAccountHistory(t *testing.T) {
	t.Run("The AccountHistory must sum debits and credit", func(t *testing.T) {
		accountName, err := vos.NewAccountName("liability:stone:clients:user-1")
		assert.Nil(t, err)

		entriesHistory := make([]vos.EntryHistory, 4)
		for i := range entriesHistory {
			entryHistory, _ := vos.NewEntryHistory(vos.CreditOperation, 150, time.Now())
			entriesHistory[i] = *entryHistory
		}

		fakeAccountHistory, _ := vos.NewAccountHistory(*accountName, entriesHistory...)
		useCase := newFakeGetAccountHistory(&fakeAccountHistory, nil)
		accountHistory, err := useCase.GetAccountHistory(context.Background(), *accountName)

		assert.Nil(t, err)
		assert.Equal(t, fakeAccountHistory.Account.Name(), accountHistory.Account.Name())
		assert.Equal(t, fakeAccountHistory.EntriesHistory, accountHistory.EntriesHistory)

		assert.Equal(t, fakeAccountHistory.TotalCredit, accountHistory.TotalCredit)
		assert.Equal(t, fakeAccountHistory.TotalDebit, accountHistory.TotalDebit)
	})
}
