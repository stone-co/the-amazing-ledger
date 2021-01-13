package vos

import (
	"testing"
	"time"

	"github.com/stone-co/the-amazing-ledger/app"
	"github.com/stretchr/testify/assert"
)

func Test_EntryHistory(t *testing.T) {
	freezeTime := time.Now()
	entryHistory, err := NewEntryHistory(CreditOperation, 100, freezeTime)

	assert.Nil(t, err)
	assert.Equal(t, freezeTime, entryHistory.CreatedAt)
	assert.Equal(t, 100, entryHistory.Amount)
	assert.Equal(t, CreditOperation, entryHistory.Operation)
}

func Test_EntryHistoryInvalidAmount(t *testing.T) {
	_, err := NewEntryHistory(CreditOperation, 0, time.Now())
	assert.Equal(t, app.ErrInvalidAmount, err)

	_, err = NewEntryHistory(CreditOperation, -1, time.Now())
	assert.Equal(t, app.ErrInvalidAmount, err)
}

func Test_EntryHistoryInvalidOperation(t *testing.T) {
	_, err := NewEntryHistory(InvalidOperation, 100, time.Now())
	assert.Equal(t, app.ErrInvalidOperation, err)
}
