package vos

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func Test_AccountHistoryTotalDebitAndCredit(t *testing.T) {
	entryHistoryOne, err := NewEntryHistory(CreditOperation, 100, time.Now())
	assert.Nil(t, err)

	entryHistoryTwo, _ := NewEntryHistory(CreditOperation, 100, time.Now())
	assert.Nil(t, err)

	entryHistoryThree, _ := NewEntryHistory(DebitOperation, 100, time.Now())
	assert.Nil(t, err)

	entryHistoryFour, _ := NewEntryHistory(DebitOperation, 100, time.Now())
	assert.Nil(t, err)

	totalCredit := entryHistoryOne.Amount + entryHistoryTwo.Amount
	totalDebit := entryHistoryThree.Amount + entryHistoryFour.Amount

	entriesHistory := []EntryHistory{*entryHistoryOne, *entryHistoryTwo, *entryHistoryThree, *entryHistoryFour}

	accountName, err := NewAccountName("liability:clients:available:account")
	assert.Nil(t, err)

	accountHistory, err := NewAccountHistory(*accountName, entriesHistory...)

	assert.Equal(t, accountHistory.TotalCredit, totalCredit)
	assert.Equal(t, accountHistory.TotalDebit, totalDebit)
	assert.Nil(t, err)
}
