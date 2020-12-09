package entities

import (
	"testing"

	"github.com/google/uuid"
	"github.com/stone-co/the-amazing-ledger/app/domain/vo"
	"github.com/stretchr/testify/assert"
)

func TestCachedAccounts_LoadOrStore(t *testing.T) {
	c := NewCachedAccounts()

	accountID := uuid.New().String()

	t.Run("New accounts started with version 1", func(t *testing.T) {
		accountInfo := c.LoadOrStore(accountID)
		assert.Equal(t, vo.NewAccountVersion, accountInfo.CurrentVersion)
	})

	t.Run("Account info is saved successfully", func(t *testing.T) {
		var version vo.Version = 1234
		c.Store(accountID, version)
		accountInfo := c.LoadOrStore(accountID)
		assert.Equal(t, version, accountInfo.CurrentVersion)
	})
}
